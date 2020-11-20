package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type GormStruct struct {
	table        *Table
	Name         string
	EmbedModel   bool
	Fields       []*GormField
	PKFields     []*GormField
	UniqueFields []*GormField
}

func NewGormStruct(
	table *Table,
	output *Output,
	prefixMapper *PrefixMapper,
) *GormStruct {
	className := table.ClassName
	if className == "" {
		tableName := table.Name
		for _, prefix := range output.GetSlice(FlagRemovePrefix) {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		className = strcase.ToCamel(tableName)

		if prefix := prefixMapper.GetPrefix(table.Group); prefix != "" {
			className = prefix + className
		}
	}

	fields := make([]*GormField, 0)
	pkFields := make([]*GormField, 0)
	uniqueFields := make([]*GormField, 0)
	gormModelColumnCount := 0
	for _, column := range table.Columns {
		field := NewGormField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
		if column.UniqueKey {
			uniqueFields = append(uniqueFields, field)
		}

		if isGormModelColumn(column.Name) {
			gormModelColumnCount++
		}
	}

	return &GormStruct{
		table:        table,
		Name:         className,
		EmbedModel:   gormModelColumnCount < len(gormModelColumns),
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
	}
}

type GormField struct {
	Column       *Column
	Name         string
	Type         string
	OverrideName bool
	Imports      []string
}

func NewGormField(column *Column) *GormField {
	var fieldType string
	importSet := NewStringSet()
	if column.Nullable {
		importSet.Add("gopkg.in/guregu/null.v4")
	}

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case ColTypeString:
		fallthrough
	case ColTypeText:
		if column.Nullable {
			fieldType = "null.String"
		} else {
			fieldType = "string"
		}
	case ColTypeBoolean:
		if column.Nullable {
			fieldType = "null.Bool"
		} else {
			fieldType = "bool"
		}
	case ColTypeLong:
		fallthrough
	case ColTypeInt:
		if column.Nullable {
			fieldType = "null.Int"
		} else {
			colSize := column.Size
			if colSize > 0 {
				if colSize <= 3 {
					fieldType = "int8"
				} else if colSize <= 5 {
					fieldType = "int16"
				} else if colSize <= 10 {
					fieldType = "int32"
				} else {
					fieldType = "int64"
				}
			} else {
				fieldType = "int64"
			}
		}
	case ColTypeDecimal:
		importSet.Add("github.com/shopspring/decimal")
		if column.Nullable {
			fieldType = "decimal.NullDecimal"
		} else {
			fieldType = "decimal.Decimal"
		}
	case ColTypeFloat:
		fallthrough
	case ColTypeDouble:
		if column.Nullable {
			fieldType = "null.Float"
		} else {
			fieldType = "float64"
		}

	case ColTypeDateTime:
		fallthrough
	case ColTypeDate:
		fallthrough
	case ColTypeTime:
		if column.Nullable {
			fieldType = "null.Time"
		} else {
			importSet.Add("time")
			fieldType = "time.Time"
		}
	case ColTypeBlob:
		fieldType = "[]byte"
	default:
		if columnType == "bit" {
			if column.Size == 1 {
				if column.Nullable {
					fieldType = "null.Bool"
				} else {
					fieldType = "bool"
				}
				break
			}
		}
		fieldType = ""
	}

	fieldName, ok := ToUpperCamel(column.Name)
	overrideName := !ok

	// replace fieldname: Id -> ID
	if strings.HasSuffix(fieldName, "Id") {
		re := regexp.MustCompile(`Id$`)
		fieldName = string(re.ReplaceAll([]byte(fieldName), []byte("ID")))
	}
	// number prefix
	if matched, _ := regexp.MatchString(`\d+.*`, fieldName); matched {
		fieldName = "_" + fieldName
		overrideName = true
	}

	return &GormField{
		Column:       column,
		Name:         fieldName,
		Type:         fieldType,
		OverrideName: overrideName,
		Imports:      importSet.Slice(),
	}
}

var gormModelColumns = [...]string{"id", "created_at", "updated_at", "deleted_at"}

func isGormModelColumn(column string) bool {
	for _, gormModelColumn := range gormModelColumns {
		if column == gormModelColumn {
			return true
		}
	}
	return false
}

