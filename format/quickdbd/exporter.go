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
		result = append(result, c.getTableDef(table))
		result = append(result, strings.Repeat("-", len(table.Name)))

		for _, column := range table.Columns {
			result = append(result, c.getColumnDef(column))
		}
		result = append(result, "")
	}

	_, err := wr.Write([]byte(strings.Join(result, "\n")))
	return err
}

func (c *Exporter) getTableDef(table *octopus.Table) string {
	if table.Description != "" {
		return fmt.Sprintf("%s # %s", table.Name, table.Description)
	} else {
		return table.Name
	}
}

func (c *Exporter) getColumnDef(col *octopus.Column) string {
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
