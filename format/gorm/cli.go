package gorm

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

const (
	FlagGormModel        = "gormModel"
	FlagGroups           = "groups"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagPackage          = "package"
	FlagPrefix           = "prefix"
	FlagRemovePrefix     = "removePrefix"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
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
			GormModel:        c.String(FlagGormModel),
			Package:          c.String(FlagPackage),
			RemovePrefixes:   strings.Split(c.String(FlagRemovePrefix), ","),
			UniqueNameSuffix: c.String(FlagUniqueNameSuffix),
		},
	}

	outputPath := c.String(FlagOutput)
	var filename string
	if ext := strings.ToLower(filepath.Ext(outputPath)); ext == ".go" {
		filename = outputPath
	} else {
		// ensure directory is created
		if _, err := util.Mkdir(outputPath); err != nil {
			return err
		}
		filename = filepath.Join(outputPath, "output.go")
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
		Usage:    "generate GORM source file(s) to `FILE`/`DIR`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagGormModel,
		Aliases: []string{"m"},
		Usage:   "set embedded base model for GORM model",
		EnvVars: []string{"OCTOPUS_GORM_MODEL"},
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "filter table groups to generate. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_GROUPS"},
	},
	&cli.StringFlag{
		Name:    FlagPackage,
		Aliases: []string{"k"},
		Usage:   "set package name",
		EnvVars: []string{"OCTOPUS_PACKAGE"},
	},
	&cli.StringFlag{
		Name:    FlagPrefix,
		Aliases: []string{"p"},
		Usage:   "set model struct name prefix",
		EnvVars: []string{"OCTOPUS_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagRemovePrefix,
		Aliases: []string{"r"},
		Usage:   "set prefixes to remove from model struct name. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_REMOVE_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagUniqueNameSuffix,
		Aliases: []string{"u"},
		Usage:   "set unique constraint name suffix",
		EnvVars: []string{"OCTOPUS_UNIQUE_NAME_SUFFIX"},
	},
}
