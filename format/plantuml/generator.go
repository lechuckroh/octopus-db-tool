package plantuml

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"io"
	"strings"
)

type Option struct {
}

type Generator struct {
	schema *octopus.Schema
	option *Option
}

func (c *Generator) Generate(wr io.Writer) error {
	result := make([]string, 0)
	refs := make([]string, 0)
	for _, table := range c.schema.Tables {
		result = append(result, fmt.Sprintf("entity %s {", table.Name))

		separatorAdded := false
		for _, column := range table.Columns {
			if !column.PrimaryKey && !separatorAdded {
				separatorAdded = true
				result = append(result, "    --")
			}
			result = append(result, fmt.Sprintf("    %s", c.getColumnDef(column)))

			// Reference
			if column.Ref != nil {
				refs = append(refs, fmt.Sprintf("%s }o-|| %s", table.Name, column.Ref.Table))
			}
		}
		result = append(result, "}")
	}

	result = append(result, refs...)

	_, err := wr.Write([]byte(strings.Join(result, "\n")))
	return err
}

func (c *Generator) getTableDef(table *octopus.Table) string {
	if table.Description != "" {
		return fmt.Sprintf("%s # %s", table.Name, table.Description)
	} else {
		return table.Name
	}
}

func (c *Generator) getColumnDef(col *octopus.Column) string {
	line := ""

	if col.NotNull {
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
