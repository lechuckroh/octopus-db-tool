package xlsx

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
)

const (
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagUseNotNullColumn = "useNotNullColumn"
)

func ImportAction(c *cli.Context) error {
	importer := Importer{}
	schema, err := importer.Import(c.String(FlagInput))
	if err != nil {
		return err
	}

	// write to file
	if bytes, err := schema.ToJson(); err != nil {
		return err
	} else {
		return util.WriteBytesToFile(c.String(FlagOutput), bytes)
	}
}

var ImportCliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagInput,
		Aliases:  []string{"i"},
		Usage:    "import xlsx from `FILE`",
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
			UseNotNullColumn: c.Bool(FlagUseNotNullColumn),
		},
	}
	return exporter.Export(c.String(FlagOutput))
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
		Usage:    "export xlsx to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagUseNotNullColumn,
		Aliases: []string{"n"},
		Usage:   "use not null column",
		EnvVars: []string{"OCTOPUS_USE_NOT_NULL_COLUMN"},
	},
}
