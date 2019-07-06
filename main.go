package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
	"time"
)

func create(c *cli.Context) error {
	args := c.Args()
	argsCount := c.NArg()
	var filename string
	if argsCount == 0 {
		filename = "db.ojson"
	} else {
		filename = args.Get(0)
	}

	app := NewApp()
	return app.Create(filename)
}

func convert(c *cli.Context) error {
	args := c.Args()
	argsCount := c.NArg()
	if argsCount == 0 {
		return cli.NewExitError("source is not set", 1)
	}
	if argsCount == 1 {
		return cli.NewExitError("target is not set", 1)
	}

	source := args.Get(0)
	target := args.Get(1)
	sourceFormat := c.String("sourceFormat")
	targetFormat := c.String("targetFormat")

	app := NewApp()
	return app.Convert(source, sourceFormat, target, targetFormat)
}

func generate(c *cli.Context) error {
	args := c.Args()
	argsCount := c.NArg()
	if argsCount == 0 {
		return cli.NewExitError("source is not set", 1)
	}
	if argsCount == 1 {
		return cli.NewExitError("target is not set", 1)
	}

	source := args.Get(0)
	target := args.Get(1)
	format := c.String("format")

	app := NewApp()
	return app.Generate(source, target, format)
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "oct"
	cliApp.Version = "0.0.1"
	cliApp.Compiled = time.Now()
	cliApp.Authors = []cli.Author{
		{Name: "Lechuck Roh"},
	}
	cliApp.Copyright = "(c) 2019 Lechuck Roh"
	cliApp.Usage = "octopus-db-tools"
	cliApp.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create `filename`",
			Action: create,
		},
		{
			Name:    "convert",
			Aliases: []string{"c"},
			Usage:   "convert `source` `target`",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "sourceFormat, sf",
					Usage:  "set source format",
					EnvVar: "OCTOPUS_CONVERT_SOURCE_FORMAT",
				},
				cli.StringFlag{
					Name:   "targetFormat, tf",
					Usage:  "set target format",
					EnvVar: "OCTOPUS_CONVERT_TARGET_FORMAT",
				},
			},
			Action: convert,
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate `source` `target`",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Usage:  "set target format",
					EnvVar: "OCTOPUS_GENERATE_FORMAT",
				},
			},
			Action: generate,
		},
	}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
