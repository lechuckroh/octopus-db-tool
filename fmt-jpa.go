package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type KotlinClass struct {
	table        *Table
	Name         string
	Fields       []*KotlinField
	PKFields     []*KotlinField
	UniqueFields []*KotlinField
}

type KotlinField struct {
	Column       *Column
	Name         string
	Type         string
	Imports      []string
	DefaultValue string
}

type JPAKotlin struct {
}

func NewKotlinClass(table *Table, output *GenOutput) *KotlinClass {
	className := table.ClassName
	if className == "" {
		tableName := table.Name
		for _, prefix := range output.PrefixesToRemove {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		className = strcase.ToCamel(tableName)
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
	case "varchar":
		fallthrough
	case "char":
		fallthrough
	case "string":
		fallthrough
	case "text":
		fieldType = "String"
		if !nullable {
			defaultValue = "\"\""
		}
	case "bool":
		fallthrough
	case "boolean":
		fieldType = "Boolean"
		if !nullable {
			defaultValue = "false"
		}
	case "bigint":
		fallthrough
	case "long":
		fieldType = "Long"
		if !nullable {
			defaultValue = "0L"
		}
	case "int":
		fallthrough
	case "integer":
		fallthrough
	case "smallint":
		fieldType = "Int"
		if !nullable {
			defaultValue = "0"
		}
	case "float":
		fieldType = "Float"
		if !nullable {
			defaultValue = "0.0F"
		}
	case "number":
		fallthrough
	case "double":
		fieldType = "Double"
		if !nullable {
			defaultValue = "0.0"
		}
	case "datetime":
		fieldType = "LocalDateTime"
		importSet.Add("java.time.LocalDateTime")
		if !nullable {
			defaultValue = "LocalDateTime.now()"
		}
	case "date":
		fieldType = "LocalDate"
		importSet.Add("java.time.LocalDate")
		if !nullable {
			defaultValue = "0.0"
		}
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

	return &KotlinField{
		Column:       column,
		Name:         strcase.ToLowerCamel(column.Name),
		Type:         fieldType,
		DefaultValue: defaultValue,
		Imports:      importSet.Slice(),
	}
}

func (k *JPAKotlin) Generate(schema *Schema, output *GenOutput, useDataClass bool) error {
	// Create package directory
	entityDir := path.Join(append([]string{output.Directory}, strings.Split(output.Package, ".")...)...)
	if err := os.MkdirAll(entityDir, 0777); err != nil {
		return err
	}
	log.Printf("[MKDIR] %s", entityDir)

	// Create repos package directory
	reposDir := ""
	if output.ReposPackage != "" {
		reposDir = path.Join(append([]string{output.Directory}, strings.Split(output.ReposPackage, ".")...)...)
		if err := os.MkdirAll(reposDir, 0777); err != nil {
			return err
		}
		log.Printf("[MKDIR] %s", reposDir)
	}

	if !useDataClass {
		// Generate AbstractJpaPersistable.kt
		if err := k.generateAbstractJpaPersistable(entityDir); err != nil {
			return err
		}
	}

	indent := "    "
	for _, table := range schema.Tables {
		class := NewKotlinClass(table, output)
		var idClass string
		pkFieldCount := len(class.PKFields)

		contents := make([]string, 0)
		classLines := make([]string, 0)
		appendLine := func(line string) {
			classLines = append(classLines, line)
		}

		// package
		if output.Package != "" {
			contents = append(contents, fmt.Sprintf("package %s", output.Package), "")
		}
		importSet := NewStringSet()
		importSet.Add("javax.persistence.*")

		// class
		appendLine("")
		appendLine("@Entity")
		appendLine(fmt.Sprintf("@Table(name = \"%s\")", table.Name))
		if pkFieldCount > 1 {
			idClass = class.Name + "PK"
			appendLine(fmt.Sprintf("@IdClass(%s::class)", idClass))
		}

		classDef := fmt.Sprintf("class %s(", class.Name)
		if useDataClass {
			appendLine("data " + classDef)
		} else {
			appendLine(classDef)
		}

		// fields
		fieldCount := len(class.Fields)
		for i, field := range class.Fields {
			column := field.Column
			if column.PrimaryKey {
				appendLine(indent + "@Id")
				if idClass == "" {
					idClass = field.Type
				}
			}
			if column.AutoIncremental {
				appendLine(indent + "@GeneratedValue(strategy = GenerationType.IDENTITY)")
			}
			if column.Type == "text" {
				appendLine(indent + "@Type(type = \"text\")")
				importSet.Add("org.hibernate.annotations.Type")
			}

			// @Column attributes
			attributes := make([]string, 0)
			if !column.Nullable {
				attributes = append(attributes, "nullable = false")
			}
			if field.Type == "String" && column.Size > 0 {
				attributes = append(attributes, fmt.Sprintf("length = %d", column.Size))
			}
			if len(attributes) > 0 {
				appendLine(indent + fmt.Sprintf("@Column(%s)", strings.Join(attributes, ", ")))
			}

			// @CreationTimestamp
			if column.Type == "datetime" && field.Name == "createdAt" {
				appendLine(indent + "@CreationTimestamp")
				importSet.Add("org.hibernate.annotations.CreationTimestamp")
			}
			// @UpdateTimestamp
			if column.Type == "datetime" && field.Name == "updatedAt" {
				appendLine(indent + "@UpdateTimestamp")
				importSet.Add("org.hibernate.annotations.UpdateTimestamp")
			}

			line := fmt.Sprintf("var %s: %s", field.Name, field.Type)

			// set default value
			if field.DefaultValue != "" {
				line = line + " = " + field.DefaultValue
			}

			if i < fieldCount-1 {
				appendLine(indent + line + ",")
				appendLine("")
			} else {
				appendLine(indent + line)
			}

			// import
			importSet.AddAll(field.Imports)
		}

		if useDataClass {
			appendLine(")")
		} else {
			appendLine("")
			appendLine(fmt.Sprintf(") : AbstractJpaPersistable<%s>()", idClass))
		}
		appendLine("")

		// Composite Key
		idClassLines := make([]string, 0)
		if pkFieldCount > 1 {
			addLine := func(s string) { idClassLines = append(idClassLines, s) }

			importSet.Add("java.io.Serializable")
			addLine(fmt.Sprintf("data class %s(", idClass))
			for i, pkField := range class.PKFields {
				line := indent + fmt.Sprintf("var %s: %s = %s", pkField.Name, pkField.Type, pkField.DefaultValue)
				if i < pkFieldCount-1 {
					line = line + ","
				}
				addLine(line)
			}
			addLine("): Serializable")
			addLine("")
		}

		// contents
		for _, imp := range importSet.Slice() {
			contents = append(contents, "import "+imp)
		}
		contents = append(contents, classLines...)
		contents = append(contents, idClassLines...)

		// Write file
		outputFile := path.Join(entityDir, fmt.Sprintf("%s.kt", class.Name))
		if err := ioutil.WriteFile(outputFile, []byte(strings.Join(contents, "\n")), 0644); err != nil {
			return err
		}
		log.Printf("[WRITE] %s", outputFile)

		// Write Repos
		if output.ReposPackage != "" {
			reposClassName := fmt.Sprintf("%sRepository", class.Name)
			file := path.Join(reposDir, reposClassName + ".kt")
			lines := []string{
				"package " + output.ReposPackage,
				"",
				"import " + output.Package + ".*",
				"import org.springframework.data.jpa.repository.JpaRepository",
				"import org.springframework.data.rest.core.annotation.RepositoryRestResource",
				"",
				"@RepositoryRestResource",
				fmt.Sprintf("interface %s : JpaRepository<%s, %s>", reposClassName, class.Name, idClass),
				"",
			}
			if err := ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644); err != nil {
				return err
			}
			log.Printf("[WRITE] %s", file)
		}
	}

	return nil
}

func (k *JPAKotlin) generateAbstractJpaPersistable(outputDir string) error {
	filename := path.Join(outputDir, "AbstractJpaPersistable.kt")
	data := `package kstec.sp.api.entity

import org.springframework.data.util.ProxyUtils
import java.io.Serializable
import javax.persistence.GeneratedValue
import javax.persistence.Id
import javax.persistence.MappedSuperclass

@MappedSuperclass
abstract class AbstractJpaPersistable<T : Serializable> {
    companion object {
        private val serialVersionUID = -5554308939380869754L
    }

    @Id
    @GeneratedValue
    private var id: T? = null

    fun getId(): T? {
        return id
    }

    override fun equals(other: Any?): Boolean {
        other ?: return false

        if (this === other) return true

        if (javaClass != ProxyUtils.getUserClass(other)) return false

        other as AbstractJpaPersistable<*>

        return if (null == this.getId()) false else this.getId() == other.getId()
    }

    override fun hashCode(): Int {
        return 31
    }

    override fun toString() = "Entity of type ${this.javaClass.name} with id: $id"
}
`
	return ioutil.WriteFile(filename, []byte(data), 0644)
}
