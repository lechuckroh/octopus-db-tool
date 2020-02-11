package main

import (
	"fmt"
	"strings"
)

type TableFilterFn func(*Table) bool

type GenerateCmd struct {
}

func (cmd *GenerateCmd) Generate(input *Input, output *Output) error {
	schema, err := (&ConvertCmd{}).inputToSchema(input)
	if err != nil {
		return err
	}

	// table filter
	tableFilterFn := cmd.getTableFilterFn(output.Get(FlagGroups))

	switch output.Format {
	case FormatGraphql:
		graphql := &Graphql{}
		return graphql.Generate(schema, output, tableFilterFn)
	case FormatJpaKotlin:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output, tableFilterFn, false)
	case FormatJpaKotlinData:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output, tableFilterFn, true)
	case FormatLiquibase:
		liquibase := &Liquibase{}
		return liquibase.Generate(schema, output, tableFilterFn)
	}

	return fmt.Errorf("unsupported output format: %s", output.Format)
}

func (cmd *GenerateCmd) getTableFilterFn(groups string) TableFilterFn {
	if groups == "" {
		return nil
	}

	groupSlice := strings.Split(groups, ",")
	return func(table *Table) bool {
		for _, group := range groupSlice {
			if table.Group == group {
				return true
			}
		}
		return false
	}
}
