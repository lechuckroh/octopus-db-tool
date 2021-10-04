package protobuf

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
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

type PbMessage struct {
	Name      string
	Fields    []*PbField
	Relations []*PbRelation
	Imports   []string
}

type PbField struct {
	Type    string
	Name    string
	Tag     int
	Rule    ProtoFieldRule
	Comment string
	Default string
	Import  string
}

type PbRelation struct {
	Type     string
	Name     string
	Tag      int
	Repeated bool
}

type Option struct {
	PrefixMapper     *common.PrefixMapper
	TableFilter      octopus.TableFilterFn
	RemovePrefixes   []string
	Package          string
	GoPackage        string
	FilePath         string
	RelationTagStart int
	RelationTagDecr  bool
}

type Generator struct {
	schema              *octopus.Schema
	option              *Option
	messageNameMemoizer func(interface{}) string
}

func newGenerator(schema *octopus.Schema, option *Option) *Generator {
	gen := &Generator{
		schema: schema,
		option: option,
	}
	gen.messageNameMemoizer = util.NewStringMemoizer(func(table interface{}) string {
		return gen.getMessageName(table.(*octopus.Table))
	})

	return gen
}

func (t *Generator) Generate(wr io.Writer) error {
	option := t.option

	pkg := option.Package
	goPkg := option.GoPackage

	var pbMessages []*PbMessage
	for _, table := range t.schema.Tables {
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}
		pbMessages = append(pbMessages, t.newPbMessage(table))
	}

	return t.generateProto(wr, pbMessages, pkg, goPkg)
}

func (t *Generator) generateProto(
	wr io.Writer,
	messages []*PbMessage,
	pkg string,
	goPkg string,
) error {
	// custom functions
	funcMap := template.FuncMap{}

	tplText := `{{"" -}}
syntax = "proto3";

{{if .Package}}package {{.Package}};{{end}}
{{range .Options}}
option {{.}};
{{end}}
{{- range .Imports}}
import "{{.}}";
{{end}}

{{- range .Messages}}
message {{.Name}} {
  {{- range .Fields}}
  {{.Type}} {{.Name}} = {{.Tag}};
  {{- end}}
  {{- range .Relations}}
  {{.Type}} {{.Name}} = {{.Tag}};
  {{- end}}
}
{{end}}`

	// parse template
	tmpl, err := util.NewTemplate("protobuf", tplText, funcMap)
	if err != nil {
		return err
	}

	// import
	imports := util.NewStringSet()
	for _, message := range messages {
		imports.AddAll(message.Imports)
	}

	// options
	var options []string
	if goPkg != "" {
		options = append(options, fmt.Sprintf("go_package = \"%s\"", goPkg))
	}

	data := struct {
		Package  string
		Options  []string
		Imports  []string
		Messages []*PbMessage
	}{
		Package:  pkg,
		Options:  options,
		Imports:  imports.Slice(),
		Messages: messages,
	}
	return tmpl.Execute(wr, &data)
}

// getMessageName returns protobuf message name from table name.
func (t *Generator) getMessageName(table *octopus.Table) string {
	tableName := table.Name

	for _, prefix := range t.option.RemovePrefixes {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	messageName := strcase.ToCamel(tableName)

	if prefixMapper := t.option.PrefixMapper; prefixMapper != nil {
		if prefix := prefixMapper.GetPrefix(table.Group); prefix != "" {
			messageName = prefix + messageName
		}
	}
	return messageName
}

func (t *Generator) newPbMessage(table *octopus.Table) *PbMessage {
	// field list
	imports := util.NewStringSet()
	var fields []*PbField
	var relations []*PbRelation
	maxTag := 0
	for i, column := range table.Columns {
		tag := i + 1
		maxTag = tag
		field := newPbField(tag, column)
		fields = append(fields, field)
		if imp := field.Import; imp != "" {
			imports.Add(imp)
		}

		// reference
		if ref := column.Ref; ref != nil {
			refTable := t.schema.TableByName(ref.Table)
			refMessageType := t.messageNameMemoizer(refTable)
			refName := strcase.ToLowerCamel(refMessageType)
			repeated := ref.Relationship == octopus.RefOneToMany

			relations = append(relations, &PbRelation{
				Type:     refMessageType,
				Name:     refName,
				Tag:      -1,
				Repeated: repeated,
			})
		}
	}

	// shift tag start if tag range is overlapped
	relationCount := len(relations)
	relTagStart := t.option.RelationTagStart
	relTagDecr := t.option.RelationTagDecr
	if relTagDecr {
		minStartTag := relTagStart - relationCount + 1
		if minStartTag <= maxTag {
			relTagStart = maxTag + relationCount
		}
	} else {
		if relTagStart < 0 || relTagStart <= maxTag {
			relTagStart = maxTag + 1
		}
	}

	// update tag of relations
	for i, relation := range relations {
		relation.Tag = relTagStart + i
	}

	return &PbMessage{
		Name:      t.getMessageName(table),
		Fields:    fields,
		Relations: relations,
		Imports:   imports.Slice(),
	}
}

func newPbField(tag int, column *octopus.Column) *PbField {
	var fieldType string
	var imp string

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case octopus.ColTypeChar:
		fallthrough
	case octopus.ColTypeVarchar:
		fallthrough
	case octopus.ColTypeText8:
		fallthrough
	case octopus.ColTypeText16:
		fallthrough
	case octopus.ColTypeText24:
		fallthrough
	case octopus.ColTypeText32:
		fieldType = "string"
	case octopus.ColTypeBoolean:
		fieldType = "bool"
	case octopus.ColTypeInt64:
		fieldType = "int64"
	case octopus.ColTypeInt8:
		fallthrough
	case octopus.ColTypeInt16:
		fallthrough
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		fieldType = "int32"
	case octopus.ColTypeDecimal:
		fieldType = "double"
	case octopus.ColTypeFloat:
		fieldType = "float"
	case octopus.ColTypeDouble:
		fieldType = "double"
	case octopus.ColTypeDateTime:
		fieldType = "google.protobuf.Timestamp"
		imp = "google/protobuf/timestamp.proto"
	default:
		fieldType = column.Type
		log.Printf("unsupported columnType: %s", columnType)
	}

	fieldName, _ := util.ToLowerCamel(column.Name)

	return &PbField{
		Type:    fieldType,
		Name:    fieldName,
		Tag:     tag,
		Rule:    Optional,
		Comment: column.Description,
		Default: column.DefaultValue,
		Import:  imp,
	}
}
