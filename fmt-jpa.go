package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"io"
	"log"
	"path"
	"strings"
	"text/template"
)

type KotlinClass struct {
	table        *Table
	Name         string
	Annotations  []string
	Fields       []*KotlinField
	PKFields     []*KotlinField
	UniqueFields []*KotlinField
}

type KotlinField struct {
	Column       *Column
	Name         string
	OverrideName bool
	Type         string
	Imports      []string
	DefaultValue string
}

func NewKotlinClass(
	table *Table,
	output *Output,
	annoMapper *AnnotationMapper,
	prefixMapper *PrefixMapper,
) *KotlinClass {
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

	fields := make([]*KotlinField, 0)
	pkFields := make([]*KotlinField, 0)
	uniqueFields := make([]*KotlinField, 0)
	for _, column := range table.Columns {
		field := NewKotlinField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
		if column.UniqueKey {
			uniqueFields = append(uniqueFields, field)
		}
	}

	return &KotlinClass{
		table:        table,
		Name:         className,
		Annotations:  annoMapper.GetAnnotations(table.Group),
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
	}
}

func NewKotlinField(column *Column) *KotlinField {
	var fieldType string
	var defaultValue string
	nullable := column.Nullable

	importSet := NewStringSet()

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case ColTypeString:
		fallthrough
	case ColTypeText:
		fieldType = "String"
		if !nullable {
			defaultValue = "\"\""
		}
	case ColTypeBoolean:
		fieldType = "Boolean"
		if !nullable {
			defaultValue = "false"
		}
	case ColTypeLong:
		fieldType = "Long"
		if !nullable {
			defaultValue = "0L"
		}
	case ColTypeInt:
		fieldType = "Int"
		if !nullable {
			defaultValue = "0"
		}
	case ColTypeDecimal:
		fieldType = "BigDecimal"
		importSet.Add("java.math.BigDecimal")
		if !nullable {
			defaultValue = "BigDecimal.ZERO"
		}
	case ColTypeFloat:
		fieldType = "Float"
		if !nullable {
			defaultValue = "0.0F"
		}
	case ColTypeDouble:
		fieldType = "Double"
		if !nullable {
			defaultValue = "0.0"
		}
	case ColTypeDateTime:
		fieldType = "Timestamp"
		importSet.Add("java.sql.Timestamp")
		if !nullable {
			defaultValue = "Timestamp(System.currentTimeMillis())"
		}
	case ColTypeDate:
		fieldType = "LocalDate"
		importSet.Add("java.time.LocalDate")
		if !nullable {
			defaultValue = "LocalDate.now()"
		}
	case ColTypeTime:
		fieldType = "LocalTime"
		importSet.Add("java.time.LocalTime")
		if !nullable {
			defaultValue = "LocalTime.now()"
		}
	case ColTypeBlob:
		fieldType = "Blob"
		importSet.Add("java.sql.Blob")
	default:
		if columnType == "bit" {
			if column.Size == 1 {
				fieldType = "Boolean"
				if !nullable {
					defaultValue = "false"
				}
				break
			}
		}
		fieldType = "Any"
	}
	if nullable {
		fieldType = fieldType + "?"
	}

	fieldName, ok := ToLowerCamel(column.Name)

	return &KotlinField{
		Column:       column,
		Name:         fieldName,
		OverrideName: !ok,
		Type:         fieldType,
		DefaultValue: defaultValue,
		Imports:      importSet.Slice(),
	}
}

type JPAKotlinTpl struct {
	schema       *Schema
	classes      []*KotlinClass
	output       *Output
	annoMapper   *AnnotationMapper
	prefixMapper *PrefixMapper
}

func NewJPAKotlinTpl(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	annoMapper *AnnotationMapper,
	prefixMapper *PrefixMapper,
) *JPAKotlinTpl {
	// populate KotlinClass
	classes := make([]*KotlinClass, 0)
	for _, table := range schema.Tables {
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		classes = append(classes, NewKotlinClass(table, output, annoMapper, prefixMapper))
	}

	return &JPAKotlinTpl{
		schema:       schema,
		classes:      classes,
		output:       output,
		annoMapper:   annoMapper,
		prefixMapper: prefixMapper,
	}
}

