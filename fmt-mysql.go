package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Mysql struct {
}

func (m *Mysql) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return m.FromString(data)
}

func (m *Mysql) FromString(data []byte) error {
	// TODO
	return errors.New("not implemented")
}

func (m *Mysql) ToSchema() (*Schema, error) {
	// TODO
	return nil, errors.New("not implemented")
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
				fmt.Sprintf(indent + "%s %s %s",
					m.quote(column.Name),
					m.toMysqlColumnType(column),
					strings.Join(params, " ")))
		}

		if len(primaryKeys) > 0 {
			lines = append(lines,
				fmt.Sprintf(indent + "PRIMARY KEY (%s)", strings.Join(primaryKeys, ", ")))
		}
		if len(uniqueKeys) > 0 {
			lines = append(lines,
				fmt.Sprintf(indent + "UNIQUE KEY %s (%s)",
					m.quote(table.Name + "_UNIQUE"),
					strings.Join(uniqueKeys, ", ")))
		}
		body := strings.Join(lines, ",\n")

		tableDef := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);", m.quote(table.Name), body)
		result = append(result, tableDef)
	}

	return []byte(strings.Join(result, "\n")), nil
}

func (m* Mysql) toMysqlColumnType(col *Column) string {
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
	case ColTypeFloat:
		if col.Size > 0 {
			if col.Scale > 0 {
				return fmt.Sprintf("decimal(%d, %d)", col.Size, col.Scale)
			} else {
				return fmt.Sprintf("decimal(%d)", col.Size)
			}
		} else {
			return "float"
		}
	case ColTypeDouble:
		if col.Size > 0 {
			if col.Scale > 0 {
				return fmt.Sprintf("decimal(%d, %d)", col.Size, col.Scale)
			} else {
				return fmt.Sprintf("decimal(%d)", col.Size)
			}
		} else {
			return "double"
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