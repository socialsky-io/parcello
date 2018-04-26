// Command Line Interface of Embedo.
package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/phogolabs/parcello"
	"github.com/urfave/cli"
)

const (
	// ErrCodeArg is returned when an invalid argument is passed to CLI
	ErrCodeArg = 101
)

func main() {
	app := &cli.App{
		Name:                 "parcello",
		HelpName:             "parcello",
		Usage:                "Golang Resource Bundler and Embedder",
		UsageText:            "parcello [global options]",
		Version:              "0.7",
		BashComplete:         cli.DefaultAppComplete,
		EnableBashCompletion: true,
		Writer:               os.Stdout,
		ErrWriter:            os.Stderr,
		Action:               run,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "quiet, q",
				Usage: "Disable logging",
			},
			cli.BoolFlag{
				Name:  "recursive, r",
				Usage: "Embed the resources recursively",
			},
			cli.StringFlag{
				Name:  "resource-dir, d",
				Usage: "Path to directory",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "bundle-path, b",
				Usage: "Path to the bundle",
				Value: ".",
			},
			cli.StringSliceFlag{
				Name:  "ignore, i",
				Usage: "Ignore file name",
			},
			cli.BoolTFlag{
				Name:  "include-docs",
				Usage: "Include API documentation in generated source code",
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	return embed(ctx)
}

func embed(ctx *cli.Context) error {
	resourceDir, err := filepath.Abs(ctx.GlobalString("resource-dir"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	bundlePath, err := filepath.Abs(ctx.GlobalString("bundle-path"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	_, packageName := filepath.Split(bundlePath)

	embedder := &parcello.Embedder{
		Logger:     logger(ctx),
		FileSystem: parcello.Dir(resourceDir),
		Composer: &parcello.Generator{
			FileSystem: parcello.Dir(bundlePath),
			Config: &parcello.GeneratorConfig{
				Package:     packageName,
				InlcudeDocs: ctx.BoolT("include-docs"),
			},
		},
		Compressor: &parcello.ZipCompressor{
			Config: &parcello.CompressorConfig{
				Logger:         logger(ctx),
				Filename:       "resource",
				IgnorePatterns: ctx.GlobalStringSlice("ignore"),
				Recurive:       ctx.GlobalBool("recursive"),
			},
		},
	}

	if err := embedder.Embed(); err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	return nil
}

func logger(ctx *cli.Context) io.Writer {
	if ctx.GlobalBool("quiet") {
		return ioutil.Discard
	}

	return os.Stdout
}
