package jpa

import (
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"strings"
)

type KotlinClass struct {
	table        *octopus.Table
	Name         string
	Annotations  []string
	Fields       []*KotlinField
	PKFields     []*KotlinField
	UniqueFields []*KotlinField
}

type KotlinField struct {
	Column       *octopus.Column
	Name         string
	OverrideName bool
	Type         string
	Imports      []string
	DefaultValue string
}

func NewKotlinClass(
	table *octopus.Table,
	option *KtOption,
) *KotlinClass {
	className := table.ClassName
	if className == "" {
		tableName := table.Name
		for _, prefix := range option.RemovePrefixes {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		className = strcase.ToCamel(tableName)

		if prefix := option.PrefixMapper.GetPrefix(table.Group); prefix != "" {
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
		Annotations:  option.AnnoMapper.GetAnnotations(table.Group),
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
	}
}

func NewKotlinField(column *octopus.Column) *KotlinField {
	var fieldType string
	var defaultValue string
	nullable := column.Nullable

	importSet := util.NewStringSet()

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case octopus.ColTypeString:
		fallthrough
	case octopus.ColTypeText:
		fieldType = "String"
		if !nullable {
			defaultValue = "\"\""
		}
	case octopus.ColTypeBoolean:
		fieldType = "Boolean"
		if !nullable {
			defaultValue = "false"
		}
	case octopus.ColTypeLong:
		fieldType = "Long"
		if !nullable {
			defaultValue = "0L"
		}
	case octopus.ColTypeInt:
		fieldType = "Int"
		if !nullable {
			defaultValue = "0"
		}
	case octopus.ColTypeDecimal:
		fieldType = "BigDecimal"
		importSet.Add("java.math.BigDecimal")
		if !nullable {
			defaultValue = "BigDecimal.ZERO"
		}
	case octopus.ColTypeFloat:
		fieldType = "Float"
		if !nullable {
			defaultValue = "0.0F"
		}
	case octopus.ColTypeDouble:
		fieldType = "Double"
		if !nullable {
			defaultValue = "0.0"
		}
	case octopus.ColTypeDateTime:
		fieldType = "Timestamp"
		importSet.Add("java.sql.Timestamp")
		if !nullable {
			defaultValue = "Timestamp(System.currentTimeMillis())"
		}
	case octopus.ColTypeDate:
		fieldType = "LocalDate"
		importSet.Add("java.time.LocalDate")
		if !nullable {
			defaultValue = "LocalDate.now()"
		}
	case octopus.ColTypeTime:
		fieldType = "LocalTime"
		importSet.Add("java.time.LocalTime")
		if !nullable {
			defaultValue = "LocalTime.now()"
		}
	case octopus.ColTypeBlob:
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

	fieldName, ok := util.ToLowerCamel(column.Name)

	return &KotlinField{
		Column:       column,
		Name:         fieldName,
		OverrideName: !ok,
		Type:         fieldType,
		DefaultValue: defaultValue,
		Imports:      importSet.Slice(),
	}
}
