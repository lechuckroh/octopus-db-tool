package graphql

import (
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/urfave/cli/v2"
	"strings"
)

const (
	FlagGraphqlPackage = "graphqlPackage"
	FlagGroups         = "groups"
	FlagInput          = "input"
	FlagOutput         = "output"
	FlagPrefix         = "prefix"
	FlagRemovePrefix   = "removePrefix"
)

func Action(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := Generator{
		schema: schema,
		option: &Option{
			PrefixMapper:   common.NewPrefixMapper(c.String(FlagPrefix)),
			TableFilter:    octopus.GetTableFilterFn(c.String(FlagGroups)),
			RemovePrefixes: strings.Split(c.String(FlagRemovePrefix), ","),
		},
	}
	return gen.Generate(c.String(FlagOutput))
}

var CliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagInput,
		Aliases:  []string{"i"},
		Usage:    "input octopus schema `FILE`",
		EnvVars:  []string{"OCTOPUS_INPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "generate graphql to `DIR`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagGraphqlPackage,
		Aliases: []string{"p"},
		Usage:   "set target graphql package name",
		EnvVars: []string{"OCTOPUS_GRAPHQL_PACKAGE"},
	},
}
