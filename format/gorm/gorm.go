package gorm

import (
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"regexp"
	"strings"
)

type GoStruct struct {
	table        *octopus.Table
	Name         string
	EmbedModel   bool
	Fields       []*GoField
	PKFields     []*GoField
	UniqueFields []*GoField
}

func NewGoStruct(
	table *octopus.Table,
	option *Option,
) *GoStruct {
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

	fields := make([]*GoField, 0)
	pkFields := make([]*GoField, 0)
	uniqueFields := make([]*GoField, 0)
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
	}

	return &GoStruct{
		table:        table,
		Name:         className,
		EmbedModel:   gormModelColumnCount < len(gormModelColumns),
		Fields:       fields,
		PKFields:     pkFields,
		UniqueFields: uniqueFields,
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
	if column.Nullable {
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
		if column.Nullable {
			fieldType = "null.String"
		} else {
			fieldType = "string"
		}
	case octopus.ColTypeBoolean:
		if column.Nullable {
			fieldType = "null.Bool"
		} else {
			fieldType = "bool"
		}
	case octopus.ColTypeInt8:
		if column.Nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int8"
		}
	case octopus.ColTypeInt16:
		if column.Nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int16"
		}
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		if column.Nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int32"
		}
	case octopus.ColTypeInt64:
		if column.Nullable {
			fieldType = "null.Int"
		} else {
			fieldType = "int64"
		}
	case octopus.ColTypeDecimal:
		importSet.Add("github.com/shopspring/decimal")
		if column.Nullable {
			fieldType = "decimal.NullDecimal"
		} else {
			fieldType = "decimal.Decimal"
		}
	case octopus.ColTypeFloat:
		fallthrough
	case octopus.ColTypeDouble:
		if column.Nullable {
			fieldType = "null.Float"
		} else {
			fieldType = "float64"
		}

	case octopus.ColTypeDateTime:
		fallthrough
	case octopus.ColTypeDate:
		fallthrough
	case octopus.ColTypeTime:
		if column.Nullable {
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

	fieldName, ok := util.ToUpperCamel(column.Name)
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