type GormTpl struct {
	schema       *Schema
	structs      []*GormStruct
	output       *Output
	prefixMapper *PrefixMapper
}

func NewGormTpl(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
) *GormTpl {
	structs := make([]*GormStruct, 0)
	for _, table := range schema.Tables {
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		structs = append(structs, NewGormStruct(table, output, prefixMapper))
	}

	return &GormTpl{
		schema:       schema,
		structs:      structs,
		output:       output,
		prefixMapper: prefixMapper,
	}
}

func (g *GormTpl) mkdir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	log.Printf("[MKDIR] %s", dir)
	return dir, nil
}

func (g *GormTpl) Generate() error {
	output := g.output
	gormModel := output.Get(FlagGormModel)
	if gormModel == "" {
		gormModel = "gorm.Model"
	}
	pkg := output.Get(FlagPackage)
	if pkg == "" {
		pkg = "main"
	}

	// write to single file if extension is '.go'
	var outputDir, filename string
	if ext := strings.ToLower(filepath.Ext(output.FilePath)); ext == ".go" {
		outputDir = filepath.Dir(output.FilePath)
		filename = output.FilePath
	} else {
		outputDir = output.FilePath
		filename = filepath.Join(output.FilePath, "output.go")
	}

	// ensure directory is created
	if _, err := g.mkdir(outputDir); err != nil {
		return err
	}

	gormStructs := make([]*GormStruct, 0)
	for _, table := range g.schema.Tables {
		gormStructs = append(gormStructs, NewGormStruct(table, output, g.prefixMapper))
	}

	buf := new(bytes.Buffer)

	// create import set
	importSet := NewStringSet()
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
	tplString := `{{"" -}}
package {{.Package}}

import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)
`
	tmpl, err := template.New("gormHeader").Parse(tplString)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(buf, &GormTplHeaderData{
		Package: pkg,
		Imports: importSet.Slice(),
	}); err != nil {
		return err
	}

	// write structs
	for _, gormStruct := range gormStructs {
		// write GORM struct
		if err := g.GenerateStruct(buf, gormStruct); err != nil {
			return err
		}
	}

	// write buffer to file
	if err := writeStringToFile(filename, buf.String()); err != nil {
		return err
	}

	return nil
}

type GormTplHeaderData struct {
	Package string
	Imports []string
}

type GormTplData struct {
	Package       string
	Struct        *GormStruct
	Table         *Table
	UniqueCstName string
	Fields        []*GormTplFieldData
}

type GormTplFieldData struct {
	Name string
	Type string
	Tag  string
}

func (g *GormTpl) GenerateStruct(
	wr io.Writer,
	gormStruct *GormStruct,
) error {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// unique constraint name
	uniqueCstName := ""
	uniqueFieldNames := make([]string, 0)
	for _, field := range gormStruct.UniqueFields {
		uniqueFieldNames = append(uniqueFieldNames, Quote(field.Name, "'"))
	}
	if len(uniqueFieldNames) > 1 {
		uniqueCstName = gormStruct.table.Name + g.output.Get(FlagUniqueNameSuffix)
	}

	// fields
	tplFields := make([]*GormTplFieldData, 0)
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

		if column.Type == "string" && column.Size > 0 {
			gormTags = append(gormTags, fmt.Sprintf("type:varchar(%d)", column.Size))
		} else if (column.Type == ColTypeDouble || column.Type == ColTypeFloat || column.Type == ColTypeDecimal) &&
			(column.Size > 0 && column.Scale > 0) {
			gormTags = append(gormTags, fmt.Sprintf("type:%s(%d,%d)", column.Type, column.Size, column.Scale))
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
		// auto_increment
		if column.AutoIncremental {
			gormTags = append(gormTags, "auto_increment")
		}
		// not null
		if !column.Nullable && !column.AutoIncremental {
			gormTags = append(gormTags, "not null")
		}

		// GORM tag
		var tag string
		if len(gormTags) > 0 {
			tag = fmt.Sprintf("`gorm:\"%s\"`", strings.Join(gormTags, ";"))
		}

		tplFields = append(tplFields, &GormTplFieldData{
			Name: field.Name,
			Type: field.Type,
			Tag:  tag,
		})
	}

	// populate template data
	data := GormTplData{
		Package:       g.output.Get(FlagPackage),
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
	{{.Name}} {{.Type}} {{.Tag}}
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
