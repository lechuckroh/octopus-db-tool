package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type QuickDBD struct {

}

func (f *QuickDBD) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return f.FromString(data)
}

func (f *QuickDBD) FromString(data []byte) error {
	// TODO
	return errors.New("not implemented")
}

func (f *QuickDBD) ToFile(schema *Schema, filename string) error {
	data, err := f.ToString(schema)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (f *QuickDBD) ToString(schema *Schema) ([]byte, error) {
	result := make([]string, 0)
	for _, table := range schema.Tables {
		result = append(result, f.getTableDef(table))
		result = append(result, strings.Repeat("-", len(table.Name)))

		for _, column := range table.Columns {
			result = append(result, f.getColumnDef(column))
		}
		result = append(result, "")
	}

	return []byte(strings.Join(result, "\n")), nil
}

func (f *QuickDBD) getTableDef(table *Table) string {
	if table.Description != "" {
		return fmt.Sprintf("%s # %s", table.Name, table.Description)
	} else {
		return table.Name
	}
}

func (f *QuickDBD) getColumnDef(col *Column) string {
	params := make([]string, 0)
	params = append(params, col.Type)

	if col.PrimaryKey {
		params = append(params, "PK")
	}
	if col.UniqueKey {
		params = append(params, "UNIQUE")
	}
	if col.AutoIncremental {
		params = append(params, "AUTOINCREMENT")
	}
	if col.Nullable {
		params = append(params, "NULLABLE")
	}
	if col.DefaultValue != "" {
		params = append(params, fmt.Sprintf("default=%s", col.DefaultValue))
	}
	if col.Ref != nil {
		ref := col.Ref
		params = append(params, fmt.Sprintf("FK >- %s.%s", ref.Table, ref.Column))
	}
	if col.Description != "" {
		params = append(params, fmt.Sprintf("# %s", col.Description))
	}

	return fmt.Sprintf("%s %s", col.Name, strings.Join(params, " "))
}