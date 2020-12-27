package mysql

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"strings"
	"text/template"
)

type ExportOption struct {
	TableFilter      octopus.TableFilterFn
	UniqueNameSuffix string
}

type Exporter struct {
	schema *octopus.Schema
	option *ExportOption
}

func (c *Exporter) Export(wr io.Writer) error {
	tplText := `{{"" -}}
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
	tpl, err := util.NewTemplate("mysqlDDL", tplText, funcMap)
	if err != nil {
		return err
	}

	for _, table := range c.schema.Tables {
		if err := c.exportTable(wr, tpl, table); err != nil {
			return err
		}
	}

	return nil
}

// exportTable exports octopus table to mysql DDL
func (c *Exporter) exportTable(
	wr io.Writer,
	tpl *template.Template,
	table *octopus.Table,
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

	data := struct {
		Name      string
		Columns   []string
		PK        string
		UniqueKey string
	}{
		Name:    table.Name,
		Columns: columns,
	}
	if len(pkColumns) > 0 {
		data.PK = fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", "))
	}
	if len(uniqueColumns) > 0 {
		data.UniqueKey = fmt.Sprintf("UNIQUE KEY %s (%s)",
			c.quote(table.Name+c.option.UniqueNameSuffix), strings.Join(uniqueColumns, ", "))
	}

	return tpl.Execute(wr, data)
}

func (c *Exporter) quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (c *Exporter) toMysqlColumnType(col *octopus.Column) string {
	switch col.Type {
	case octopus.ColTypeChar:
		return fmt.Sprintf("char(%d)", col.Size)
	case octopus.ColTypeVarchar:
		return fmt.Sprintf("varchar(%d)", col.Size)
	case octopus.ColTypeText8:
		return "tinytext"
	case octopus.ColTypeText16:
		return "text"
	case octopus.ColTypeText24:
		return "mediumtext"
	case octopus.ColTypeText32:
		return "longtext"
	case octopus.ColTypeBoolean:
		return "bit(1)"
	case octopus.ColTypeInt8:
		return "tinyint"
	case octopus.ColTypeInt16:
		return "smallint"
	case octopus.ColTypeInt24:
		return "mediumint"
	case octopus.ColTypeInt32:
		return "tinyint"
	case octopus.ColTypeInt64:
		return "bigint"
	case octopus.ColTypeDecimal:
		if col.Size > 0 {
			if col.Scale > 0 {
				return fmt.Sprintf("decimal(%d, %d)", col.Size, col.Scale)
			} else {
				return fmt.Sprintf("decimal(%d)", col.Size)
			}
		} else {
			return col.Type
		}
	case octopus.ColTypeFloat:
		return "float"
	case octopus.ColTypeDouble:
		return "double"
	case octopus.ColTypeDateTime:
		return "datetime"
	case octopus.ColTypeDate:
		return "date"
	case octopus.ColTypeTime:
		return "time"
	case octopus.ColTypeBlob8:
		return "tinyblob"
	case octopus.ColTypeBlob16:
		return "blob"
	case octopus.ColTypeBlob24:
		return "mediumblob"
	case octopus.ColTypeBlob32:
		return "longblob"
	default:
		return col.Type
	}
}

func (c *Exporter) columnConstraints(column *octopus.Column) string {
	constraints := make([]string, 0)

	if !column.Nullable {
		constraints = append(constraints, "NOT NULL")
	}

	if column.AutoIncremental {
		constraints = append(constraints, "AUTO_INCREMENT")
	}

	if column.DefaultValue != "" {
		defaultValue := column.DefaultValue
		if util.IsStringType(column.Type) {
			defaultValue = fmt.Sprintf("'%s'", defaultValue)
		}
		constraints = append(constraints, "DEFAULT "+defaultValue)
	}

	if column.Description != "" {
		constraints = append(constraints, fmt.Sprintf("COMMENT '%s'", column.Description))
	}

	return strings.Join(constraints, " ")
}
