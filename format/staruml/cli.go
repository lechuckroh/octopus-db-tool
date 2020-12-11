package staruml

import (
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
)

const (
	FlagInput  = "input"
	FlagOutput = "output"
)

func ImportAction(c *cli.Context) error {
	importer := Importer{}
	schema, err := importer.ImportFile(c.String(FlagInput))
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
		Usage:    "import input starUML from `FILE`",
		EnvVars:  []string{"OCTOPUS_INPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagOutput,
		Aliases:  []string{"o"},
		Usage:    "output octopus schema to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
}
