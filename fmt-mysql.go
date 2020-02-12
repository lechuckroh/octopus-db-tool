package main

import (
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
	"io"
	"io/ioutil"
	"strings"
)

type Mysql struct {
	schema *Schema
}

func (m *Mysql) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return m.FromString(data)
}

func (m *Mysql) FromString(data []byte) error {
	tokens := sqlparser.NewStringTokenizer(string(data))

	tables := make([]*Table, 0)
	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if err != nil {
			m.schema = nil
			return err
		}

		switch stmt.(type) {
		case *sqlparser.DDL:
			ddl := stmt.(*sqlparser.DDL)
			columns := make([]*Column, 0)

			tableSpec := ddl.TableSpec

			if tableSpec != nil {
				pkSet := NewStringSet()
				uniqueSet := NewStringSet()
				for _, idx := range tableSpec.Indexes {
					info := idx.Info
					if info.Primary {
						for _, c := range idx.Columns {
							pkSet.Add(c.Column.String())
						}
					} else if info.Unique {
						for _, c := range idx.Columns {
							uniqueSet.Add(c.Column.String())
						}
					}
				}

				for _, col := range ddl.TableSpec.Columns {
					name := col.Name.String()
					nullable := !bool(col.Type.NotNull)
					defaultValue := SQLValToString(col.Type.Default, "")
					if nullable && defaultValue == "null" {
						defaultValue = ""
					}
					comment := SQLValToString(col.Type.Comment, "")
					columns = append(columns, &Column{
						Name:            name,
						Type:            m.fromColumnType(col.Type),
						Description:     comment,
						Size:            uint16(SQLValToInt(col.Type.Length, 0)),
						Scale:           uint16(SQLValToInt(col.Type.Scale, 0)),
						Nullable:        nullable,
						PrimaryKey:      pkSet.Contains(name),
						UniqueKey:       uniqueSet.Contains(name),
						AutoIncremental: bool(col.Type.Autoincrement),
						DefaultValue:    defaultValue,
					})
				}
				tables = append(tables, &Table{
					Name:    ddl.NewName.Name.String(),
					Columns: columns,
				})
			}
		}
	}

	m.schema = &Schema{
		Tables: tables,
	}

	return nil
}

func (m *Mysql) ToSchema() (*Schema, error) {
	if m.schema == nil {
		return nil, errors.New("schema is not read")
	}
	return m.schema, nil
}

func (m *Mysql) ToFile(schema *Schema, filename string) error {
	data, err := m.ToString(schema)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (m *Mysql) quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (m *Mysql) ToString(schema *Schema) ([]byte, error) {
	result := make([]string, 0)

	indent := "  "
	for _, table := range schema.Tables {
		lines := make([]string, 0)

		primaryKeys := make([]string, 0)
		uniqueKeys := make([]string, 0)
		for _, column := range table.Columns {
			params := make([]string, 0)

			if column.PrimaryKey {
				primaryKeys = append(primaryKeys, m.quote(column.Name))
			}
			if column.UniqueKey {
				uniqueKeys = append(uniqueKeys, m.quote(column.Name))
			}

			if !column.Nullable {
				params = append(params, "NOT NULL")
			}

			if column.AutoIncremental {
				params = append(params, "AUTO_INCREMENT")
			}

			if column.DefaultValue != "" {
				defaultValue := column.DefaultValue
				if IsStringType(column.Type) {
					defaultValue = fmt.Sprintf("'%s'", defaultValue)
				}
				params = append(params, "DEFAULT "+defaultValue)
			}

			if column.Description != "" {
				params = append(params, fmt.Sprintf("COMMENT '%s'", column.Description))
			}

			lines = append(lines,
				fmt.Sprintf(indent+"%s %s %s",
					m.quote(column.Name),
					m.toMysqlColumnType(column),
					strings.Join(params, " ")))
		}

		if len(primaryKeys) > 0 {
			lines = append(lines,
				fmt.Sprintf(indent+"PRIMARY KEY (%s)", strings.Join(primaryKeys, ", ")))
		}
		if len(uniqueKeys) > 0 {
			lines = append(lines,
				fmt.Sprintf(indent+"UNIQUE KEY %s (%s)",
					m.quote(table.Name+"_UNIQUE"),
					strings.Join(uniqueKeys, ", ")))
		}
		body := strings.Join(lines, ",\n")

		tableDef := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);", m.quote(table.Name), body)
		result = append(result, tableDef)
	}

	return []byte(strings.Join(result, "\n")), nil
}

func (m *Mysql) toMysqlColumnType(col *Column) string {
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

func (m *Mysql) fromColumnType(colType sqlparser.ColumnType) string {
	switch colType.Type {
	case "varchar":
		return ColTypeString
	case "longtext":
		fallthrough
	case "text":
		return ColTypeText
	case "bit":
		return ColTypeBoolean
	case "bigint":
		return ColTypeLong
	case "int":
		return ColTypeInt
	case "decimal":
		return ColTypeDecimal
	case "float":
		return ColTypeFloat
	case "double":
		return ColTypeDouble
	case "datetime":
		return ColTypeDateTime
	case "date":
		return ColTypeDate
	case "time":
		return ColTypeTime
	case "blob":
		return ColTypeBlob
	default:
		return colType.Type
	}
}
