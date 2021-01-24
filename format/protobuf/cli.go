package protobuf

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/urfave/cli/v2"
	"strings"
)

const (
	FlagGoPackage        = "goPackage"
	FlagGroups           = "groups"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagPackage          = "package"
	FlagPrefix           = "prefix"
	FlagRemovePrefix     = "removePrefix"
	FlogRelationTagDecr  = "relationTagDecr"
	FlogRelationTagStart = "relationTagStart"
)

func Action(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := newGenerator(
		schema,
		&Option{
			PrefixMapper:     common.NewPrefixMapper(c.String(FlagPrefix)),
			TableFilter:      octopus.GetTableFilterFn(c.String(FlagGroups)),
			RemovePrefixes:   strings.Split(c.String(FlagRemovePrefix), ","),
			Package:          c.String(FlagPackage),
			GoPackage:        c.String(FlagGoPackage),
			FilePath:         c.String(FlagOutput),
			RelationTagStart: -1,
			RelationTagDecr:  false,
		},
	)
	buf := new(bytes.Buffer)
	if err = gen.Generate(buf); err != nil {
		return err
	}

	// write to file
	return util.WriteStringToFile(c.String(FlagOutput), buf.String())
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
		Usage:    "generate protobuf definition to `FILE`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagGoPackage,
		Usage:   "set go package name",
		EnvVars: []string{"OCTOPUS_GO_PACKAGE"},
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "filter table groups to generate. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_GROUPS"},
	},
	&cli.StringFlag{
		Name:    FlagPackage,
		Aliases: []string{"p"},
		Usage:   "set package name",
		EnvVars: []string{"OCTOPUS_PACKAGE"},
	},
	&cli.StringFlag{
		Name:    FlagPrefix,
		Aliases: []string{"f"},
		Usage:   "set proto message name prefix",
		EnvVars: []string{"OCTOPUS_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagRemovePrefix,
		Aliases: []string{"d"},
		Usage:   "set prefixes to remove from message name. set multiple values with comma separated.",
		EnvVars: []string{"OCTOPUS_REMOVE_PREFIX"},
	},
	&cli.BoolFlag{
		Name:    FlogRelationTagDecr,
		Usage:   "set relation tags decremental from `relationTagStart`",
		EnvVars: []string{"OCTOPUS_RELATION_TAG_DECR"},
	},
	&cli.StringFlag{
		Name:    FlogRelationTagStart,
		Aliases: []string{"s"},
		Usage:   "set relation tags start index. set -1 to start from last of fields.",
		EnvVars: []string{"OCTOPUS_RELATION_TAG_START"},
	},
}
