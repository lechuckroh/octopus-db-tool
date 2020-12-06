package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type MysqlExportOption struct {
	TableFilter      TableFilterFn
	UniqueNameSuffix string
}

type MysqlExport struct {
	schema *Schema
}

type MysqlExportTmplData struct {
	Name      string
	Columns   []string
	PK        string
	UniqueKey string
}

func (c *MysqlExport) Export(wr io.Writer, option *MysqlExportOption) error {
	tmplText := `{{"" -}}
CREATE TABLE IF NOT EXISTS {{.Name}} (
{{range .Columns}}  {{.}},
{{end}}
{{- if ne .PK ""}}  {{.PK}}{{if ne .UniqueKey ""}},{{end}}
{{- end}}
{{if ne .UniqueKey ""}}  {{.UniqueKey}}
{{- end}}
);
`
	funcMap := template.FuncMap{}
	tmpl, err := NewTemplate("mysqlDDL", tmplText, funcMap)
	if err != nil {
		return err
	}

	for _, table := range c.schema.Tables {
		if err := c.exportTable(wr, option, tmpl, table); err != nil {
			return err
		}
	}

	return nil
}

// exportTable exports octopus table to mysql DDL
func (c *MysqlExport) exportTable(
	wr io.Writer,
	option *MysqlExportOption,
	tmpl *template.Template,
	table *Table,
) error {
	columns := make([]string, 0)
	pkColumns := make([]string, 0)
	uniqueColumns := make([]string, 0)
	for _, column := range table.Columns {
		params := make([]string, 0)
		params = append(params, column.Name)
		params = append(params, c.toMysqlColumnType(column))
		constraints := c.columnConstraints(column)
		if constraints != "" {
			params = append(params, constraints)
		}
		columns = append(columns, strings.Join(params, " "))

		if column.PrimaryKey {
			pkColumns = append(pkColumns, c.quote(column.Name))
		}
		if column.UniqueKey {
			uniqueColumns = append(uniqueColumns, c.quote(column.Name))
		}
	}

	data := MysqlExportTmplData{
		Name:    table.Name,
		Columns: columns,
	}
	if len(pkColumns) > 0 {
		data.PK = fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", "))
	}
	if len(uniqueColumns) > 0 {
		data.UniqueKey = fmt.Sprintf("UNIQUE KEY %s (%s)",
			c.quote(table.Name+option.UniqueNameSuffix), strings.Join(uniqueColumns, ", "))
	}

	return tmpl.Execute(wr, data)
}

func (c *MysqlExport) quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (c *MysqlExport) toMysqlColumnType(col *Column) string {
	switch col.Type {
	case ColTypeString:
		return fmt.Sprintf("varchar(%d)", col.Size)
	case ColTypeText:
		return "text"
	case ColTypeBoolean:
		return "bit(1)"
	case ColTypeLong:
		return "bigint"
	case ColTypeInt:
		return "int"
	case ColTypeDecimal:
		fallthrough
	case ColTypeFloat:
		fallthrough
	case ColTypeDouble:
		if col.Size > 0 {
			if col.Scale > 0 {
				return fmt.Sprintf("decimal(%d, %d)", col.Size, col.Scale)
			} else {
				return fmt.Sprintf("decimal(%d)", col.Size)
			}
		} else {
			return col.Type
		}
	case ColTypeDateTime:
		return "datetime"
	case ColTypeDate:
		return "date"
	case ColTypeTime:
		return "time"
	case ColTypeBlob:
		return "blob"
	default:
		return col.Type
	}
}

func (c *MysqlExport) columnConstraints(column *Column) string {
	constraints := make([]string, 0)

	if !column.Nullable {
		constraints = append(constraints, "NOT NULL")
	}

	if column.AutoIncremental {
		constraints = append(constraints, "AUTO_INCREMENT")
	}

	if column.DefaultValue != "" {
		defaultValue := column.DefaultValue
		if IsStringType(column.Type) {
			defaultValue = fmt.Sprintf("'%s'", defaultValue)
		}
		constraints = append(constraints, "DEFAULT "+defaultValue)
	}

	if column.Description != "" {
		constraints = append(constraints, fmt.Sprintf("COMMENT '%s'", column.Description))
	}

	return strings.Join(constraints, " ")
}
