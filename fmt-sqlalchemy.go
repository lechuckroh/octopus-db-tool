package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type SaClass struct {
	table        *Table
	Name         string
	Fields       []*SaField
	PKFields     []*SaField
	UniqueFields []*SaField
}

type SaField struct {
	Column       *Column
	Name         string
	OverrideName bool
	Type         string
	Imports      []string
}

type SqlAlchemy struct {
}

func NewSaClass(
	table *Table,
	output *Output,
	prefixMapper *PrefixMapper,
) *SaClass {
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

	fields := make([]*SaField, 0)
	pkFields := make([]*SaField, 0)
	uniqueFields := make([]*SaField, 0)
	for _, column := range table.Columns {
		field := NewSaField(column)
		fields = append(fields, field)

		if column.PrimaryKey {
			pkFields = append(pkFields, field)
		}
		if column.UniqueKey {
			uniqueFields = append(uniqueFields, field)
		}
	}

	return &SaClass{
		table:        table,
		Name:         className,
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
	}
}

func NewSaField(column *Column) *SaField {
	var fieldType string
	importSet := NewStringSet()

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case ColTypeString:
		fallthrough
	case ColTypeText:
		fieldType = "String"
	case ColTypeBoolean:
		fieldType = "Boolean"
	case ColTypeLong:
		fieldType = "BigInteger"
	case ColTypeInt:
		fieldType = "Integer"
	case ColTypeDecimal:
		fieldType = "Numeric"
	case ColTypeFloat:
		fieldType = "Float"
	case ColTypeDouble:
		fieldType = "Float"
	case ColTypeDateTime:
		fieldType = "DateTime"
	case ColTypeDate:
		fieldType = "Date"
	case ColTypeTime:
		fieldType = "Time"
	case ColTypeBlob:
		fieldType = "LargeBinary"
	default:
		if columnType == "bit" {
			if column.Size == 1 {
				fieldType = "Boolean"
				break
			}
		}
		fieldType = ""
	}

	if fieldType != "" {
		importSet.Add(fieldType)
	}

	fieldName, ok := ToLowerSnake(column.Name)

	// check python reserved words
	reservedWord := IsPythonReservedWord(fieldName)
	if reservedWord {
		fieldName = fieldName + "_"
	}

	return &SaField{
		Column:       column,
		Name:         fieldName,
		OverrideName: reservedWord || !ok,
		Type:         fieldType,
		Imports:      importSet.Slice(),
	}
}

func (sa *SqlAlchemy) mkdir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	log.Printf("[MKDIR] %s", dir)
	return dir, nil
}

