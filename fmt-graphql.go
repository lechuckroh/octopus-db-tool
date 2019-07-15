package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type GraphqlClass struct {
	table    *Table
	Name     string
	Fields   []*GraphqlField
	PKFields []*GraphqlField
}

type GraphqlField struct {
	Column *Column
	Name   string
	Type   string
}

type Graphql struct {
}

func NewGraphqlClass(table *Table, output *GenOutput) *GraphqlClass {
	className := table.ClassName
	if className == "" {
		tableName := table.Name
		for _, prefix := range output.PrefixesToRemove {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		className = strcase.ToCamel(tableName)
	}

	fields := make([]*GraphqlField, 0)
	pkFields := make([]*GraphqlField, 0)
	for _, column := range table.Columns {
		field := NewGraphqlField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
	}

	if len(pkFields) == 1 {
		pkFields[0].Type = "ID!"
	}

	return &GraphqlClass{
		table:  table,
		Name:   className,
		Fields: fields,
	}
}

func NewGraphqlField(column *Column) *GraphqlField {
	var fieldType string
	nullable := column.Nullable
	columnType := strings.ToLower(column.Type)
	switch columnType {
	case "datetime":
		fallthrough
	case "date":
		fallthrough
	case "varchar":
		fallthrough
	case "char":
		fallthrough
	case "string":
		fallthrough
	case "text":
		fieldType = "String"
	case "bool":
		fallthrough
	case "boolean":
		fieldType = "Boolean"
	case "bigint":
		fallthrough
	case "long":
		fallthrough
	case "int":
		fallthrough
	case "integer":
		fallthrough
	case "smallint":
		fieldType = "Int"
	case "float":
		fallthrough
	case "number":
		fallthrough
	case "double":
		fieldType = "Float"
	default:
		fieldType = "String"
	}
	if !nullable {
		fieldType = fieldType + "!"
	}

	return &GraphqlField{
		Column: column,
		Name:   strcase.ToLowerCamel(column.Name),
		Type:   fieldType,
	}
}

func (l *Graphql) Generate(schema *Schema, output *GenOutput) error {
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

	appendLine(0, `schema {
    query: Query
}
`)

	classes := make([]*GraphqlClass, 0)

	appendLine(0, "type Query {")
	for _, table := range schema.Tables {
		class := NewGraphqlClass(table, output)
		classes = append(classes, class)

		appendLine(1, fmt.Sprintf("%sList: [%s]", strcase.ToLowerCamel(class.Name), class.Name))
	}
	appendLine(0, "}")
	appendLine(0, "")

	for _, class := range classes {
		appendLine(0, fmt.Sprintf("type %s {", class.Name))

		for _, field := range class.Fields {
			appendLine(1, fmt.Sprintf("%s: %s", field.Name, field.Type))
		}
		appendLine(0, "}")
		appendLine(0, "")
	}

	// Write file
	outputFile := path.Join(output.Directory,
		fmt.Sprintf("%s-%s.graphqls", schema.Name, schema.Version))

	if err := ioutil.WriteFile(outputFile, []byte(strings.Join(contents, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", outputFile)

	return nil
}

func (l *Graphql) getType(column *Column) string {
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
