package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type PlantUML struct {

}

func (f *PlantUML) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return f.FromString(data)
}

func (f *PlantUML) FromString(data []byte) error {
	// TODO
	return errors.New("not implemented")
}

func (f *PlantUML) ToSchema() (*Schema, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (f *PlantUML) ToFile(schema *Schema, filename string) error {
	data, err := f.ToString(schema)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (f *PlantUML) ToString(schema *Schema) ([]byte, error) {
	result := make([]string, 0)
	refs := make([]string, 0)
	for _, table := range schema.Tables {
		result = append(result, fmt.Sprintf("entity %s {", table.Name))

		separatorAdded := false
		for _, column := range table.Columns {
			if !column.PrimaryKey && !separatorAdded{
				separatorAdded = true
				result = append(result, "    --")
			}
			result = append(result, fmt.Sprintf("    %s", f.getColumnDef(column)))

			// Reference
			if column.Ref != nil {
				refs = append(refs, fmt.Sprintf("%s }o-|| %s", table.Name, column.Ref.Table))
			}
		}
		result = append(result, "}")
	}

	result = append(result, refs...)

	return []byte(strings.Join(result, "\n")), nil
}

func (f *PlantUML) getTableDef(table *Table) string {
	if table.Description != "" {
		return fmt.Sprintf("%s # %s", table.Name, table.Description)
	} else {
		return table.Name
	}
}

func (f *PlantUML) getColumnDef(col *Column) string {
	line := ""

	if !col.Nullable {
		line += "* "
	}

	line += col.Name + ": " + col.Type

	if col.PrimaryKey {
		line += " <<PK>>"
	}
	if col.UniqueKey {
		line += " <<UQ>>"
	}
	if col.Ref != nil {
		line += " <<FK>>"
	}
	return line
}
