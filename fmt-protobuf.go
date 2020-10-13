package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"io"
	"log"
	"strings"
	"text/template"
)

type ProtoFieldRule string

const (
	Required ProtoFieldRule = "required"
	Optional ProtoFieldRule = "optional"
	Repeated ProtoFieldRule = "repeated"
)

type ProtoMessage struct {
	Name    string
	Fields  []*ProtoField
	Imports []string
}

type ProtoField struct {
	Type    string
	Name    string
	Number  int
	Rule    ProtoFieldRule
	Comment string
	Default string
	Import  string
}

func NewProtobufMessage(table *Table, output *Output, prefixMapper *PrefixMapper) *ProtoMessage {
	imports := NewStringSet()
	fields := make([]*ProtoField, 0)
	for i, column := range table.Columns {
		field := NewProtoField(i+1, column)
		fields = append(fields, field)
		if field.Import != "" {
			imports.Add(field.Import)
		}
	}

	tableName := table.Name
	for _, prefix := range output.GetSlice(FlagRemovePrefix) {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	messageName := strcase.ToCamel(tableName)

	if prefixMapper != nil {
		if prefix := prefixMapper.GetPrefix(table.Group); prefix != "" {
			messageName = prefix + messageName
		}
	}

	return &ProtoMessage{
		Name:    messageName,
		Fields:  fields,
		Imports: imports.Slice(),
	}
}

func NewProtoField(number int, column *Column) *ProtoField {
	var fieldType string
	var imp string

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case ColTypeString:
		fallthrough
	case ColTypeText:
		fieldType = "string"
	case ColTypeBoolean:
		fieldType = "bool"
	case ColTypeLong:
		fieldType = "int64"
	case ColTypeInt:
		if column.Size <= 10 {
			fieldType = "int32"
		} else {
			fieldType = "int64"
		}
	case ColTypeDecimal:
		fieldType = "double"
	case ColTypeFloat:
		fieldType = "float"
	case ColTypeDouble:
		fieldType = "double"
	case ColTypeDateTime:
		fieldType = "google.protobuf.Timestamp"
		imp = "google/protobuf/timestamp.proto"
	default:
		if columnType == "bit" {
			if column.Size == 1 {
				fieldType = "bool"
				break
			}
		}
		fieldType = column.Type
		log.Printf("unsupported columnType: %s", columnType)
	}

	fieldName, _ := ToLowerCamel(column.Name)

	return &ProtoField{
		Type:    fieldType,
		Name:    fieldName,
		Number:  number,
		Rule:    Optional,
		Comment: column.Description,
		Default: column.DefaultValue,
		Import:  imp,
	}
}

type ProtobufTpl struct {
	schema       *Schema
	messages     []*ProtoMessage
	output       *Output
	prefixMapper *PrefixMapper
}

func NewProtobufTpl(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
) *ProtobufTpl {
	messages := make([]*ProtoMessage, 0)
	for _, table := range schema.Tables {
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		messages = append(messages, NewProtobufMessage(table, output, prefixMapper))
	}

	return &ProtobufTpl{
		schema:       schema,
		output:       output,
		messages:     messages,
		prefixMapper: prefixMapper,
	}
}

func (t *ProtobufTpl) Generate() error {
	output := t.output

	pkg := output.Get(FlagPackage)
	goPkg := output.Get(FlagGoPackage)

	buf := new(bytes.Buffer)
	if err := t.GenerateProto(buf, t.messages, pkg, goPkg); err != nil {
		return err
	}
	filename := output.FilePath
	if err := writeStringToFile(filename, buf.String()); err != nil {
		return err
	}

	return nil
}

func (t *ProtobufTpl) GenerateProto(
	wr io.Writer,
	messages []*ProtoMessage,
	pkg string,
	goPkg string,
) error {
	// custom functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"hasNext": func(field *KotlinField, fields []*KotlinField) bool {
			return field != fields[len(fields)-1]
		},
	}

	tplString := `{{"" -}}
syntax = "proto3";

package {{.Package}};
{{range .Options}}
option {{.}};
{{end}}
{{- range .Imports}}
import "{{.}}";
{{end}}

{{- range .Messages}}
message {{.Name}} {
  {{- range .Fields}}
  {{.Type}} {{.Name}} = {{.Number}};
  {{- end}}
}
{{end}}`

	// parse template
	tmpl, err := template.New("protobuf").Funcs(funcMap).Parse(tplString)
	if err != nil {
		return err
	}

	// import
	imports := NewStringSet()
	for _, message := range messages {
		imports.AddAll(message.Imports)
	}

	// options
	options := make([]string, 0)
	if goPkg != "" {
		options = append(options, fmt.Sprintf("go_package = \"%s\"", goPkg))
	}


	data := struct {
		Package  string
		Options  []string
		Imports  []string
		Messages []*ProtoMessage
	}{
		Package:  pkg,
		Options:  options,
		Imports:  imports.Slice(),
		Messages: messages,
	}
	return tmpl.Execute(wr, &data)
}
