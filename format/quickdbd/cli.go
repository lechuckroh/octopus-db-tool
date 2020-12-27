package quickdbd

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
)

const (
	FlagInput  = "input"
	FlagOutput = "output"
)

func ExportAction(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	exporter := Exporter{
		schema: schema,
		option: &ExportOption{},
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
		Usage:    "export quickdbd to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
}