func (k *JPAKotlinTpl) getFieldAnnotations(class *KotlinClass, field *KotlinField) []string {
	annotations := make([]string, 0)

	column := field.Column
	colAttrs := make([]string, 0)
	if column.PrimaryKey {
		annotations = append(annotations, "@Id")
	}
	if column.AutoIncremental {
		annotations = append(annotations, "@GeneratedValue(strategy = GenerationType.AUTO)")
	}
	if column.Type == ColTypeText {
		annotations = append(annotations, "@Lob")
	}

	// @VRelation
	if k.output.Get(FlagRelation) == "VRelation" {
		if ref := column.Ref; ref != nil {
			targetClassName := k.getClassNameByTableName(ref.Table)
			if len(targetClassName) == 0 {
				log.Fatalf("Relation not found. %s::%s -> %s", class.Name, field.Name, ref.Table)
			}

			annotations = append(annotations,
				fmt.Sprintf("@VRelation(cls = \"%s\", field = \"%s\")",
					targetClassName,
					strcase.ToLowerCamel(ref.Column)))
		}
	}

	if field.OverrideName {
		colAttrs = append(colAttrs, fmt.Sprintf("name = \"%s\"", column.Name))
	}
	if !column.Nullable {
		colAttrs = append(colAttrs, "nullable = false")
	}
	if column.Type == ColTypeString && column.Size > 0 {
		colAttrs = append(colAttrs, fmt.Sprintf("length = %d", column.Size))
	}
	if column.Type == ColTypeDouble || column.Type == ColTypeFloat || column.Type == ColTypeDecimal {
		if column.Size > 0 {
			colAttrs = append(colAttrs, fmt.Sprintf("precision = %d", column.Size))
		}
		if column.Scale > 0 {
			colAttrs = append(colAttrs, fmt.Sprintf("scale = %d", column.Scale))
		}
	}
	// @CreationTimestamp
	if column.Type == ColTypeDateTime && field.Name == "createdAt" {
		annotations = append(annotations, "@CreationTimestamp")
		colAttrs = append(colAttrs, "updatable = false")
	}
	// @UpdateTimestamp
	if column.Type == ColTypeDateTime && field.Name == "updatedAt" {
		annotations = append(annotations, "@UpdateTimestamp")
	}
	if len(colAttrs) > 0 {
		annotations = append(annotations, fmt.Sprintf("@Column(%s)", strings.Join(colAttrs, ", ")))
	}
	return annotations
}

type JPAKotlinTplData struct {
	Package              string
	Class                *KotlinClass
	SuperClass           string
	Annotations          []string
	Table                *Table
	IdEntityField        *KotlinField
	IdClassName          string
	UniqueConstraintName string
	UniqueFieldNames     []string
	HasUniqueFields      bool
	UseDataClass         bool
	Imports              []string
	JavaImports          []string
}

// GenerateEntityClass generates entity class
func (k *JPAKotlinTpl) GenerateEntityClass(
	wr io.Writer,
	class *KotlinClass,
) error {
	// custom functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"fieldAnnotations": func(field *KotlinField) []string {
			return k.getFieldAnnotations(class, field)
		},
		"hasNext": func(field *KotlinField, fields []*KotlinField) bool {
			return field != fields[len(fields)-1]
		},
	}

	// template
	tplString := `{{"" -}}
package {{.Package}}

{{range .Imports}}import {{.}}
{{end}}
{{range .JavaImports}}import {{.}}
{{end}}
{{range .Annotations}}{{.}}{{end}}
@Entity
{{- if .HasUniqueFields}}
@Table(name="{{.Table.Name}}", uniqueConstraints = [
    UniqueConstraint(name = "{{.UniqueConstraintName}}", columnNames = [{{join .UniqueFieldNames ", "}}])
])
{{- else}}
@Table(name = "{{.Table.Name}}")
{{- end}}
{{- if ne .IdClassName ""}}
@IdClass({{.IdClassName}}::class)
{{- end}}
data class {{.Class.Name}}(
{{- range .Class.Fields}}
    {{- $annotations := fieldAnnotations .}}
    {{- range $annotations}}
        {{. -}}
    {{end}}
    {{- if .Column.Nullable}}
        var {{.Name}}: {{.Type}}{{if hasNext . $.Class.Fields}},{{end}}
    {{- else}}
		{{- if eq . $.IdEntityField}}
        override var {{.Name}}: {{.Type}} = {{.DefaultValue}},
		{{- else}}
        var {{.Name}}: {{.Type}} = {{.DefaultValue}}{{if hasNext . $.Class.Fields}},{{end}}
		{{- end}}
    {{- end}}
{{end}}
{{- if eq .SuperClass ""}}
)
{{- else}}
): {{.SuperClass}}
{{- end}}

{{- if ne .IdClassName ""}}
data class {{.IdClassName}}(
{{- range .Class.PKFields}}
        var {{.Name}}: {{.Type}} = {{.DefaultValue}}{{if hasNext . $.Class.PKFields}},{{end}}
{{- end}}
): Serializable
{{- end}}
`

	// parse template
	tmpl, err := template.New("jpaKotlinEntity").Funcs(funcMap).Parse(tplString)
	if err != nil {
		return err
	}

	// PK
	var idClassName string
	pkFieldCount := len(class.PKFields)
	if pkFieldCount > 1 {
		idClassName = class.Name + "PK"
	}

	// idEntity field
	var idEntityField *KotlinField
	idEntityInterfaceName := k.output.Get(FlagIdEntity)
	if idEntityInterfaceName != "" && pkFieldCount == 1 && class.PKFields[0].Name == "id" {
		idEntityField = class.PKFields[0]
	}

	// annotations
	annotations := k.annoMapper.GetAnnotations(class.table.Group)

	// super class
	superClass := ""
	if idEntityField != nil {
		superClass = fmt.Sprintf("%s<%s>", idEntityInterfaceName, idEntityField.Type)
	}

	// unique constraint
	uniqueFieldNames := make([]string, 0)
	for _, field := range class.UniqueFields {
		uniqueFieldNames = append(uniqueFieldNames, Quote(field.Name, "\""))
	}

	// imports
	importSet := NewStringSet()
	javaImportSet := NewStringSet()
	javaImportSet.Add("javax.persistence.*")
	if pkFieldCount > 1 {
		javaImportSet.Add("java.io.Serializable")
	}

	for _, field := range class.Fields {
		column := field.Column
		for _, imp := range field.Imports {
			if strings.HasPrefix(imp, "java.") {
				javaImportSet.Add(imp)
			} else {
				importSet.Add(imp)
			}
		}

		// @CreationTimestamp
		if column.Type == ColTypeDateTime && field.Name == "createdAt" {
			importSet.Add("org.hibernate.annotations.CreationTimestamp")
		}
		// @UpdateTimestamp
		if column.Type == ColTypeDateTime && field.Name == "updatedAt" {
			importSet.Add("org.hibernate.annotations.UpdateTimestamp")
		}
	}

	// populate template data
	data := JPAKotlinTplData{
		Package:              k.output.Get(FlagPackage),
		Class:                class,
		SuperClass:           superClass,
		Annotations:          annotations,
		Table:                class.table,
		IdEntityField:        idEntityField,
		IdClassName:          idClassName,
		UniqueConstraintName: class.table.Name + k.output.Get(FlagUniqueNameSuffix),
		UniqueFieldNames:     uniqueFieldNames,
		HasUniqueFields:      len(uniqueFieldNames) > 0,
		Imports:              importSet.Slice(),
		JavaImports:          javaImportSet.Slice(),
	}

	return tmpl.Execute(wr, &data)
}

