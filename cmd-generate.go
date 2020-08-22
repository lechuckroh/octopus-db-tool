package main

import (
	"fmt"
	"strings"
)

type TableFilterFn func(*Table) bool

type GenerateCmd struct {
}

func (cmd *GenerateCmd) Generate(input *Input, output *Output) error {
	schema, err := input.ToSchema()
	if err != nil {
		return err
	}

	// table filter
	tableFilterFn := cmd.getTableFilterFn(output.Get(FlagGroups))

	// annotation mapper
	annoMapper := newAnnotationMapper(output.Get(FlagAnnotation))
	// prefix mapper
	prefixMapper := newPrefixMapper(output.Get(FlagPrefix))

	switch output.Format {
	case FormatGorm:
		gorm := &Gorm{}
		return gorm.Generate(schema, output, tableFilterFn, prefixMapper)
	case FormatGraphql:
		graphql := &Graphql{}
		return graphql.Generate(schema, output, tableFilterFn, prefixMapper)
	case FormatJpaKotlin:
		jpa := NewJPAKotlin(schema, output, tableFilterFn, annoMapper, prefixMapper, false)
		return jpa.Generate()
	case FormatJpaKotlinData:
		jpa := NewJPAKotlin(schema, output, tableFilterFn, annoMapper, prefixMapper, true)
		return jpa.Generate()
	case FormatJpaKotlinTpl:
		jpa := NewJPAKotlinTpl(schema, output, tableFilterFn, annoMapper, prefixMapper, false)
		return jpa.Generate()
	case FormatJpaKotlinTplData:
		jpa := NewJPAKotlinTpl(schema, output, tableFilterFn, annoMapper, prefixMapper, true)
		return jpa.Generate()
	case FormatLiquibase:
		liquibase := &Liquibase{}
		return liquibase.Generate(schema, output, tableFilterFn)
	case FormatSqlalchemy:
		sqlAlchemy := &SqlAlchemy{}
		return sqlAlchemy.Generate(schema, output, tableFilterFn, prefixMapper)
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

type AnnotationMapper struct {
	anno    string
	annoMap map[string][]string
	useMap  bool
}

func newAnnotationMapper(annotation string) *AnnotationMapper {
	annoMap := make(map[string][]string)

	// populate annoMap
	if strings.Contains(annotation, ":") {
		for _, annoToken := range strings.Split(annotation, ",") {
			kv := strings.SplitN(annoToken, ":", 2)
			group := kv[0]
			annotations := strings.Split(kv[1], ";")
			annoMap[group] = annotations
		}
	}

	return &AnnotationMapper{
		anno:    annotation,
		annoMap: annoMap,
		useMap:  len(annoMap) > 0,
	}
}

func (m *AnnotationMapper) GetAnnotations(group string) []string {
	if m.useMap {
		if annotations, ok := m.annoMap[group]; ok {
			return annotations
		}
		return []string{}
	} else {
		return []string{m.anno}
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
