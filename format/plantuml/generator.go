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
			result = append(result, fmt.Sprintf("    %s", getColumnDef(column)))

			// Reference
			if ref := column.Ref; ref != nil {
				relationship := getRelationshipType(ref)
				refs = append(refs, fmt.Sprintf("%s %s %s", table.Name, relationship, ref.Table))
			}
		}
		result = append(result, "}")
	}

	result = append(result, refs...)

	_, err := wr.Write([]byte(strings.Join(result, "\n")))
	return err
}

func getTableDef(table *octopus.Table) string {
	if table.Description != "" {
		return fmt.Sprintf("%s # %s", table.Name, table.Description)
	} else {
		return table.Name
	}
}

func getColumnDef(col *octopus.Column) string {
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

func getRelationshipType(ref *octopus.Reference) string {
	switch ref.Relationship {
	case octopus.RefManyToOne:
		return "}o-||"
	case octopus.RefOneToMany:
		return "||-o{"
	case octopus.RefOneToOne:
		return "||-||"
	default:
		return "}o-||"
	}
}
