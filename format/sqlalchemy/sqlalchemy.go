package sqlalchemy

import (
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"strings"
)

type SaClass struct {
	table        *octopus.Table
	Name         string
	Fields       []*SaField
	PKFields     []*SaField
	UniqueFields []*SaField
}

type SaField struct {
	Column       *octopus.Column
	Name         string
	OverrideName bool
	Type         string
	Imports      []string
}

func NewSaClass(
	table *octopus.Table,
	option *Option,
) *SaClass {
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

	var fields []*SaField
	var pkFields []*SaField
	var uniqueFields []*SaField
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

func NewSaField(column *octopus.Column) *SaField {
	var fieldType string
	importSet := util.NewStringSet()

	columnType := strings.ToLower(column.Type)
	switch columnType {
	case octopus.ColTypeChar:
		fallthrough
	case octopus.ColTypeVarchar:
		fallthrough
	case octopus.ColTypeText8:
		fallthrough
	case octopus.ColTypeText16:
		fallthrough
	case octopus.ColTypeText24:
		fallthrough
	case octopus.ColTypeText32:
		fieldType = "String"
	case octopus.ColTypeBoolean:
		fieldType = "Boolean"
	case octopus.ColTypeInt8:
		fallthrough
	case octopus.ColTypeInt16:
		fallthrough
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		fieldType = "Integer"
	case octopus.ColTypeInt64:
		fieldType = "BigInteger"
	case octopus.ColTypeDecimal:
		fieldType = "Numeric"
	case octopus.ColTypeFloat:
		fieldType = "Float"
	case octopus.ColTypeDouble:
		fieldType = "Float"
	case octopus.ColTypeDateTime:
		fieldType = "DateTime"
	case octopus.ColTypeDate:
		fieldType = "Date"
	case octopus.ColTypeTime:
		fieldType = "Time"
	case octopus.ColTypeBlob8:
		fallthrough
	case octopus.ColTypeBlob16:
		fallthrough
	case octopus.ColTypeBlob24:
		fallthrough
	case octopus.ColTypeBlob32:
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

	fieldName, ok := util.ToLowerSnake(column.Name)

	// check python reserved words
	reservedWord := util.IsPythonReservedWord(fieldName)
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
