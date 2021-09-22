package plantuml

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

const (
	FlagInput  = "input"
	FlagOutput = "output"
)

func Action(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := Generator{
		schema: schema,
		option: &Option{},
	}

	outputPath := c.String(FlagOutput)
	extSet := util.NewStringSet(".wsd", ".pu", ".puml", ".plantuml", ".iuml")
	var filename string
	if ext := strings.ToLower(filepath.Ext(outputPath)); extSet.Contains(ext) {
		filename = outputPath
	} else {
		// ensure directory is created
		if _, err := util.Mkdir(outputPath); err != nil {
			return err
		}
		filename = filepath.Join(outputPath, "output.plantuml")
	}

	buf := new(bytes.Buffer)
	if err = gen.Generate(buf); err != nil {
		return err
	}

	// write to file
	return util.WriteStringToFile(filename, buf.String())
}

var CliFlags = []cli.Flag{
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
		Usage:    "geneate plantUML to `FILE` or `DIRECTORY`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
}
