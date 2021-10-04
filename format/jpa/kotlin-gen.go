package jpa

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"log"
	"path"
	"strings"
	"text/template"
)

type KtOption struct {
	AnnoMapper       *common.AnnotationMapper
	PrefixMapper     *common.PrefixMapper
	TableFilter      octopus.TableFilterFn
	IdEntity         string
	Package          string
	Relation         string
	RemovePrefixes   []string
	ReposPackage     string
	UniqueNameSuffix string
}

type KtGenerator struct {
	schema  *octopus.Schema
	classes []*KotlinClass
	option  *KtOption
}

func NewKtGenerator(
	schema *octopus.Schema,
	option *KtOption,
) *KtGenerator {
	// populate KotlinClass
	var classes []*KotlinClass
	for _, table := range schema.Tables {
		// TODO: move to Generate() function
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}
		classes = append(classes, NewKotlinClass(table, option))
	}

	return &KtGenerator{
		schema:  schema,
		classes: classes,
		option:  option,
	}
}

func (c *KtGenerator) getFieldAnnotations(class *KotlinClass, field *KotlinField) []string {
	var annotations []string

	column := field.Column
	var colAttrs []string
	if column.PrimaryKey {
		annotations = append(annotations, "@Id")
	}
	if column.AutoIncremental {
		annotations = append(annotations, "@GeneratedValue(strategy = GenerationType.AUTO)")
	}
	if octopus.IsColTypeClob(column.Type) {
		annotations = append(annotations, "@Lob")
	}

	// @VRelation
	if c.option.Relation == "VRelation" {
		if ref := column.Ref; ref != nil {
			targetClassName := c.getClassNameByTableName(ref.Table)
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
	if column.NotNull {
		colAttrs = append(colAttrs, "nullable = false")
	}
	if octopus.IsColTypeString(column.Type) && column.Size > 0 {
		colAttrs = append(colAttrs, fmt.Sprintf("length = %d", column.Size))
	}
	if octopus.IsColTypeDecimal(column.Type) {
		if column.Size > 0 {
			colAttrs = append(colAttrs, fmt.Sprintf("precision = %d", column.Size))
		}
		if column.Scale > 0 {
			colAttrs = append(colAttrs, fmt.Sprintf("scale = %d", column.Scale))
		}
	}
	// @CreationTimestamp
	if column.Type == octopus.ColTypeDateTime && field.Name == "createdAt" {
		annotations = append(annotations, "@CreationTimestamp")
		colAttrs = append(colAttrs, "updatable = false")
	}
	// @UpdateTimestamp
	if column.Type == octopus.ColTypeDateTime && field.Name == "updatedAt" {
		annotations = append(annotations, "@UpdateTimestamp")
	}
	if len(colAttrs) > 0 {
		annotations = append(annotations, fmt.Sprintf("@Column(%s)", strings.Join(colAttrs, ", ")))
	}
	return annotations
}

type KotlinTplData struct {
	Package              string
	Class                *KotlinClass
	SuperClass           string
	Annotations          []string
	Table                *octopus.Table
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
func (c *KtGenerator) GenerateEntityClass(
	wr io.Writer,
	class *KotlinClass,
) error {
	// custom functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"fieldAnnotations": func(field *KotlinField) []string {
			return c.getFieldAnnotations(class, field)
		},
		"hasNext": func(field *KotlinField, fields []*KotlinField) bool {
			return field != fields[len(fields)-1]
		},
	}

	// template
	tplText := `{{"" -}}
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
    {{- if .Column.NotNull}}
		{{- if eq . $.IdEntityField}}
        override var {{.Name}}: {{.Type}} = {{.DefaultValue}},
		{{- else}}
        var {{.Name}}: {{.Type}} = {{.DefaultValue}}{{if hasNext . $.Class.Fields}},{{end}}
		{{- end}}
    {{- else}}
        var {{.Name}}: {{.Type}}{{if hasNext . $.Class.Fields}},{{end}}
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
	tpl, err := util.NewTemplate("jpaKotlinEntity", tplText, funcMap)
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
	idEntityInterfaceName := c.option.IdEntity
	if idEntityInterfaceName != "" && pkFieldCount == 1 && class.PKFields[0].Name == "id" {
		idEntityField = class.PKFields[0]
	}

	// annotations
	annotations := c.option.AnnoMapper.GetAnnotations(class.table.Group)

	// super class
	superClass := ""
	if idEntityField != nil {
		superClass = fmt.Sprintf("%s<%s>", idEntityInterfaceName, idEntityField.Type)
	}

	// unique constraint
	var uniqueFieldNames []string
	for _, field := range class.UniqueFields {
		uniqueFieldNames = append(uniqueFieldNames, util.Quote(field.Name, "\""))
	}

	// imports
	importSet := util.NewStringSet()
	javaImportSet := util.NewStringSet()
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
		if column.Type == octopus.ColTypeDateTime && field.Name == "createdAt" {
			importSet.Add("org.hibernate.annotations.CreationTimestamp")
		}
		// @UpdateTimestamp
		if column.Type == octopus.ColTypeDateTime && field.Name == "updatedAt" {
			importSet.Add("org.hibernate.annotations.UpdateTimestamp")
		}
	}

	// populate template data
	data := KotlinTplData{
		Package:              c.option.Package,
		Class:                class,
		SuperClass:           superClass,
		Annotations:          annotations,
		Table:                class.table,
		IdEntityField:        idEntityField,
		IdClassName:          idClassName,
		UniqueConstraintName: class.table.Name + c.option.UniqueNameSuffix,
		UniqueFieldNames:     uniqueFieldNames,
		HasUniqueFields:      len(uniqueFieldNames) > 0,
		Imports:              importSet.Slice(),
		JavaImports:          javaImportSet.Slice(),
	}

	return tpl.Execute(wr, &data)
}

func (c *KtGenerator) getClassNameByTableName(tableName string) string {
	for _, cls := range c.classes {
		if cls.table.Name == tableName {
			return cls.Name
		}
	}
	return ""
}

func (c *KtGenerator) Generate(outputPath string) error {
	if err := c.generateEntities(outputPath); err != nil {
		return err
	}

	if c.option.ReposPackage != "" {
		if err := c.generateRepositories(outputPath); err != nil {
			return err
		}
	}

	return nil
}

func (c *KtGenerator) generateEntities(outputPath string) error {
	outputPackage := c.option.Package

	entityDir, err := util.MkdirPackage(outputPath, outputPackage)
	if err != nil {
		return err
	}
	classes := c.classes

	// generate entity
	for _, class := range classes {
		// write entity class
		buf := new(bytes.Buffer)
		if err := c.GenerateEntityClass(buf, class); err != nil {
			return err
		}
		filename := path.Join(entityDir, fmt.Sprintf("%s.kt", class.Name))
		if err := util.WriteStringToFile(filename, buf.String()); err != nil {
			return err
		}
	}

	return nil
}


func (c *KtGenerator) generateRepositories(outputPath string) error {
	outputPackage := c.option.Package
	reposPackage := c.option.ReposPackage

	classes := c.classes

	reposDir, err := util.MkdirPackage(outputPath, reposPackage)
	if err != nil {
		return err
	}

	for _, class := range classes {
		buf := new(bytes.Buffer)
		if err := c.generateRepository(buf, class, outputPackage, reposPackage); err != nil {
			return err
		}
		reposFilename := path.Join(reposDir, fmt.Sprintf("%sRepository.kt", class.Name))
		if err := util.WriteStringToFile(reposFilename, buf.String()); err != nil {
			return err
		}
	}

	return nil
}

func (c *KtGenerator) generateRepository(
	wr io.Writer,
	class *KotlinClass,
	entityPackageName string,
	reposPackageName string,
) error {
	// custom functions
	funcMap := template.FuncMap{
		"join": strings.Join,
		"fieldAnnotations": func(field *KotlinField) []string {
			return c.getFieldAnnotations(class, field)
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
