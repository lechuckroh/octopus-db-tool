package gorm

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"regexp"
	"strings"
)

type GoAssocationField struct {
	Name       string
	Type       string
	Array      bool
	ForeignKey string
	Reference  string
}

type GoStruct struct {
	table             *octopus.Table
	Name              string
	EmbedModel        bool
	Fields            []*GoField
	PKFields          []*GoField
	UniqueFields      []*GoField
	AssociationFields []*GoAssocationField
}

type IProcessor interface {
	StructName(table *octopus.Table) string
	Reference(ref octopus.Reference) (*octopus.Table, *octopus.Column)
}

func NewGoStruct(
	table *octopus.Table,
	p IProcessor,
) *GoStruct {
	var fields []*GoField
	var pkFields []*GoField
	var uniqueFields []*GoField
	var associationFields []*GoAssocationField
	gormModelColumnCount := 0
	for _, column := range table.Columns {
		field := NewGoField(column)
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

		// reference
		if ref := column.Ref; ref != nil {
			refTable, refColumn := p.Reference(*ref)
			if refTable != nil && refColumn != nil {
				refType := p.StructName(refTable)
				associationFields = append(associationFields, &GoAssocationField{
					Name:       refType,
					Type:       refType,
					Array:      ref.Relationship == octopus.RefOneToMany,
					ForeignKey: field.Name,
					Reference:  NewGoField(refColumn).Name,
				})
			}
		}
	}

	return &GoStruct{
		table:             table,
		Name:              p.StructName(table),
		EmbedModel:        gormModelColumnCount == len(gormModelColumns),
		Fields:            fields,
		PKFields:          pkFields,
		UniqueFields:      uniqueFields,
		AssociationFields: associationFields,
	}
}

type GoField struct {
	Column       *octopus.Column
	Name         string
	Type         string
	OverrideName bool
	Imports      []string
}

func NewGoField(column *octopus.Column) *GoField {
	var fieldType string
	importSet := util.NewStringSet()
	nullable := !column.NotNull
	if nullable {
		importSet.Add("gopkg.in/guregu/null.v4")
	}

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
		if nullable {
			fieldType = "null.String"
		} else {
			fieldType = "string"
		}
	case octopus.ColTypeBoolean:
		if nullable {
			fieldType = "null.Bool"
		} else {
			fieldType = "bool"
		}
	case octopus.ColTypeInt8:
		if nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int8"
		}
	case octopus.ColTypeInt16:
		if nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int16"
		}
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		if nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int32"
		}
	case octopus.ColTypeInt64:
		if nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int64"
		}
	case octopus.ColTypeDecimal:
		importSet.Add("github.com/shopspring/decimal")
		if nullable {
			fieldType = "decimal.NullDecimal"
		} else {
			fieldType = "decimal.Decimal"
		}
	case octopus.ColTypeFloat:
		if nullable {
			fieldType = "null.Float"
		} else {
			fieldType = "float32"
		}
	case octopus.ColTypeDouble:
		if nullable {
			fieldType = "null.Float"
		} else {
			fieldType = "float64"
		}
	case octopus.ColTypeDateTime:
		fallthrough
	case octopus.ColTypeDate:
		fallthrough
	case octopus.ColTypeTime:
		if nullable {
			fieldType = "null.Time"
		} else {
			importSet.Add("time")
			fieldType = "time.Time"
		}
	case octopus.ColTypeBlob8:
		fallthrough
	case octopus.ColTypeBlob16:
		fallthrough
	case octopus.ColTypeBlob24:
		fallthrough
	case octopus.ColTypeBlob32:
		fieldType = "[]byte"
	case octopus.ColTypeBit:
		if nullable {
			fieldType = "*byte"
		} else {
			fieldType = "byte"
		}
	default:
		fieldType = "interface{}"
	}

	fieldName, ok := util.ToUpperCamel(column.Name)
	overrideName := !ok

	// replace fieldname: Id -> ID
	if strings.HasSuffix(fieldName, "Id") {
		re := regexp.MustCompile(`Id$`)
		fieldName = string(re.ReplaceAll([]byte(fieldName), []byte("ID")))
	}
	// number prefix
	if matched, _ := regexp.MatchString(`^\d+.*$`, fieldName); matched {
		fieldName = "_" + fieldName
		overrideName = true
	}

	return &GoField{
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
