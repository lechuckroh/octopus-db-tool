package dbml

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"strings"
)

type Option struct {
	TableFilter octopus.TableFilterFn
}

type Exporter struct {
	schema *octopus.Schema
	option *Option
}

func (c *Exporter) Export(wr io.Writer) error {
	// TODO: replace with template
	if bytes, err := c.exportToString(); err != nil {
		return err
	} else {
		_, err := wr.Write(bytes)
		return err
	}
}

func (c *Exporter) exportToString() ([]byte, error) {
	result := make([]string, 0)
	definedTables := make(map[string]bool)
	deferredRefs := make([]string, 0)

	for _, table := range c.schema.Tables {
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
			if column.NotNull {
				params = append(params, "not null")
			}
			if column.DefaultValue != "" {
				defaultValue := column.DefaultValue
				if util.IsStringType(column.Type) {
					defaultValue = fmt.Sprintf("'%s'", defaultValue)
				}
				params = append(params, fmt.Sprintf("default: %s", defaultValue))
			}
			if ref := column.Ref; ref != nil {
				rel := getRelationshipType(ref)

				if definedTables[ref.Table] {
					params = append(params,
						fmt.Sprintf("ref: %s %s.%s", rel, ref.Table, ref.Column))
				} else {
					deferredRefs = append(deferredRefs,
						fmt.Sprintf("Ref: %s.%s %s %s.%s",
							table.Name, column.Name, rel, ref.Table, ref.Column),
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

func getRelationshipType(ref *octopus.Reference) string {
	switch ref.Relationship {
	case octopus.RefManyToOne:
		return ">"
	case octopus.RefOneToMany:
		return "<"
	case octopus.RefOneToOne:
		return "-"
	default:
		return ">"
	}
}
