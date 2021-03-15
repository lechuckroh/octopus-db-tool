package main

import (
	"github.com/lechuckroh/octopus-db-tools/format/dbml"
	"github.com/lechuckroh/octopus-db-tools/format/diff"
	"github.com/lechuckroh/octopus-db-tools/format/gorm"
	"github.com/lechuckroh/octopus-db-tools/format/graphql"
	"github.com/lechuckroh/octopus-db-tools/format/jpa"
	"github.com/lechuckroh/octopus-db-tools/format/liquibase"
	"github.com/lechuckroh/octopus-db-tools/format/mysql"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/format/ojson"
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

const VERSION = "2.0.0-beta3"

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
				Name:   "ojson",
				Action: ojson.ImportAction,
				Flags:  ojson.ImportCliFlags,
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

func diffCommand() *cli.Command {
	return &cli.Command{
		Name: "diff",
		Subcommands: []*cli.Command{
			{
				Name:   "flyway",
				Action: diff.FlywayAction,
				Flags:  diff.CliFlags,
			},
			{
				Name:   "liquibase",
				Action: diff.LiquibaseAction,
				Flags:  diff.CliFlags,
			},
			{
				Name:   "md",
				Action: diff.MarkdownAction,
				Flags:  diff.CliFlags,
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
		diffCommand(),
		initCommand(),
		importCommand(),
		exportCommand(),
		generateCommand(),
	}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
