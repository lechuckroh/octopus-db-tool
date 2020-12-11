package main

import (
	"github.com/lechuckroh/octopus-db-tools/format/dbml"
	"github.com/lechuckroh/octopus-db-tools/format/gorm"
	"github.com/lechuckroh/octopus-db-tools/format/graphql"
	"github.com/lechuckroh/octopus-db-tools/format/jpa"
	"github.com/lechuckroh/octopus-db-tools/format/liquibase"
	"github.com/lechuckroh/octopus-db-tools/format/mysql"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/format/plantuml"
	"github.com/lechuckroh/octopus-db-tools/format/protobuf"
	"github.com/lechuckroh/octopus-db-tools/format/quickdbd"
	"github.com/lechuckroh/octopus-db-tools/format/sqlalchemy"
	"github.com/lechuckroh/octopus-db-tools/format/staruml"
	"github.com/lechuckroh/octopus-db-tools/format/xlsx"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"time"
)

const VERSION = "2.0.0-beta2"

var buildDateVersion string

func initCommand() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Action: octopus.InitAction,
		Flags:  octopus.InitCliFlags,
	}
}

func importCommand() *cli.Command {
	return &cli.Command{
		Name: "import",
		Subcommands: []*cli.Command{
			{
				Name:   "mysql",
				Action: mysql.ImportAction,
				Flags:  mysql.ImportCliFlags,
			},
			{
				Name:   "staruml",
				Action: staruml.ImportAction,
				Flags:  staruml.ImportCliFlags,
			},
			{
				Name:   "xlsx",
				Action: xlsx.ImportAction,
				Flags:  xlsx.ImportCliFlags,
			},
		},
	}
}

func exportCommand() *cli.Command {
	return &cli.Command{
		Name: "export",
		Subcommands: []*cli.Command{
			{
				Name:   "dbml",
				Action: dbml.ExportAction,
				Flags:  dbml.ExportCliFlags,
			},
			{
				Name:   "quickdbd",
				Action: quickdbd.ExportAction,
				Flags:  quickdbd.ExportCliFlags,
			},
			{
				Name:   "mysql",
				Action: mysql.ExportAction,
				Flags:  mysql.ExportCliFlags,
			},
			{
				Name:   "xlsx",
				Action: xlsx.ExportAction,
				Flags:  xlsx.ExportCliFlags,
			},
		},
	}
}

func generateCommand() *cli.Command {
	return &cli.Command{
		Name: "generate",
		Subcommands: []*cli.Command{
			{
				Name:   "gorm",
				Action: gorm.Action,
				Flags:  gorm.CliFlags,
			},
			{
				Name:   "graphql",
				Action: graphql.Action,
				Flags:  graphql.CliFlags,
			},
			{
				Name:   "kt",
				Action: jpa.KotlinAction,
				Flags:  jpa.KotlinCliFlags,
			},
			{
				Name:   "liquibase",
				Action: liquibase.Action,
				Flags:  liquibase.CliFlags,
			},
			{
				Name:   "plantuml",
				Action: plantuml.Action,
				Flags:  plantuml.CliFlags,
			},
			{
				Name:   "pb",
				Action: protobuf.Action,
				Flags:  protobuf.CliFlags,
			},
			{
				Name:   "sqlalchemy",
				Action: sqlalchemy.Action,
				Flags:  sqlalchemy.CliFlags,
			},
		},
	}
}