func (sa *SqlAlchemy) Generate(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
) error {
	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)
	useUTC := output.GetBool(FlagUseUTC)

	// write to single file if extension is '.py'
	var outputDir string
	generateSingleFile := false
	if ext := strings.ToLower(filepath.Ext(output.FilePath)); ext == ".py" {
		outputDir = filepath.Dir(output.FilePath)
		generateSingleFile = true
	} else {
		outputDir = output.FilePath
	}

	if _, err := sa.mkdir(outputDir); err != nil {
		return err
	}

	indent := strings.Repeat(" ", 4)

	classes := make([]*SaClass, 0)
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		classes = append(classes, NewSaClass(table, output, prefixMapper))
	}

	// imports from sqlalchemy
	saImportSet := NewStringSet()
	saImportSet.Add("Column")
	// imports
	importSet := NewStringSet()

	// contents to write
	contents := make([]string, 0)

	for _, class := range classes {
		table := class.table

		classLines := make([]string, 0)
		appendLine := func(lines ...string) {
			classLines = append(classLines, lines...)
		}

		// unique
		uniqueCstName := table.Name + uniqueNameSuffix
		uniqueFieldNames := make([]string, 0)
		for _, field := range class.UniqueFields {
			uniqueFieldNames = append(uniqueFieldNames, Quote(field.Name, "'"))
		}

		// class
		appendLine("", "",
			fmt.Sprintf("class %s(Base):", class.Name),
			indent+fmt.Sprintf("__tablename__ = '%s'", table.Name),
		)

		if len(uniqueFieldNames) > 1 {
			appendLine(indent + fmt.Sprintf("UniqueConstraint(%s, name='%s')", strings.Join(uniqueFieldNames, ", "), uniqueCstName))
			saImportSet.Add("UniqueConstraint")
		}
		appendLine("")

		// fields
		for _, field := range class.Fields {
			column := field.Column

			// Column attributes
			attributes := make([]string, 0)

			if field.OverrideName {
				attributes = append(attributes, Quote(column.Name, "'"))
			}

			if column.Type == "string" && column.Size > 0 {
				attributes = append(attributes, fmt.Sprintf("%s(%d)", field.Type, column.Size))
			} else if column.Type == ColTypeDouble || column.Type == ColTypeFloat || column.Type == ColTypeDecimal {
				colAttrs := make([]string, 0)
				if column.Size > 0 {
					colAttrs = append(colAttrs, fmt.Sprintf("precision=%d", column.Size))
				}
				if column.Scale > 0 {
					colAttrs = append(colAttrs, fmt.Sprintf("scale=%d", column.Scale))
				}

				attributes = append(attributes, fmt.Sprintf("%s(%s)", field.Type, strings.Join(colAttrs, ", ")))
			} else {
				attributes = append(attributes, field.Type)
			}
			// PK
			if column.PrimaryKey {
				attributes = append(attributes, "primary_key=True")
			}
			// Unique
			if column.UniqueKey {
				attributes = append(attributes, "unique=True")
			}
			// auto_increment
			if column.AutoIncremental {
				attributes = append(attributes, "autoincrement=True")
			}
			// not null
			if !column.Nullable && !column.AutoIncremental {
				attributes = append(attributes, "nullable=False")
			}
			// audit column
			if column.Name == "created_at" {
				if useUTC {
					attributes = append(attributes, "default=datetime.utcnow")
				} else {
					attributes = append(attributes, "default=datetime.now")
				}
				importSet.Add("from datetime import datetime")
			}
			if column.Name == "updated_at" {
				if useUTC {
					attributes = append(attributes, "default=datetime.utcnow", "onupdate=datetime.utcnow")
				} else {
					attributes = append(attributes, "default=datetime.now", "onupdate=datetime.now")
				}
				importSet.Add("from datetime import datetime")
			}

			appendLine(indent + fmt.Sprintf("%s = Column(%s)", field.Name, strings.Join(attributes, ", ")))

			// import
			for _, imp := range field.Imports {
				saImportSet.Add(imp)
			}
		}

		if generateSingleFile {
			contents = append(contents, classLines...)
		} else {
			contents = append(contents, sa.getHeaderLines(importSet.Slice(), saImportSet.Slice())...)
			contents = append(contents, classLines...)
			contents = append(contents, "")

			outputFile := path.Join(outputDir, fmt.Sprintf("%s.py", table.Name))
			if err := WriteLinesToFile(outputFile, contents); err != nil {
				return err
			}

			// reset slice
			contents = make([]string, 0)
			saImportSet.Clear()
			importSet.Clear()
		}
	}

	// Write to single file
	if generateSingleFile {
		finalOutput := sa.getHeaderLines(importSet.Slice(), saImportSet.Slice())
		finalOutput = append(finalOutput, contents...)
		finalOutput = append(finalOutput, "")

		if err := WriteLinesToFile(output.FilePath, finalOutput); err != nil {
			return err
		}
	}

	return nil
}

func (sa *SqlAlchemy) getHeaderLines(imports []string, saImports []string) []string {
	lines := make([]string, 0)

	if len(imports) > 0 {
		lines = append(lines, imports...)
		lines = append(lines, "")
	}

	lines = append(lines,
		"from sqlalchemy import "+strings.Join(saImports, ", "),
		"from sqlalchemy.ext.declarative import declarative_base",
		"from sqlalchemy_repr import RepresentableBase",
		"",
		"Base = declarative_base(cls=RepresentableBase)",
	)

	return lines
}
