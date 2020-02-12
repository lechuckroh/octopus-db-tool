package main

import (
	"fmt"
	"github.com/gertd/go-pluralize"
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

func NewKotlinClass(
	table *Table,
	output *Output,
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

	return &KotlinField{
		Column:       column,
		Name:         strcase.ToLowerCamel(column.Name),
		Type:         fieldType,
		DefaultValue: defaultValue,
		Imports:      importSet.Slice(),
	}
}

func (k *JPAKotlin) mkdir(basedir, pkgName string) (string, error) {
	if pkgName == "" {
		return "", nil
	}
	dir := path.Join(append([]string{basedir}, strings.Split(pkgName, ".")...)...)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}
	log.Printf("[MKDIR] %s", dir)
	return dir, nil
}

func (k *JPAKotlin) Generate(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
	prefixMapper *PrefixMapper,
	useDataClass bool,
) error {
	outputPackage := output.Get(FlagPackage)
	reposPackage := output.Get(FlagReposPackage)
	graphqlPackage := output.Get(FlagGraphqlPackage)
	relation := output.Get(FlagRelation)
	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)

	entityDir, err := k.mkdir(output.FilePath, outputPackage)
	if err != nil {
		return err
	}
	reposDir, err := k.mkdir(output.FilePath, reposPackage)
	if err != nil {
		return err
	}
	graphqlDir, err := k.mkdir(output.FilePath, graphqlPackage)
	if err != nil {
		return err
	}

	if !useDataClass {
		// Generate AbstractJpaPersistable.kt
		if err := k.generateAbstractJpaPersistable(entityDir, outputPackage); err != nil {
			return err
		}
	}

	indent := "    "

	classes := make([]*KotlinClass, 0)
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}
		classes = append(classes, NewKotlinClass(table, output, prefixMapper))
	}

	getClassNameByTable := func(table string) string {
		for _, cls := range classes {
			if cls.table.Name == table {
				return cls.Name
			}
		}
		return ""
	}

	for _, class := range classes {
		table := class.table

		var idClass string
		pkFieldCount := len(class.PKFields)

		contents := make([]string, 0)
		classLines := make([]string, 0)
		appendLine := func(line string) {
			classLines = append(classLines, line)
		}

		// package
		if outputPackage != "" {
			contents = append(contents, fmt.Sprintf("package %s", outputPackage), "")
		}
		importSet := NewStringSet()
		importSet.Add("javax.persistence.*")

		// unique
		uniqueCstName := table.Name + uniqueNameSuffix
		uniqueFieldNames := make([]string, 0)
		for _, field := range class.UniqueFields {
			uniqueFieldNames = append(uniqueFieldNames, Quote(field.Name, "\""))
		}

		// class
		appendLine("")
		appendLine("@Entity")
		if len(uniqueFieldNames) == 0 {
			appendLine(fmt.Sprintf("@Table(name = \"%s\")", table.Name))
		} else {
			appendLine(fmt.Sprintf("@Table(name = \"%s\", uniqueConstraints = [\n    UniqueConstraint(name = \"%s\", columnNames = [%s])\n])",
				table.Name, uniqueCstName, strings.Join(uniqueFieldNames, ", ")))
		}
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
				appendLine(indent + "@GeneratedValue(strategy = GenerationType.AUTO)")
			}
			if column.Type == "text" {
				appendLine(indent + "@Lob")
			}

			// @VRelation
			if relation == "VRelation" {
				if ref := column.Ref; ref != nil {
					targetClassName := getClassNameByTable(ref.Table)
					if len(targetClassName) == 0 {
						log.Fatalf("Relation not found. %s::%s -> %s",
							class.Name, field.Name, ref.Table)
					}

					appendLine(indent +
						fmt.Sprintf("@VRelation(cls = \"%s\", field = \"%s\")",
							targetClassName,
							strcase.ToLowerCamel(ref.Column)))
				}
			}

			// @Column attributes
			attributes := make([]string, 0)
			if !column.Nullable {
				attributes = append(attributes, "nullable = false")
			}
			if column.Type == "string" && column.Size > 0 {
				attributes = append(attributes, fmt.Sprintf("length = %d", column.Size))
			}
			if column.Type == ColTypeDouble || column.Type == ColTypeFloat || column.Type == ColTypeDecimal {
				if column.Size > 0 {
					attributes = append(attributes, fmt.Sprintf("precision = %d", column.Size))
				}
				if column.Scale > 0 {
					attributes = append(attributes, fmt.Sprintf("scale = %d", column.Scale))
				}
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
		if reposDir != "" {
			reposClassName := fmt.Sprintf("%sRepository", class.Name)
			lines := []string{
				"package " + reposPackage,
				"",
				"import " + outputPackage + ".*",
				"import org.springframework.data.repository.PagingAndSortingRepository",
				"import org.springframework.stereotype.Repository",
				"",
				"@Repository",
				fmt.Sprintf("interface %s : PagingAndSortingRepository<%s, %s>", reposClassName, class.Name, idClass),
				"",
			}
			if err := k.writeLines(path.Join(reposDir, reposClassName+".kt"), lines); err != nil {
				return err
			}
		}
	}

	// write graphql
	if graphqlDir != "" {
		contents := []string{
			"package " + graphqlPackage,
			"",
			"import com.coxautodev.graphql.tools.GraphQLQueryResolver",
			"import org.springframework.stereotype.Component",
			"import java.util.*",
			"import " + outputPackage + ".*",
			"import " + reposPackage + ".*",
			"",
			"@Component",
			"class Query(",
		}
		ctorArgs := make([]string, 0)
		for _, class := range classes {
			ctorArgs = append(ctorArgs,
				fmt.Sprintf("        private val %sRepos: %sRepository", strcase.ToLowerCamel(class.Name), class.Name))
		}
		contents = append(contents, strings.Join(ctorArgs, ",\n"))
		contents = append(contents, ") : GraphQLQueryResolver {")

		client := pluralize.NewClient()
		for _, class := range classes {
			lowerClassName := strcase.ToLowerCamel(class.Name)

			contents = append(contents, "")
			contents = append(contents, fmt.Sprintf("    fun %s(): Iterable<%s> {", client.Plural(lowerClassName), class.Name))
			contents = append(contents, fmt.Sprintf("        return %sRepos.findAll()", lowerClassName))
			contents = append(contents, "    }")
		}
		contents = append(contents, "}")
		if err := k.writeLines(path.Join(graphqlDir, "Query.kt"), contents); err != nil {
			return err
		}
	}

	return nil
}

func (k *JPAKotlin) writeLines(filename string, lines []string) error {
	if err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", filename)
	return nil
}

func (k *JPAKotlin) generateAbstractJpaPersistable(outputDir string, packageName string) error {
	filename := path.Join(outputDir, "AbstractJpaPersistable.kt")
	data := fmt.Sprintf(`package %s

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
`, packageName)
	return ioutil.WriteFile(filename, []byte(data), 0644)
}
