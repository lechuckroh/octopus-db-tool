package common

import (
	"path/filepath"
	"strings"
)

const (
	FormatDbdiagramIo     = "dbdiagram.io"
	FormatGorm            = "gorm"
	FormatGraphql         = "graphql"
	FormatJpaKotlin       = "jpa-kotlin"
	FormatLiquibase       = "liquibase"
	FormatOctopus         = "octopus"
	FormatOptiStudio      = "opti-studio"
	FormatPlantuml        = "plantuml"
	FormatProtobuf        = "protobuf"
	FormatQuickdbd        = "quickdbd"
	FormatSchemaConverter = "schema-converter"
	FormatSqlalchemy      = "sqlalchemy"
	FormatSqlH2           = "h2"
	FormatSqlMysql        = "mysql"
	FormatSqlOracle       = "oracle"
	FormatSqlPostgresql   = "postgresql"
	FormatSqlSqlite3      = "sqlite3"
	FormatSqlSqlserver    = "sqlserver"
	FormatStaruml2        = "staruml2"
	FormatXlsx            = "xlsx"
)

func GetFileFormat(filename string) string {
	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".graphql":
		fallthrough
	case ".graphqls":
		return FormatGraphql
	case ".mdj":
		return FormatStaruml2
	case ".ojson":
		return FormatOctopus
	case ".plantuml":
		return FormatPlantuml
	case ".schema":
		return FormatSchemaConverter
	case ".xlsx":
		return FormatXlsx
	default:
		return ""
	}
}

func GetFileFormatIfNotSet(fileFormat string, filename string) string {
	if fileFormat != "" {
		return fileFormat
	} else {
		return GetFileFormat(filename)
	}
}
