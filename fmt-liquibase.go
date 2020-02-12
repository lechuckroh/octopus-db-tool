package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Liquibase struct {
}

func (l *Liquibase) Generate(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
) error {
	// Create directory
	if err := os.MkdirAll(output.FilePath, 0777); err != nil {
		return err
	}
	log.Printf("[MKDIR] %s", output.FilePath)

	indent := "  "
	contents := make([]string, 0)
	appendLine := func(indentLevel int, line string) {
		contents = append(contents,
			strings.Repeat(indent, indentLevel)+line,
		)
	}

	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)

	appendLine(0, "databaseChangeLog:")
	appendLine(1, "- objectQuotingStrategy: QUOTE_ALL_OBJECTS")
	i := 1
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}

		uniques := table.GetUniqueColumnNames()
		primaryKeys := table.GetPrimaryKeyColumnNames()
		pkCount := len(primaryKeys)
		uniqueCount := len(uniques)

		appendLine(1, "- changeSet:")
		appendLine(3, fmt.Sprintf("id: %d", i))
		appendLine(3, "author: "+schema.Author)
		appendLine(3, "changes:")
		appendLine(4, "- createTable:")
		appendLine(6, "tableName: "+table.Name)
		appendLine(6, "columns:")
		for _, column := range table.Columns {
			appendLine(7, "- column:")
			appendLine(9, "name: "+column.Name)
			appendLine(9, "type: "+l.getType(column))

			// auto_incremental
			if column.AutoIncremental {
				appendLine(9, "autoIncrement: true")
			}

			// constraints
			constraints := make([]string, 0)
			if column.PrimaryKey && pkCount == 1 {
				constraints = append(constraints, "primaryKey: true")
			} else if !column.Nullable {
				constraints = append(constraints, "nullable: false")
			}
			if column.UniqueKey && uniqueCount == 1 {
				constraints = append(constraints, "unique: true")
			}

			if len(constraints) > 0 {
				appendLine(9, "constraints:")
				for _, constraint := range constraints {
					appendLine(10, constraint)
				}
			}

			// default value
			if column.DefaultValue != "" {
				if IsStringType(column.Type) {
					appendLine(9, "defaultValue: \""+column.DefaultValue+"\"")
				} else if IsBooleanType(column.Type) {
					appendLine(9, "defaultValueBoolean: "+column.DefaultValue)
				} else if IsNumericType(column.Type) {
					appendLine(9, "defaultValueNumeric: "+column.DefaultValue)
				}
			}
		}
		i++

		// Primary Key
		if pkCount >= 2 {
			appendLine(1, "- changeSet:")
			appendLine(3, fmt.Sprintf("id: %d", i))
			appendLine(3, "author: "+schema.Author)
			appendLine(3, "changes:")
			appendLine(4, "- addPrimaryKey:")
			appendLine(6, "columnNames: "+strings.Join(primaryKeys, ", "))
			appendLine(6, "tableName: "+table.Name)

			i++
		}
		// Unique Constraint
		if uniqueCount >= 1 {
			appendLine(1, "- changeSet:")
			appendLine(3, fmt.Sprintf("id: %d", i))
			appendLine(3, "author: "+schema.Author)

			if uniqueCount == 1 {
				appendLine(3, "preConditions:")
				appendLine(4, "onError: CONTINUE")
				appendLine(4, "onFail: CONTINUE")
				appendLine(4, "dbms:")
				appendLine(5, "type: derby, h2, mssql, mariadb, mysql, postgresql, sqlite")
			}

			appendLine(3, "changes:")
			appendLine(4, "- addUniqueConstraint:")
			appendLine(6, "columnNames: "+strings.Join(uniques, ", "))
			appendLine(6, "constraintName: "+table.Name+uniqueNameSuffix)
			appendLine(6, "tableName: "+table.Name)

			i++
		}
	}

	// Write file
	outputFile := path.Join(output.FilePath,
		fmt.Sprintf("%s-%s.yaml", schema.Name, schema.Version))

	if err := ioutil.WriteFile(outputFile, []byte(strings.Join(contents, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", outputFile)

	return nil
}

func (l *Liquibase) getType(column *Column) string {
	typ := ""
	switch strings.ToLower(column.Type) {
	case ColTypeString:
		typ = "varchar"
	case ColTypeText:
		typ = "clob"
	case ColTypeBoolean:
		typ = "boolean"
	case ColTypeLong:
		typ = "bigint"
	case ColTypeInt:
		typ = "int"
	case ColTypeDecimal:
		typ = "decimal"
	case ColTypeFloat:
		typ = "float"
	case ColTypeDouble:
		typ = "double"
	case ColTypeDateTime:
		typ = "datetime"
	case ColTypeDate:
		typ = "date"
	case ColTypeTime:
		typ = "time"
	case ColTypeBlob:
		typ = "blob"
	default:
		typ = column.Type
	}
	if column.Size > 0 {
		return fmt.Sprintf("%s(%d)", typ, column.Size)
	} else {
		return typ
	}
}