func (k *JPAKotlinTpl) getClassNameByTableName(tableName string) string {
	for _, cls := range k.classes {
		if cls.table.Name == tableName {
			return cls.Name
		}
	}
	return ""
}

func (k *JPAKotlinTpl) Generate() error {
	output := k.output
	outputPackage := output.Get(FlagPackage)
	reposPackage := output.Get(FlagReposPackage)

	entityDir, err := mkdirPackage(output.FilePath, outputPackage)
	if err != nil {
		return err
	}
	reposDir, err := mkdirPackage(output.FilePath, reposPackage)
	if err != nil {
		return err
	}

	classes := k.classes

	for _, class := range classes {
		// write entity class
		buf := new(bytes.Buffer)
		if err := k.GenerateEntityClass(buf, class); err != nil {
			return err
		}
		filename := path.Join(entityDir, fmt.Sprintf("%s.kt", class.Name))
		if err := writeStringToFile(filename, buf.String()); err != nil {
			return err
		}

		// write repository
		buf = new(bytes.Buffer)
		if err := k.generateRepository(buf, class, outputPackage, reposPackage); err != nil {
			return err
		}
		reposFilename := path.Join(reposDir, fmt.Sprintf("%sRepository.kt", class.Name))
		if err := writeStringToFile(reposFilename, buf.String()); err != nil {
			return err
		}
	}

	return nil
}

func (k *JPAKotlinTpl) generateRepository(
	wr io.Writer,
	class *KotlinClass,
	entityPackageName string,
	reposPackageName string,
) error {
	// custom functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"fieldAnnotations": func(field *KotlinField) []string {
			return k.getFieldAnnotations(class, field)
		},
		"hasNext": func(field *KotlinField, fields []*KotlinField) bool {
			return field != fields[len(fields)-1]
		},
	}

	tplString := `{{"" -}}
package {{.ReposPackage}}

import {{.EntityPackage}}.*
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface {{.ClassName}}Repository : JpaRepository<{{.ClassName}}, {{.IdClassName}}>
`

	// parse template
	tmpl, err := template.New("jpaKotlinEntity").Funcs(funcMap).Parse(tplString)
	if err != nil {
		return err
	}

	// FIXME: duplicated
	var idClassName string
	pkFieldCount := len(class.PKFields)
	if pkFieldCount > 1 {
		idClassName = class.Name + "PK"
	} else {
		idClassName = class.PKFields[0].Type
	}

	data := struct {
		EntityPackage string
		ReposPackage  string
		ClassName     string
		IdClassName   string
	}{
		EntityPackage: entityPackageName,
		ReposPackage:  reposPackageName,
		ClassName:     class.Name,
		IdClassName:   idClassName,
	}
	return tmpl.Execute(wr, &data)
}
