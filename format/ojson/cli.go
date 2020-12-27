package ojson

import (
	"github.com/urfave/cli/v2"
)

const (
	FlagInput  = "input"
	FlagOutput = "output"
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
		Usage:    "import octopus v1 schema from `FILE`",
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
