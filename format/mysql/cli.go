package mysql

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
)

const (
	FlagGroups           = "groups"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
)

func ImportAction(c *cli.Context) error {
	importer := Importer{
		option: &ImportOption{},
	}
	schema, err := importer.ImportFile(c.String(FlagInput))
	if err != nil {
		return err
	}

	// write to file
	return schema.ToFile(c.String(FlagOutput))
}

var ImportCliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagInput,
		Aliases:  []string{"i"},
		Usage:    "import mysql DDL from `FILE`",
		EnvVars:  []string{"OCTOPUS_INPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "write octopus schema to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
}

func ExportAction(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	exporter := Exporter{
		schema: schema,
		option: &ExportOption{
			TableFilter:      octopus.GetTableFilterFn(c.String(FlagGroups)),
			UniqueNameSuffix: c.String(FlagUniqueNameSuffix),
		},
	}
	buf := new(bytes.Buffer)
	if err = exporter.Export(buf); err != nil {
		return err
	}

	// write to file
	return util.WriteStringToFile(c.String(FlagOutput), buf.String())
}

var ExportCliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagInput,
		Aliases:  []string{"i"},
		Usage:    "read octopus schema from `FILE`",
		EnvVars:  []string{"OCTOPUS_INPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "export mysql DDL to `FILE`",
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
		Name:    FlagUniqueNameSuffix,
		Aliases: []string{"u"},
		Usage:   "set unique constraint name suffix",
		EnvVars: []string{"OCTOPUS_UNIQUE_NAME_SUFFIX"},
	},
}
