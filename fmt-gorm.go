package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type GormClass struct {
	table        *Table
	Name         string
	EmbedModel   bool
	Fields       []*GormField
	PKFields     []*GormField
	UniqueFields []*GormField
}

type GormField struct {
	Column       *Column
	Name         string
	Type         string
	OverrideName bool
	Imports      []string
}

type Gorm struct {
}

func getGormModelColumns() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at"}
}

func NewGormClass(
	table *Table,
	output *Output,
	prefixMapper *PrefixMapper,
) *GormClass {
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
	modelColumnSet := NewStringSet(getGormModelColumns()...)
	for _, column := range table.Columns {
		field := NewGormField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
		if column.UniqueKey {
			uniqueFields = append(uniqueFields, field)
		}

		modelColumnSet.Remove(column.Name)
	}

	return &GormClass{
		table:        table,
		Name:         className,
		EmbedModel:   modelColumnSet.Size() == 0,
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
	}
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

func (g *Gorm) mkdir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	log.Printf("[MKDIR] %s", dir)
	return dir, nil
}

func (g *Gorm) Generate(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
) error {
	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)
	gormModel := output.Get(FlagGormModel)
	if gormModel == "" {
		gormModel = "gorm.Model"
	}
	pkg := output.Get(FlagPackage)
	if pkg == "" {
		pkg = "main"
	}

	// write to single file if extension is '.go'
	var outputDir string
	generateSingleFile := false
	if ext := strings.ToLower(filepath.Ext(output.FilePath)); ext == ".go" {
		outputDir = filepath.Dir(output.FilePath)
		generateSingleFile = true
	} else {
		outputDir = output.FilePath
	}

	if _, err := g.mkdir(outputDir); err != nil {
		return err
	}

	indent := strings.Repeat(" ", 4)

	classes := make([]*GormClass, 0)
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		classes = append(classes, NewGormClass(table, output, prefixMapper))
	}

	// imports
	importSet := NewStringSet()

	// contents to write
	contents := make([]string, 0)

	// embedded model columns
	embeddedModelColumns := NewStringSet(getGormModelColumns()...)

	for _, class := range classes {
		table := class.table

		classLines := make([]string, 0)
		appendLine := func(lines ...string) {
			classLines = append(classLines, lines...)
		}

		// embedded model
		if class.EmbedModel {
			importSet.Add("github.com/jinzhu/gorm")
		}

		// unique
		uniqueCstName := ""
		uniqueFieldNames := make([]string, 0)
		for _, field := range class.UniqueFields {
			uniqueFieldNames = append(uniqueFieldNames, Quote(field.Name, "'"))
		}
		if len(uniqueFieldNames) > 1 {
			uniqueCstName = table.Name + uniqueNameSuffix
		}

		// struct
		appendLine("", "",
			fmt.Sprintf("type %s struct {", class.Name),
		)
		if class.EmbedModel {
			appendLine(indent + "gorm.Model")
		}

		// fields
		for _, field := range class.Fields {
			column := field.Column

			// embedded model column
			if class.EmbedModel && embeddedModelColumns.Contains(column.Name){
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
			if len(gormTags) == 0 {
				appendLine(indent + fmt.Sprintf("%s %s", field.Name, field.Type))
			} else {
				tag := fmt.Sprintf("`gorm:\"%s\"`", strings.Join(gormTags, ";"))
				appendLine(indent + fmt.Sprintf("%s %s %s", field.Name, field.Type, tag))
			}

			// import
			for _, imp := range field.Imports {
				importSet.Add(imp)
			}
		}
		classLines = append(classLines,
			"}",
			fmt.Sprintf("func (c *%s) TableName() string { return \"%s\" }", class.Name, table.Name))

		if generateSingleFile {
			contents = append(contents, classLines...)
		} else {
			contents = append(contents, fmt.Sprintf("package %s", pkg))
			contents = append(contents, g.getHeaderLines(indent, importSet.Slice())...)
			contents = append(contents, classLines...)
			contents = append(contents, "")

			outputFile := path.Join(outputDir, fmt.Sprintf("%s.go", table.Name))
			if err := WriteLinesToFile(outputFile, contents); err != nil {
				return err
			}

			// reset slice
			contents = make([]string, 0)
			importSet.Clear()
		}
	}

	// Write to single file
	if generateSingleFile {
		finalOutput := []string{
			fmt.Sprintf("package %s", pkg),
			"",
		}
		finalOutput = append(finalOutput, g.getHeaderLines(indent, importSet.Slice())...)
		finalOutput = append(finalOutput, contents...)
		finalOutput = append(finalOutput, "")

		if err := WriteLinesToFile(output.FilePath, finalOutput); err != nil {
			return err
		}
	}

	return nil
}

func (g *Gorm) getHeaderLines(indent string, imports []string) []string {
	lines := make([]string, 0)

	if len(imports) > 0 {
		lines = append(lines, "import (")
		for _, imp := range imports {
			lines = append(lines, indent+Quote(imp, "\""))
		}
		lines = append(lines, ")")
	}

	return lines
}
