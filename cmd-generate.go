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

	// prefix mapper
	prefixMapper := newPrefixMapper(output.Get(FlagPrefix))

	switch output.Format {
	case FormatGraphql:
		graphql := &Graphql{}
		return graphql.Generate(schema, output, tableFilterFn, prefixMapper)
	case FormatJpaKotlin:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output, tableFilterFn, prefixMapper, false)
	case FormatJpaKotlinData:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output, tableFilterFn, prefixMapper, true)
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

type PrefixMapper struct {
	prefix    string
	prefixMap map[string]string
	useMap    bool
}

func newPrefixMapper(prefix string) *PrefixMapper {
	prefixMap := make(map[string]string)

	// populate prefixMap
	if strings.Contains(prefix, ":") {
		for _, prefixToken := range strings.Split(prefix, ",") {
			kv := strings.SplitN(prefixToken, ":", 2)
			group := kv[0]
			prefixValue := kv[1]
			prefixMap[group] = prefixValue
		}
	}

	return &PrefixMapper{
		prefix:    prefix,
		prefixMap: prefixMap,
		useMap:    len(prefixMap) > 0,
	}
}

func (p *PrefixMapper) GetPrefix(group string) string {
	if p.useMap {
		if prefix, ok := p.prefixMap[group]; ok {
			return prefix
		}
		return ""
	} else {
		return p.prefix
	}
}
