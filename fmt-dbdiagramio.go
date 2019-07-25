package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type DBDiagramIO struct {
}

func (f *DBDiagramIO) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return f.FromString(data)
}

func (f *DBDiagramIO) FromString(data []byte) error {
	// TODO
	return errors.New("not implemented")
}

func (f *DBDiagramIO) ToSchema() (*Schema, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (f *DBDiagramIO) ToFile(schema *Schema, filename string) error {
	data, err := f.ToString(schema)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (f *DBDiagramIO) ToString(schema *Schema) ([]byte, error) {
	result := make([]string, 0)
	definedTables := make(map[string]bool)
	deferredRefs := make([]string, 0)

	for _, table := range schema.Tables {
		result = append(result, fmt.Sprintf("Table %s {", table.Name))

		for _, column := range table.Columns {
			params := make([]string, 0)

			if column.PrimaryKey {
				params = append(params, "pk")
			}
			if column.UniqueKey {
				params = append(params, "unique")
			}
			//if column.AutoIncremental {
			//	params = append(params, "auto increment")
			//}
			if !column.Nullable {
				params = append(params, "not null")
			}
			if column.DefaultValue != "" {
				defaultValue := column.DefaultValue
				if IsStringType(column.Type) {
					defaultValue = fmt.Sprintf("'%s'", defaultValue)
				}
				params = append(params, fmt.Sprintf("default: %s", defaultValue))
			}
			if column.Ref != nil {
				ref := column.Ref
				if definedTables[ref.Table] {
					params = append(params,
						fmt.Sprintf("ref: > %s.%s", ref.Table, ref.Column))
				} else {
					deferredRefs = append(deferredRefs,
						fmt.Sprintf("Ref: %s.%s > %s.%s",
							table.Name, column.Name, ref.Table, ref.Column),
					)
				}
			}
			if column.Description != "" {
				params = append(params, fmt.Sprintf("note: \"%s\"", column.Description))
			}

			columnType := column.Type
			if column.Size > 0 {
				if column.Scale > 0 {
					columnType = fmt.Sprintf("%s(%d,%d)", column.Type, column.Size, column.Scale)
				} else {
					columnType = fmt.Sprintf("%s(%d)", column.Type, column.Size)
				}
			}

			if len(params) == 0 {
				result = append(result, fmt.Sprintf("  %s %s", column.Name, columnType))
			} else {
				result = append(result, fmt.Sprintf("  %s %s [%s]", column.Name, columnType, strings.Join(params, ", ")))
			}
		}
		result = append(result, "}")
		result = append(result, "")

		definedTables[table.Name] = true
	}

	result = append(result, strings.Join(deferredRefs, "\n"))

	return []byte(strings.Join(result, "\n")), nil
}
