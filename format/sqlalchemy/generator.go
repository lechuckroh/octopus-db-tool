package sqlalchemy

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Option struct {
	PrefixMapper     *common.PrefixMapper
	TableFilter      octopus.TableFilterFn
	RemovePrefixes   []string
	UniqueNameSuffix string
	UseUTC           bool
}

type Generator struct {
	schema *octopus.Schema
	option *Option
}

func (c *Generator) mkdir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	log.Printf("[MKDIR] %s", dir)
	return dir, nil
}

func (c *Generator) Generate(outputPath string) error {
	option := c.option
	uniqueNameSuffix := option.UniqueNameSuffix
	useUTC := option.UseUTC

	// write to single file if extension is '.py'
	var outputDir string
	generateSingleFile := false
	if ext := strings.ToLower(filepath.Ext(outputPath)); ext == ".py" {
		outputDir = filepath.Dir(outputPath)
		generateSingleFile = true
	} else {
		outputDir = outputPath
	}

	if _, err := c.mkdir(outputDir); err != nil {
		return err
	}

	indent := strings.Repeat(" ", 4)

	classes := make([]*SaClass, 0)
	for _, table := range c.schema.Tables {
		// filter table
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}
		classes = append(classes, NewSaClass(table, option))
	}

	// imports from sqlalchemy
	saImportSet := util.NewStringSet()
	saImportSet.Add("Column")
	// imports
	importSet := util.NewStringSet()

	// contents to write
	contents := make([]string, 0)
	useTZDateTime := false

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
			uniqueFieldNames = append(uniqueFieldNames, util.Quote(field.Name, "'"))
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
			lcColumnType := strings.ToLower(column.Type)

			// Column attributes
			attributes := make([]string, 0)

			if field.OverrideName {
				attributes = append(attributes, util.Quote(column.Name, "'"))
			}

			if (lcColumnType == octopus.ColTypeVarchar || lcColumnType == octopus.ColTypeChar) && column.Size > 0 {
				attributes = append(attributes, fmt.Sprintf("%s(%d)", field.Type, column.Size))
			} else if lcColumnType == octopus.ColTypeDouble || lcColumnType == octopus.ColTypeFloat || lcColumnType == octopus.ColTypeDecimal {
				colAttrs := make([]string, 0)
				if column.Size > 0 {
					colAttrs = append(colAttrs, fmt.Sprintf("precision=%d", column.Size))
				}
				if column.Scale > 0 {
					colAttrs = append(colAttrs, fmt.Sprintf("scale=%d", column.Scale))
				}

				attributes = append(attributes, fmt.Sprintf("%s(%s)", field.Type, strings.Join(colAttrs, ", ")))
			} else if lcColumnType == octopus.ColTypeDateTime {
				if column.Name == "created_at" {
					if useUTC {
						attributes = append(attributes, field.Type, "default=datetime.utcnow")
					} else {
						attributes = append(attributes, field.Type, "default=datetime.now")
					}
				} else if column.Name == "updated_at" {
					if useUTC {
						attributes = append(attributes, field.Type, "default=datetime.utcnow", "onupdate=datetime.utcnow")
					} else {
						attributes = append(attributes, field.Type, "default=datetime.now", "onupdate=datetime.now")
					}
				} else {
					if useUTC {
						attributes = append(attributes, "TZDateTime")
						useTZDateTime = true
						saImportSet.Add("TypeDecorator")
						importSet.Add("from datetime import timezone")
					} else {
						attributes = append(attributes, field.Type)
					}
				}
				importSet.Add("from datetime import datetime")
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
			if column.NotNull && !column.AutoIncremental {
				attributes = append(attributes, "nullable=False")
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
			contents = append(contents, c.getHeaderLines(importSet.Slice(), saImportSet.Slice())...)
			contents = append(contents, classLines...)
			contents = append(contents, "")

			outputFile := path.Join(outputDir, fmt.Sprintf("%s.py", table.Name))
			if err := util.WriteLinesToFile(outputFile, contents); err != nil {
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
		finalOutput := c.getHeaderLines(importSet.Slice(), saImportSet.Slice())
		if useTZDateTime {
			finalOutput = append(finalOutput, c.getTZDateTimeLines()...)
		}
		finalOutput = append(finalOutput, contents...)
		finalOutput = append(finalOutput, "")

		if err := util.WriteLinesToFile(outputPath, finalOutput); err != nil {
			return err
		}
	}

	return nil
}

func (c *Generator) GenerateClass(wr io.Writer, class *SaClass) error {
	// TODO: implement
	return nil
}

func (c *Generator) getTZDateTimeLines() []string {
	return []string{
		"",
		"",
		"class TZDateTime(TypeDecorator):",
		"    impl = DateTime",
		"",
		"    def process_bind_param(self, value, dialect):",
		"        if value is not None:",
		"            if not value.tzinfo:",
		"                raise TypeError(\"tzinfo is required\")",
		"            value = value.astimezone(timezone.utc).replace(tzinfo=None)",
		"        return value",
		"",
		"    def process_result_value(self, value, dialect):",
		"        if value is not None:",
		"            value = value.replace(tzinfo=timezone.utc)",
		"        return value",
	}
}

func (c *Generator) getHeaderLines(imports []string, saImports []string) []string {
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
