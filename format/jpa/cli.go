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
	FlagUseUTC           = "useUTC"
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
	&cli.StringFlag{
		Name:    FlagAnnotation,
		Aliases: []string{"a"},
		Usage:   "Custom Entity class annotation. `FORMAT`: '{group1}:{annotations1}[,{group2}:{annotations2}]'",
		EnvVars: []string{"OCTOPUS_ANNOTATION"},
	},
	&cli.StringFlag{
		Name:    FlagGroups,
		Aliases: []string{"g"},
		Usage:   "Filter table groups to generate. `GROUPS` are separated by comma",
		EnvVars: []string{"OCTOPUS_GROUPS"},
	},
	&cli.StringFlag{
		Name:    FlagIdEntity,
		Aliases: []string{"e"},
		Usage:   "Interface `NAME` with 'id' field",
		EnvVars: []string{"OCTOPUS_ID_ENTITY"},
	},
	&cli.StringFlag{
		Name:    FlagPackage,
		Aliases: []string{"p"},
		Usage:   "Entity class `PACKAGE` name",
		EnvVars: []string{"OCTOPUS_PACKAGE"},
	},
	&cli.StringFlag{
		Name:    FlagPrefix,
		Aliases: []string{"f"},
		Usage:   "Class name prefix. `FORMAT`: '{group1}:{prefix1}[,{group2}:{prefix2}]'",
		EnvVars: []string{"OCTOPUS_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagRelation,
		Aliases: []string{"l"},
		Usage:   "Virtual relation `ANNOTATION` type. Available values: VRelation",
		EnvVars: []string{"OCTOPUS_RELATION"},
	},
	&cli.StringFlag{
		Name:    FlagRemovePrefix,
		Aliases: []string{"d"},
		Usage:   "Table `PREFIXES` to remove from class name. Multiple prefixes are separated by comma",
		EnvVars: []string{"OCTOPUS_REMOVE_PREFIX"},
	},
	&cli.StringFlag{
		Name:    FlagReposPackage,
		Aliases: []string{"r"},
		Usage:   "Repository class `PACKAGE` name. Generated if not empty.",
		EnvVars: []string{"OCTOPUS_REPOS_PACKAGE"},
	},
	&cli.StringFlag{
		Name:    FlagUniqueNameSuffix,
		Aliases: []string{"q"},
		Usage:   "Unique constraint name `SUFFIX`.",
		EnvVars: []string{"OCTOPUS_UNIQUE_NAME_SUFFIX"},
	},
	&cli.BoolFlag{
		Name:    FlagUseUTC,
		Aliases: []string{"u"},
		Usage:   "Set to use UTC for audit columns ('created_at', 'updated_at').",
		EnvVars: []string{"OCTOPUS_USE_UTC"},
	},
}
