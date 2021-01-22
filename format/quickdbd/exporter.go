package quickdbd

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"io"
	"strings"
)

type ExportOption struct {
}

type Exporter struct {
	schema *octopus.Schema
	option *ExportOption
}

func (c *Exporter) Export(wr io.Writer) error {
	result := make([]string, 0)
	for _, table := range c.schema.Tables {
		result = append(result, getTableDef(table))
		result = append(result, strings.Repeat("-", len(table.Name)))

		for _, column := range table.Columns {
			result = append(result, getColumnDef(column))
		}
		result = append(result, "")
	}

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
	if !col.NotNull {
		params = append(params, "NULLABLE")
	}
	if col.DefaultValue != "" {
		params = append(params, fmt.Sprintf("default=%s", col.DefaultValue))
	}
	if ref := col.Ref; ref != nil {
		rel := getRelationshipType(ref)
		params = append(params, fmt.Sprintf("FK %s %s.%s", rel, ref.Table, ref.Column))
	}
	if col.Description != "" {
		params = append(params, fmt.Sprintf("# %s", col.Description))
	}

	return fmt.Sprintf("%s %s", col.Name, strings.Join(params, " "))
}

func getRelationshipType(ref *octopus.Reference) string {
	switch ref.Relationship {
	case octopus.RefManyToOne:
		return ">-"
	case octopus.RefOneToMany:
		return "-<"
	case octopus.RefOneToOne:
		return "-"
	default:
		return ">-"
	}
}
