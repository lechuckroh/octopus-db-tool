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


func (l *Liquibase) Generate(schema *Schema, output *GenOutput) error {
	// Create directory
	if err := os.MkdirAll(output.Directory, 0777); err != nil {
		return err
	}
	log.Printf("[MKDIR] %s", output.Directory)

	indent := "  "
	contents := make([]string, 0)
	appendLine := func(indentLevel int, line string) {
		contents = append(contents,
			strings.Repeat(indent, indentLevel)+line,
		)
	}

	appendLine(0, "databaseChangeLog:")
	i := 1
	for _, table := range schema.Tables {
		uniques := make([]string, 0)
		primaryKeys := make([]string, 0)
		for _, column := range table.Columns {
			if column.UniqueKey {
				uniques = append(uniques, column.Name)
			}
			if column.PrimaryKey {
				primaryKeys = append(primaryKeys, column.Name)
			}
		}
		pkCount := len(primaryKeys)
		uniqueCount := len(uniques)

		appendLine(0, "- changeSet:")
		appendLine(2, fmt.Sprintf("id: %d", i))
		appendLine(2, "author: "+schema.Author)
		appendLine(2, "changes:")
		appendLine(2, "- createTable:")
		appendLine(4, "tableName: "+table.Name)
		appendLine(4, "columns:")
		for _, column := range table.Columns {
			appendLine(4, "- column:")
			appendLine(6, "name: "+column.Name)
			appendLine(6, "type: "+l.getType(column))

			// auto_incremental
			if column.AutoIncremental {
				appendLine(6, "autoIncrement: true")
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
				appendLine(6, "constraints:")
				for _, constraint := range constraints {
					appendLine(7, constraint)
				}
			}

			// default value
			if column.DefaultValue != "" {
				if IsStringType(column.Type) {
					appendLine(6, "defaultValue: \""+column.DefaultValue+"\"")
				} else if IsBooleanType(column.Type) {
					appendLine(6, "defaultValueBoolean: "+column.DefaultValue)
				} else if IsNumericType(column.Type) {
					appendLine(6, "defaultValueNumeric: "+column.DefaultValue)
				}
			}
		}
		i++

		// Primary Key
		if pkCount >= 2 {
			appendLine(0, "- changeSet:")
			appendLine(2, fmt.Sprintf("id: %d", i))
			appendLine(2, "author: "+schema.Author)
			appendLine(2, "changes:")
			appendLine(2, "- addPrimaryKey:")
			appendLine(4, "columnNames: "+strings.Join(primaryKeys, ", "))
			appendLine(4, "tableName: "+table.Name)

			i++
		}
		// Unique Constraint
		if uniqueCount >= 2 {
			appendLine(0, "- changeSet:")
			appendLine(2, fmt.Sprintf("id: %d", i))
			appendLine(2, "author: "+schema.Author)
			appendLine(2, "changes:")
			appendLine(2, "- addUniqueConstraint:")
			appendLine(4, "columnNames: "+strings.Join(uniques, ", "))
			appendLine(4, "constraintName: "+table.Name+output.UniqueNameSuffix)
			appendLine(4, "tableName: "+table.Name)

			i++
		}
	}

	// Write file
	outputFile := path.Join(output.Directory,
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
	case "string":
		fallthrough
	case "varchar":
		typ = "varchar"
	case "char":
		typ = "char"
	case "text":
		typ = "clob"
	case "bool":
		fallthrough
	case "boolean":
		typ = "boolean"
	case "bigint":
		fallthrough
	case "long":
		typ = "bigint"
	case "int":
		fallthrough
	case "integer":
		typ = "int"
	case "smallint":
		typ = "smallint"
	case "float":
		typ = "float"
	case "number":
		fallthrough
	case "double":
		typ = "double"
	case "datetime":
		typ = "datetime"
	case "date":
		typ = "date"
	case "blob":
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
