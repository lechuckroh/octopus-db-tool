package main

import (
	"errors"
	"io"
)

type ProtoFieldRule string

const (
	Required ProtoFieldRule = "required"
	Optional ProtoFieldRule = "optional"
	Repeated ProtoFieldRule = "repeated"
)

type ProtoMessage struct {
	Fields []*ProtoField
}

type ProtoField struct {
	Type    string
	Name    string
	Number  uint16
	Rule    ProtoFieldRule
	Comment string
	Default string
}

func NewProtobufMessage(table *Table, output *Output, mapper *PrefixMapper) *ProtoMessage {
	fields := make([]*ProtoField, 0)
	for i, column := range table.Columns {
		fields = append(fields, &ProtoField{
			Type:    column.Type,
			Name:    column.Name,
			Number:  uint16(i + 1),
			Rule:    Optional,
			Comment: column.Description,
			Default: column.DefaultValue,
		})
	}

	return &ProtoMessage{fields}
}

type ProtobufTpl struct {
	schema        *Schema
	output        *Output
	tableFilterFn TableFilterFn
	prefixMapper  *PrefixMapper
}

func NewProtobufTpl(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
) *ProtobufTpl {
	return &ProtobufTpl{
		schema:        schema,
		output:        output,
		tableFilterFn: tableFilterFn,
		prefixMapper:  prefixMapper,
	}
}

func (t *ProtobufTpl) Generate() error {
	// TODO:
	return errors.New("not implemented")
}

func (t *ProtobufTpl) GenerateProto(wr io.Writer, msg *ProtoMessage) error {
	// TODO:
	return errors.New("not implemented")
}

type ProtobufTplData struct {
	Package   string
	GoPackage string
}
