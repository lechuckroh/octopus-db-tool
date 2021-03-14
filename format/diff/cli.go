package diff

import (
	"github.com/lechuckroh/octopus-db-tools/format/liquibase"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/urfave/cli/v2"
)

const (
	FlagAuthor           = "author"
	FlagFrom             = "from"
	FlagGroups           = "groups"
	FlagOutput           = "output"
	FlagTo               = "to"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
	FlagUseComments      = "comments"
)

func LiquibaseAction(c *cli.Context) error {
	fromSchema, err := octopus.LoadSchema(c.String(FlagFrom))
	if err != nil {
		return err
	}
	toSchema, err := octopus.LoadSchema(c.String(FlagTo))
	if err != nil {
		return err
	}

	return liquibase.NewDiff(
		&liquibase.DiffOption{
			Author:  c.String(FlagAuthor),
			DiffFrom:    fromSchema,
			DiffTo:      toSchema,
			TableFilter: octopus.GetTableFilterFn(c.String(FlagGroups)),
			UseComments: c.Bool(FlagUseComments),
		},
	).Generate(c.String(FlagOutput))
}

var CliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    FlagAuthor,
		Aliases: []string{"a"},
		Usage:   "diff author",
		EnvVars: []string{"OCTOPUS_AUTHOR"},
	},
	&cli.StringFlag{
		Name:     FlagFrom,
		Aliases:  []string{"f"},
		Usage:    "octopus schema to compare 'from'",
		EnvVars:  []string{"OCTOPUS_FROM"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "filter table groups to compare. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_GROUPS"},
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "diff output `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagTo,
		Aliases:  []string{"t"},
		Usage:    "octopus schema to compare 'to'",
		EnvVars:  []string{"OCTOPUS_TO"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagUniqueNameSuffix,
		Aliases: []string{"u"},
		Usage:   "set unique constraint name suffix",
		EnvVars: []string{"OCTOPUS_UNIQUE_NAME_SUFFIX"},
	},
	&cli.BoolFlag{
		Name:    FlagUseComments,
		Aliases: []string{"c"},
		Usage:   "set true to compare column comments",
		EnvVars: []string{"OCTOPUS_COMMENTS"},
	},
}
