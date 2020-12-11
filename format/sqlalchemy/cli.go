package sqlalchemy

import (
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/urfave/cli/v2"
	"strings"
)

const (
	FlagGroups           = "groups"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagPrefix           = "prefix"
	FlagRemovePrefix     = "removePrefix"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
	FlagUseUTC           = "useUTC"
)

func Action(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := Generator{
		schema: schema,
		option: &Option{
			PrefixMapper:     common.NewPrefixMapper(c.String(FlagPrefix)),
			TableFilter:      octopus.GetTableFilterFn(c.String(FlagGroups)),
			RemovePrefixes:   strings.Split(c.String(FlagRemovePrefix), ","),
			UniqueNameSuffix: c.String(FlagUniqueNameSuffix),
			UseUTC:           c.Bool(FlagUseUTC),
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
		Usage:    "generate python files to `PATH`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "filter table groups to generate. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_GROUPS"},
	},
	&cli.StringFlag{
		Name:    FlagPrefix,
		Aliases: []string{"p"},
		Usage:   "set entity class name prefix",
		EnvVars: []string{"OCTOPUS_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagRemovePrefix,
		Aliases: []string{"r"},
		Usage:   "set prefixes to remove from entity class name. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_REMOVE_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagUniqueNameSuffix,
		Aliases: []string{"u"},
		Usage:   "set unique constraint name suffix",
		EnvVars: []string{"OCTOPUS_UNIQUE_NAME_SUFFIX"},
	},
	&cli.StringFlag{
		Name:    FlagUseUTC,
		Aliases: []string{"t"},
		Usage:   "use UTC for audit column default value",
		EnvVars: []string{"OCTOPUS_USE_UTC"},
	},
}
