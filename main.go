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
	FilePath string
	Format   string
	Options  map[string]string
}

func (o *Output) Get(name string) string {
	return o.Options[name]
}
func (o *Output) GetBool(name string) bool {
	return o.Options[name] == "true"
}
func (o *Output) GetSlice(name string) []string {
	return strings.Split(o.Options[name], ",")
}

func getFlagValues(c *cli.Context) map[string]string {
	result := make(map[string]string)

	for _, flagName := range c.FlagNames() {
		flagValue := c.String(flagName)
		result[flagName] = flagValue
	}

	return result
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
	return cmd.Create(&Output{
		FilePath: filename,
		Format:   "",
		Options:  getFlagValues(c),
	})
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

	inputFormat := GetFileFormat(c.String(FlagSourceFormat), inputFilename)
	if inputFormat == "" {
		return errors.New("cannot find sourceFormat")
	}
	outputFormat := GetFileFormat(c.String(FlagTargetFormat), outputFilename)
	if outputFormat == "" {
		return errors.New("cannot find targetFormat")
	}

	input := &Input{
		Filename: inputFilename,
		Format:   inputFormat,
	}
	output := &Output{
		FilePath: outputFilename,
		Format:   outputFormat,
		Options:  getFlagValues(c),
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
	inputFormat := GetFileFormat(c.String(FlagSourceFormat), inputFilename)
	if inputFormat == "" {
		return errors.New("cannot find sourceFormat")
	}

	input := &Input{
		Filename: inputFilename,
		Format:   inputFormat,
	}

	output := &Output{
		FilePath: args.Get(1),
		Format:   c.String(FlagTargetFormat),
		Options:  getFlagValues(c),
	}

	cmd := &GenerateCmd{}
	return cmd.Generate(input, output)
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "oct"
	cliApp.Version = "1.0.10"
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
					Name:   FlagSourceFormat,
					Usage:  "set source format",
					EnvVar: "OCTOPUS_SOURCE_FORMAT",
				},
				cli.StringFlag{
					Name:   FlagTargetFormat,
					Usage:  "set target format",
					EnvVar: "OCTOPUS_TARGET_FORMAT",
				},
				cli.StringFlag{
					Name:   FlagNotNull,
					Usage:  "use 'not null' instead of 'nullable'",
					EnvVar: "OCTOPUS_NOT_NULL",
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
					Name:   FlagSourceFormat,
					Usage:  "set source format",
					EnvVar: "OCTOPUS_SOURCE_FORMAT",
				},
				cli.StringFlag{
					Name:   FlagTargetFormat,
					Usage:  "set target format",
					EnvVar: "OCTOPUS_TARGET_FORMAT",
				},
				cli.StringFlag{
					Name:   FlagPackage,
					Usage:  "set target package name",
					EnvVar: "OCTOPUS_PACKAGE",
				},
				cli.StringFlag{
					Name:   FlagReposPackage,
					Usage:  "set target repository package name",
					EnvVar: "OCTOPUS_REPOS_PACKAGE",
				},
				cli.StringFlag{
					Name:   FlagRelation,
					Usage:  "set relation annotation type",
					EnvVar: "OCTOPUS_RELATION",
				},
				cli.StringFlag{
					Name:   FlagGraphqlPackage,
					Usage:  "set target graphql package name",
					EnvVar: "OCTOPUS_GRAPHQL_PACKAGE",
				},
				cli.StringFlag{
					Name:   FlagRemovePrefix,
					Usage:  "set prefixes to remove. set multiple values with comma separated.",
					EnvVar: "OCTOPUS_REMOVE_PREFIX",
				},
				cli.StringFlag{
					Name:   FlagUniqueNameSuffix,
					Usage:  "set unique constraint name suffix",
					EnvVar: "OCTOPUS_UNIQUE_NAME_SUFFIX",
				},
				cli.StringFlag{
					Name:   FlagGroups,
					Usage:  "filter table groups to generate",
					EnvVar: "OCTOPUS_GROUPS",
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
