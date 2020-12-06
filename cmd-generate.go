package main

import (
	"fmt"
	"strings"
)

type GenerateCmd struct {
}

func (cmd *GenerateCmd) Generate(input *Input, output *Output) error {
	schema, err := input.ToSchema()
	if err != nil {
		return err
	}

	// table filter
	tableFilterFn := getTableFilterFn(output.Get(FlagGroups))

	// annotation mapper
	annoMapper := newAnnotationMapper(output.Get(FlagAnnotation))
	// prefix mapper
	prefixMapper := newPrefixMapper(output.Get(FlagPrefix))

	switch output.Format {
	case FormatGorm:
		gorm := NewGormTpl(schema, output, tableFilterFn, prefixMapper)
		return gorm.Generate()
	case FormatGraphql:
		graphql := &Graphql{}
		return graphql.Generate(schema, output, tableFilterFn, prefixMapper)
	case FormatJpaKotlin:
		jpa := NewJPAKotlinTpl(schema, output, tableFilterFn, annoMapper, prefixMapper)
		return jpa.Generate()
	case FormatLiquibase:
		liquibase := &Liquibase{}
		return liquibase.Generate(schema, output, tableFilterFn)
	case FormatProtobuf:
		protobuf := NewProtobufTpl(schema, output, tableFilterFn, prefixMapper)
		return protobuf.Generate()
	case FormatSqlalchemy:
		sqlAlchemy := &SqlAlchemy{}
		return sqlAlchemy.Generate(schema, output, tableFilterFn, prefixMapper)
	}

	return fmt.Errorf("unsupported output format: %s", output.Format)
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
