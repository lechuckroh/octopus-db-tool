package protobuf

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"text/template"
)

type ProtoFieldRule string

const (
	Required ProtoFieldRule = "required"
	Optional ProtoFieldRule = "optional"
	Repeated ProtoFieldRule = "repeated"
)

type Option struct {
	PrefixMapper   *common.PrefixMapper
	TableFilter    octopus.TableFilterFn
	RemovePrefixes []string
	Package        string
	GoPackage      string
	FilePath       string
}

type Generator struct {
	schema   *octopus.Schema
	messages []*PbMessage
	option   *Option
}

func NewGenerator(
	schema *octopus.Schema,
	option *Option,
) *Generator {
	messages := make([]*PbMessage, 0)
	for _, table := range schema.Tables {
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}
		messages = append(messages, NewPbMessage(table, option))
	}

	return &Generator{
		schema:   schema,
		option:   option,
		messages: messages,
	}
}

func (t *Generator) Generate(wr io.Writer) error {
	option := t.option

	pkg := option.Package
	goPkg := option.GoPackage

	return t.generateProto(wr, t.messages, pkg, goPkg)
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
	options := make([]string, 0)
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
