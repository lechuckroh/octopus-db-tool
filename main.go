package main

import (
	"errors"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Input struct {
	Filename string
	Format   string
}

type Output struct {
	Filename string
	Format   string
}

type GenOutput struct {
	Format           string
	Directory        string
	Package          string
	ReposPackage     string
	Relation         string
	GraphqlPackage   string
	PrefixesToRemove []string
	UniqueNameSuffix string
}

func create(c *cli.Context) error {
	args := c.Args()
	argsCount := c.NArg()
	var filename string
	if argsCount == 0 {
		filename = "db.ojson"
	} else {
		filename = args.Get(0)
	}

	cmd := &CreateCmd{}
	return cmd.Create(&Output{Filename: filename})
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

	inputFilename := args.Get(0)
	outputFilename := args.Get(1)

	inputFormat := GetFileFormat(c.String("sourceFormat"), inputFilename)
	if inputFormat == "" {
		return errors.New("cannot find sourceFormat")
	}
	outputFormat := GetFileFormat(c.String("targetFormat"), outputFilename)
	if outputFormat == "" {
		return errors.New("cannot find targetFormat")
	}

	input := &Input{
		Filename: inputFilename,
		Format:   inputFormat,
	}
	output := &Output{
		Filename: outputFilename,
		Format:   outputFormat,
	}

	cmd := &ConvertCmd{}
	return cmd.Convert(input, output)
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

	inputFilename := args.Get(0)
	inputFormat := GetFileFormat(c.String("sourceFormat"), inputFilename)
	if inputFormat == "" {
		return errors.New("cannot find sourceFormat")
	}

	input := &Input{
		Filename: inputFilename,
		Format:   inputFormat,
	}

	prefixesToRemove := strings.Split(c.String("removePrefix"), ",")

	output := &GenOutput{
		Directory:        args.Get(1),
		Format:           c.String("targetFormat"),
		Package:          c.String("package"),
		ReposPackage:     c.String("reposPackage"),
		Relation:         c.String("relation"),
		GraphqlPackage:   c.String("graphqlPackage"),
		UniqueNameSuffix: c.String("uniqueNameSuffix"),
		PrefixesToRemove: prefixesToRemove,
	}

	cmd := &GenerateCmd{}
	return cmd.Generate(input, output)
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "oct"
	cliApp.Version = "1.0.7"
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
			Action:  create,
		},
		{
			Name:    "convert",
			Aliases: []string{"c"},
			Usage:   "convert `source` `target`",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "sourceFormat, sf",
					Usage:  "set source format",
					EnvVar: "OCTOPUS_SOURCE_FORMAT",
				},
				cli.StringFlag{
					Name:   "targetFormat, tf",
					Usage:  "set target format",
					EnvVar: "OCTOPUS_TARGET_FORMAT",
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
					Name:   "sourceFormat, sf",
					Usage:  "set source format",
					EnvVar: "OCTOPUS_SOURCE_FORMAT",
				},
				cli.StringFlag{
					Name:   "targetFormat, tf",
					Usage:  "set target format",
					EnvVar: "OCTOPUS_TARGET_FORMAT",
				},
				cli.StringFlag{
					Name:   "package, p",
					Usage:  "set target package name",
					EnvVar: "OCTOPUS_PACKAGE",
				},
				cli.StringFlag{
					Name:   "reposPackage",
					Usage:  "set target repository package name",
					EnvVar: "OCTOPUS_REPOS_PACKAGE",
				},
				cli.StringFlag{
					Name:   "relation",
					Usage:  "set relation annotation type",
					EnvVar: "OCTOPUS_RELATION",
				},
				cli.StringFlag{
					Name:   "graphqlPackage",
					Usage:  "set target graphql package name",
					EnvVar: "OCTOPUS_GRAPHQL_PACKAGE",
				},
				cli.StringFlag{
					Name:   "removePrefix",
					Usage:  "set prefixes to remove. set multiple values with comma separated.",
					EnvVar: "OCTOPUS_REMOVE_PREFIX",
				},
				cli.StringFlag{
					Name:   "uniqueNameSuffix",
					Usage:  "set unique constraint name suffix",
					EnvVar: "OCTOPUS_UNIQUE_NAME_SUFFIX",
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
