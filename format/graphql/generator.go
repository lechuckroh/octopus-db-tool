package graphql

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Class struct {
	table    *octopus.Table
	Name     string
	Fields   []*Field
	PKFields []*Field
}

type Field struct {
	Column *octopus.Column
	Name   string
	Type   string
}

type Option struct {
	TableFilter    octopus.TableFilterFn
	PrefixMapper   *common.PrefixMapper
	RemovePrefixes []string
}

type Generator struct {
	schema *octopus.Schema
	option *Option
}

func NewClass(
	table *octopus.Table,
	option *Option,
) *Class {
	className := table.ClassName
	if className == "" {
		tableName := table.Name
		prefixes := option.RemovePrefixes
		for _, prefix := range prefixes {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		className = strcase.ToCamel(tableName)

		if prefix := option.PrefixMapper.GetPrefix(table.Group); prefix != "" {
			className = prefix + className
		}
	}

	fields := make([]*Field, 0)
	pkFields := make([]*Field, 0)
	for _, column := range table.Columns {
		field := NewField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
	}

	if len(pkFields) == 1 {
		pkFields[0].Type = "ID!"
	}

	return &Class{
		table:  table,
		Name:   className,
		Fields: fields,
	}
}

func NewField(column *octopus.Column) *Field {
	var fieldType string
	nullable := column.Nullable
	columnType := strings.ToLower(column.Type)
	switch columnType {
	case octopus.ColTypeDateTime:
		fallthrough
	case octopus.ColTypeDate:
		fallthrough
	case octopus.ColTypeTime:
		fallthrough
	case octopus.ColTypeChar:
		fallthrough
	case octopus.ColTypeVarchar:
		fieldType = "String"
	case octopus.ColTypeBoolean:
		fieldType = "Boolean"
	case octopus.ColTypeInt8:
		fallthrough
	case octopus.ColTypeInt16:
		fallthrough
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		fallthrough
	case octopus.ColTypeInt64:
		fieldType = "Int"
	case octopus.ColTypeDecimal:
		fallthrough
	case octopus.ColTypeFloat:
		fallthrough
	case octopus.ColTypeDouble:
		fieldType = "Float"
	default:
		log.Printf("unknown column type: '%s', column: %s", column.Type, column.Name)
		fieldType = "String"
	}
	if !nullable {
		fieldType = fieldType + "!"
	}

	return &Field{
		Column: column,
		Name:   strcase.ToLowerCamel(column.Name),
		Type:   fieldType,
	}
}

func (c *Generator) Generate(outputPath string) error {
	// TODO: generate with template

	// Create directory
	if err := os.MkdirAll(outputPath, 0777); err != nil {
		return err
	}
	log.Printf("[MKDIR] %s", outputPath)

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

	classes := make([]*Class, 0)

	client := pluralize.NewClient()

	appendLine(0, "type Query {")
	for _, table := range c.schema.Tables {
		// filter table
		if c.option.TableFilter != nil && !c.option.TableFilter(table) {
			continue
		}

		class := NewClass(table, c.option)
		classes = append(classes, class)

		lowerClassName := strcase.ToLowerCamel(class.Name)
		appendLine(1, fmt.Sprintf("%s: [%s]", client.Plural(lowerClassName), class.Name))
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
	outputFile := path.Join(outputPath,
		fmt.Sprintf("%s-%s.graphqls", c.schema.Name, c.schema.Version))

	if err := ioutil.WriteFile(outputFile, []byte(strings.Join(contents, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", outputFile)

	return nil
}
