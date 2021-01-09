package gorm

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"strings"
	"text/template"
)

type Option struct {
	PrefixMapper     *common.PrefixMapper
	TableFilter      octopus.TableFilterFn
	GormModel        string
	Package          string
	RemovePrefixes   []string
	UniqueNameSuffix string
}

type Generator struct {
	schema *octopus.Schema
	option *Option
}

func (g *Generator) Generate(wr io.Writer) error {
	option := g.option
	gormModel := option.GormModel
	if gormModel == "" {
		gormModel = "gorm.Model"
	}
	pkg := option.Package
	if pkg == "" {
		pkg = "main"
	}

	gormStructs := make([]*GoStruct, 0)
	tableFilter := option.TableFilter
	for _, table := range g.schema.Tables {
		if tableFilter == nil || tableFilter(table) {
			gormStructs = append(gormStructs, NewGoStruct(table, option))
		}
	}

	// create import set
	importSet := util.NewStringSet()
	for _, gormStruct := range gormStructs {
		if gormStruct.EmbedModel {
			importSet.Add("github.com/jinzhu/gorm")
		}

		for _, field := range gormStruct.Fields {
			for _, imp := range field.Imports {
				importSet.Add(imp)
			}
		}
	}

	// generate header
	funcMap := template.FuncMap{}
	tplText := `{{"" -}}
package {{.Package}}

import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

`
	tmpl, err := util.NewTemplate("gormHeader", tplText, funcMap)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, &TplHeaderData{
		Package: pkg,
		Imports: importSet.Slice(),
	}); err != nil {
		return err
	}

	// write structs
	for _, gormStruct := range gormStructs {
		// write GORM struct
		if err := g.GenerateStruct(wr, gormStruct); err != nil {
			return err
		}
	}

	return nil
}

type TplHeaderData struct {
	Package string
	Imports []string
}

type TplData struct {
	Package       string
	Struct        *GoStruct
	Table         *octopus.Table
	UniqueCstName string
	Fields        []*TplFieldData
}

type TplFieldData struct {
	Name string
	Type string
	Tag  string
}

func (f *TplFieldData) ToString() string {
	if f.Tag == "" {
		return f.Name + " " + f.Type
	} else {
		return f.Name + " " + f.Type + " " + f.Tag
	}
}

type IndexTag struct {
	IndexName   string
	Priority    int
	SingleIndex bool
}

func (g *Generator) getGormIndexTag(indices *[]*octopus.Index, field *GoField) []*IndexTag {
	fieldColumnName := field.Column.Name

	var gormIndexTags []*IndexTag

	for _, index := range *indices {
		singleIndexColumn := len(index.Columns) == 1

		for i, col := range index.Columns {
			if fieldColumnName == col {
				gormIndexTags = append(gormIndexTags, &IndexTag{
					IndexName:   index.Name,
					Priority:    i + 1,
					SingleIndex: singleIndexColumn,
				})
				break
			}
		}
	}
	return gormIndexTags
}

func (g *Generator) GenerateStruct(
	wr io.Writer,
	gormStruct *GoStruct,
) error {
	funcMap := template.FuncMap{
		"join": strings.Join,
		"fieldToString": func(field *TplFieldData) string {
			return field.ToString()
		},
	}

	// unique constraint name
	uniqueCstName := ""
	uniqueFieldNames := make([]string, 0)
	for _, field := range gormStruct.UniqueFields {
		uniqueFieldNames = append(uniqueFieldNames, util.Quote(field.Name, "'"))
	}
	if len(uniqueFieldNames) > 1 {
		uniqueCstName = gormStruct.table.Name + g.option.UniqueNameSuffix
	}

	// fields
	tplFields := make([]*TplFieldData, 0)
	for _, field := range gormStruct.Fields {
		column := field.Column

		// embedded model column
		if gormStruct.EmbedModel && isGormModelColumn(column.Name) {
			continue
		}

		// Column gormTags
		gormTags := make([]string, 0)

		if field.OverrideName {
			gormTags = append(gormTags, fmt.Sprintf("column:%s", column.Name))
		}

		if tag := g.getGormTagByType(column); tag != "" {
			gormTags = append(gormTags, tag)
		}

		// PK
		if column.PrimaryKey {
			gormTags = append(gormTags, "primary_key")
		}
		// Unique
		if column.UniqueKey {
			if uniqueCstName == "" {
				gormTags = append(gormTags, "unique")
			} else {
				gormTags = append(gormTags, fmt.Sprintf("unique_index:%s", uniqueCstName))
			}
		}
		// Index
		for _, indexTag := range g.getGormIndexTag(&gormStruct.table.Indices, field) {
			if indexTag.SingleIndex {
				gormTags = append(gormTags, fmt.Sprintf("index:%s", indexTag.IndexName))
			} else {
				gormTags = append(gormTags, fmt.Sprintf("index:%s,priority:%d", indexTag.IndexName, indexTag.Priority))
			}
		}

		// auto_increment
		if column.AutoIncremental {
			gormTags = append(gormTags, "auto_increment")
		}
		// not null
		if column.NotNull && !column.AutoIncremental {
			gormTags = append(gormTags, "not null")
		}

		// GORM tag
		var tag string
		if len(gormTags) > 0 {
			tag = fmt.Sprintf("`gorm:\"%s\"`", strings.Join(gormTags, ";"))
		}

		tplFields = append(tplFields, &TplFieldData{
			Name: field.Name,
			Type: field.Type,
			Tag:  tag,
		})
	}

	// populate template data
	data := TplData{
		Package:       g.option.Package,
		Struct:        gormStruct,
		Table:         gormStruct.table,
		UniqueCstName: uniqueCstName,
		Fields:        tplFields,
	}

	tplString := `{{"" -}}
type {{.Struct.Name}} struct {
{{- if .Struct.EmbedModel}}
	gorm.Model
{{- end}}
{{- range .Fields}}
	{{fieldToString .}}
{{- end}}
}

func (c *{{.Struct.Name}}) TableName() string { return "{{.Table.Name}}" }

`
	tmpl, err := template.New("gormStruct").Funcs(funcMap).Parse(tplString)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, &data)
}

func (g *Generator) getGormTagByType(column *octopus.Column) string {
	size := column.Size
	switch column.Type {
	case octopus.ColTypeBit:
		if size > 0 {
			return fmt.Sprintf("type:%s(%d)", column.Type, column.Size)
		}
	case octopus.ColTypeChar:
		fallthrough
	case octopus.ColTypeVarchar:
		if size > 0 {
			return fmt.Sprintf("type:%s(%d)", column.Type, column.Size)
		}
	case octopus.ColTypeDouble:
		fallthrough
	case octopus.ColTypeFloat:
		fallthrough
	case octopus.ColTypeDecimal:
		if size > 0 && column.Scale > 0 {
			return fmt.Sprintf("type:%s(%d,%d)", column.Type, column.Size, column.Scale)
		}
	}
	return ""
}