func main() {
	cliApp := cli.NewApp()
	cliApp.EnableBashCompletion = true
	cliApp.Name = "oct"
	cliApp.Version = VERSION + buildDateVersion
	cliApp.Compiled = time.Now()
	cliApp.Authors = []*cli.Author{
		{
			Name:  "Lechuck Roh",
			Email: "lechuckroh@gmail.com",
		},
	}
	cliApp.Copyright = "(c) 2019-2020 Lechuck Roh"
	cliApp.Usage = "octopus-db-tools"
	cliApp.Commands = []*cli.Command{
		initCommand(),
		importCommand(),
		exportCommand(),
		generateCommand(),
	}
	//
	//	{
	//		Name:    "create",
	//		Aliases: []string{"c"},
	//		Usage:   "create `filename`",
	//		ExportAction:  create,
	//	},
	//	{
	//		Name:    "convert",
	//		Aliases: []string{"c"},
	//		Usage:   "convert `source` `target`",
	//		Flags: []cli.Flag{
	//			cli.StringFlag{
	//				Name:   FlagSourceFormat,
	//				Usage:  "set source format",
	//				EnvVar: "OCTOPUS_SOURCE_FORMAT",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagTargetFormat,
	//				Usage:  "set target format",
	//				EnvVar: "OCTOPUS_TARGET_FORMAT",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagNotNull,
	//				Usage:  "use 'not null' instead of 'nullable'",
	//				EnvVar: "OCTOPUS_NOT_NULL",
	//			},
	//		},
	//		ExportAction: convert,
	//	},
	//	{
	//		Name:    "generate",
	//		Aliases: []string{"g"},
	//		Usage:   "generate `source` `target`",
	//		Flags: []cli.Flag{
	//			cli.StringFlag{
	//				Name:   FlagSourceFormat,
	//				Usage:  "set source format",
	//				EnvVar: "OCTOPUS_SOURCE_FORMAT",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagTargetFormat,
	//				Usage:  "set target format",
	//				EnvVar: "OCTOPUS_TARGET_FORMAT",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagPackage,
	//				Usage:  "set target package name",
	//				EnvVar: "OCTOPUS_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagReposPackage,
	//				Usage:  "set target repository package name",
	//				EnvVar: "OCTOPUS_REPOS_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagRelation,
	//				Usage:  "set relation annotation type",
	//				EnvVar: "OCTOPUS_RELATION",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagAnnotation,
	//				Usage:  "add custom class annotations",
	//				EnvVar: "OCTOPUS_ANNOTATION",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagRemovePrefix,
	//				Usage:  "set prefixes to remove. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_REMOVE_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagPrefix,
	//				Usage:  "set prefix to add",
	//				EnvVar: "OCTOPUS_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagUniqueNameSuffix,
	//				Usage:  "set unique constraint name suffix",
	//				EnvVar: "OCTOPUS_UNIQUE_NAME_SUFFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagGroups,
	//				Usage:  "filter table groups to generate. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_GROUPS",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagDiff,
	//				Usage:  "diff octopus filename.",
	//				EnvVar: "OCTOPUS_DIFF",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagIdEntity,
	//				Usage:  "set IdEntity interface name",
	//				EnvVar: "OCTOPUS_ID_ENTITY",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagUseUTC,
	//				Usage:  "use UTC for audit column default value",
	//				EnvVar: "OCTOPUS_USE_UTC",
	//			},
	//			cli.StringFlag{
	//				Name:   FlagUseComments,
	//				Usage:  "generate column comments",
	//				EnvVar: "OCTOPUS_USE_COMMENTS",
	//			},
	//		},
	//		ExportAction: generate,
	//	},
	//	{
	//		Name:  "jpa-kotlin",
	//		Usage: "jpa-kotlin `source` `target`",
	//		Flags: []cli.Flag{
	//			cli.StringFlag{
	//				Name:   jpa.FlagPackage,
	//				Usage:  "set kotlin package name",
	//				EnvVar: "OCTOPUS_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagReposPackage,
	//				Usage:  "set target repository package name",
	//				EnvVar: "OCTOPUS_REPOS_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagRelation,
	//				Usage:  "set virtual relation annotation type",
	//				EnvVar: "OCTOPUS_RELATION",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagAnnotation,
	//				Usage:  "add custom kotlin class annotations",
	//				EnvVar: "OCTOPUS_ANNOTATION",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagRemovePrefix,
	//				Usage:  "set prefixes to remove from kotlin class name. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_REMOVE_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagPrefix,
	//				Usage:  "set kotlin class name prefix",
	//				EnvVar: "OCTOPUS_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagUniqueNameSuffix,
	//				Usage:  "set unique constraint name suffix",
	//				EnvVar: "OCTOPUS_UNIQUE_NAME_SUFFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagGroups,
	//				Usage:  "filter table groups to generate. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_GROUPS",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagIdEntity,
	//				Usage:  "set kotlin interface name with `id` field.",
	//				EnvVar: "OCTOPUS_ID_ENTITY",
	//			},
	//			cli.StringFlag{
	//				Name:   jpa.FlagUseUTC,
	//				Usage:  "use UTC for audit column default value",
	//				EnvVar: "OCTOPUS_USE_UTC",
	//			},
	//		},
	//		ExportAction: generateJpaKotlin,
	//	},
	//	{
	//		Name:  "protobuf",
	//		Usage: "protobuf `source` `target`",
	//		Flags: []cli.Flag{
	//			cli.StringFlag{
	//				Name:   protobuf.FlagPackage,
	//				Usage:  "set package name",
	//				EnvVar: "OCTOPUS_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   protobuf.FlagGoPackage,
	//				Usage:  "set golang package name",
	//				EnvVar: "OCTOPUS_GO_PACKAGE",
	//			},
	//			cli.StringFlag{
	//				Name:   protobuf.FlagRemovePrefix,
	//				Usage:  "set prefixes to remove from message name. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_REMOVE_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   protobuf.FlagPrefix,
	//				Usage:  "set message name prefix",
	//				EnvVar: "OCTOPUS_PREFIX",
	//			},
	//			cli.StringFlag{
	//				Name:   protobuf.FlagGroups,
	//				Usage:  "filter table groups to generate. set multiple values with comma separated.",
	//				EnvVar: "OCTOPUS_GROUPS",
	//			},
	//		},
	//		ExportAction: generateProtobuf,
	//	},
	//}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
