package liquibase

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/urfave/cli/v2"
)

const (
	FlagDiff             = "diff"
	FlagGroups           = "groups"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
	FlagUseComments      = "comments"
)

func Action(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := Generator{
		schema: schema,
		option: &Option{
			TableFilter:      octopus.GetTableFilterFn(c.String(FlagGroups)),
			DiffTarget:       c.String(FlagDiff),
			UniqueNameSuffix: c.String(FlagUniqueNameSuffix),
			UseComments:      c.Bool(FlagUseComments),
		},
	}
	return gen.Generate(c.String(FlagOutput))
}

var CliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagInput,
		Aliases:  []string{"i"},
		Usage:    "load input octopus schema from `FILE`",
		EnvVars:  []string{"OCTOPUS_INPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "export liquibase to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagDiff,
		Aliases: []string{"d"},
		Usage:   "",
		EnvVars: []string{"OCTOPUS_DIFF"},
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "filter table groups to generate. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_GROUPS"},
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
		Usage:   "set true to generate column comments",
		EnvVars: []string{"OCTOPUS_COMMENTS"},
	},
}
