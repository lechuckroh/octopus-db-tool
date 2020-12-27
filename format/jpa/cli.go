package jpa

import (
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/urfave/cli/v2"
	"strings"
)

const (
	FlagAnnotation       = "annotation"
	FlagGroups           = "groups"
	FlagIdEntity         = "idEntity"
	FlagInput            = "input"
	FlagOutput           = "output"
	FlagPackage          = "package"
	FlagPrefix           = "prefix"
	FlagRelation         = "relation"
	FlagRemovePrefix     = "removePrefix"
	FlagReposPackage     = "reposPackage"
	FlagUniqueNameSuffix = "uniqueNameSuffix"
)

func KotlinAction(c *cli.Context) error {
	schema, err := octopus.LoadSchema(c.String(FlagInput))
	if err != nil {
		return err
	}

	gen := NewKtGenerator(schema, &KtOption{
		AnnoMapper:       common.NewAnnotationMapper(c.String(FlagAnnotation)),
		PrefixMapper:     common.NewPrefixMapper(c.String(FlagPrefix)),
		TableFilter:      octopus.GetTableFilterFn(c.String(FlagGroups)),
		IdEntity:         c.String(FlagIdEntity),
		Package:          c.String(FlagPackage),
		Relation:         c.String(FlagRelation),
		RemovePrefixes:   strings.Split(c.String(FlagRemovePrefix), ","),
		ReposPackage:     c.String(FlagReposPackage),
		UniqueNameSuffix: c.String(FlagUniqueNameSuffix),
	})
	return gen.Generate(c.String(FlagOutput))
}

var KotlinCliFlags = []cli.Flag{
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
		Usage:    "generate kotlin files to `DIR`",
		EnvVars:  []string{"OCTOPUS_OUTPUT"},
		Required: true,
	},
	// TODO: add flags
}
