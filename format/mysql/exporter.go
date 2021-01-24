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

type exportTplData struct {
	Name        string
	Definitions []string
}

func (c *Exporter) Export(wr io.Writer) error {
	tplText := `{{"" -}}
CREATE TABLE IF NOT EXISTS {{.Name}} (
{{range .Definitions}}  {{.}}
{{end}});
`
	funcMap := template.FuncMap{}
	tpl, err := util.NewTemplate("mysqlDDL", tplText, funcMap)
	if err != nil {
		return err
	}

	for _, table := range c.schema.Tables {
		if c.option.TableFilter != nil && !c.option.TableFilter(table) {
			continue
		}
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
	var definitions []string
	var pkColumns []string
	var uniqueColumns []string
	for _, column := range table.Columns {
		var params []string
		params = append(params, column.Name)
		params = append(params, c.toMysqlColumnType(column))
		constraints := c.columnConstraints(column)
		if constraints != "" {
			params = append(params, constraints)
		}
		definitions = append(definitions, strings.Join(params, " "))

		if column.PrimaryKey {
			pkColumns = append(pkColumns, c.quote(column.Name))
		}
		if column.UniqueKey {
			uniqueColumns = append(uniqueColumns, c.quote(column.Name))
		}
	}

	if len(pkColumns) > 0 {
		definitions = append(definitions, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", ")))
	}
	if len(uniqueColumns) > 0 {
		definitions = append(definitions, fmt.Sprintf("UNIQUE KEY %s (%s)",
			c.quote(table.Name+c.option.UniqueNameSuffix), strings.Join(uniqueColumns, ", ")))
	}

	for _, index := range table.Indices {
		definitions = append(definitions, fmt.Sprintf("INDEX %s (%s)",
			c.quote(index.Name), util.QuoteAndJoin(index.Columns, "`", ", ")))
	}

	// append ',' except last
	for i, definition := range definitions {
		if i < len(definitions) -1 {
			definitions[i] = definition + ","
		}
	}

	data := exportTplData{
		Name:        table.Name,
		Definitions: definitions,
	}

	return tpl.Execute(wr, data)
}

func (c *Exporter) quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (c *Exporter) formatColumnType(colType string, col *octopus.Column) string {
	if col.Size > 0 {
		if col.Scale > 0 {
			return fmt.Sprintf("%s(%d,%d)", colType, col.Size, col.Scale)
		} else {
			return fmt.Sprintf("%s(%d)", colType, col.Size)
		}
	}
	return colType
}

func (c *Exporter) toMysqlColumnType(col *octopus.Column) string {
	switch col.Type {
	case octopus.ColTypeBinary:
		return "binary"
	case octopus.ColTypeBit:
		return c.formatColumnType("bit", col)
	case octopus.ColTypeBlob8:
		return "tinyblob"
	case octopus.ColTypeBlob16:
		return "blob"
	case octopus.ColTypeBlob24:
		return "mediumblob"
	case octopus.ColTypeBlob32:
		return "longblob"
	case octopus.ColTypeBoolean:
		return "bit(1)"
	case octopus.ColTypeChar:
		return c.formatColumnType("char", col)
	case octopus.ColTypeDate:
		return "date"
	case octopus.ColTypeDateTime:
		return "datetime"
	case octopus.ColTypeDecimal:
		return c.formatColumnType("decimal", col)
	case octopus.ColTypeDouble:
		return c.formatColumnType("double", col)
	case octopus.ColTypeEnum:
		return "enum(" + util.QuoteAndJoin(col.Values, "'", ", ") + ")"
	case octopus.ColTypeFloat:
		return c.formatColumnType("float", col)
	case octopus.ColTypeGeometry:
		return "geometry"
	case octopus.ColTypeInt8:
		return c.formatColumnType("tinyint", col)
	case octopus.ColTypeInt16:
		return c.formatColumnType("smallint", col)
	case octopus.ColTypeInt24:
		return c.formatColumnType("mediumint", col)
	case octopus.ColTypeInt32:
		return c.formatColumnType("int", col)
	case octopus.ColTypeInt64:
		return c.formatColumnType("bigint", col)
	case octopus.ColTypeJSON:
		return "json"
	case octopus.ColTypePoint:
		return "point"
	case octopus.ColTypeSet:
		return "set(" + util.QuoteAndJoin(col.Values, "'", ", ") + ")"
	case octopus.ColTypeText8:
		return "tinytext"
	case octopus.ColTypeText16:
		return "text"
	case octopus.ColTypeText24:
		return "mediumtext"
	case octopus.ColTypeText32:
		return "longtext"
	case octopus.ColTypeTime:
		return "time"
	case octopus.ColTypeVarbinary:
		return "varbinary"
	case octopus.ColTypeVarchar:
		return c.formatColumnType("varchar", col)
	case octopus.ColTypeYear:
		return "year"
	default:
		return col.Type
	}
}

func (c *Exporter) columnConstraints(column *octopus.Column) string {
	var constraints []string

	if column.NotNull {
		constraints = append(constraints, "NOT NULL")
	}

	if column.AutoIncremental {
		constraints = append(constraints, "AUTO_INCREMENT")
	}

	if column.DefaultValue != "" {
		defaultValue, fn := column.GetDefaultValue()
		if fn {
			defaultValue = defaultValue + "()"
		} else {
			if util.IsStringType(column.Type) {
				defaultValue = fmt.Sprintf("'%s'", defaultValue)
			}
		}
		constraints = append(constraints, "DEFAULT "+defaultValue)
	}

	if column.OnUpdate != "" {
		onUpdate, fn := column.GetOnUpdate()
		if fn {
			onUpdate = onUpdate + "()"
		} else {
			if util.IsStringType(column.Type) {
				onUpdate = fmt.Sprintf("'%s'", onUpdate)
			}
		}
		constraints = append(constraints, "ON UPDATE "+onUpdate)
	}

	if column.Description != "" {
		constraints = append(constraints, fmt.Sprintf("COMMENT '%s'", column.Description))
	}

	return strings.Join(constraints, " ")
}
