package protobuf

import (
	"github.com/iancoleman/strcase"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"log"
	"strings"
)

type PbMessage struct {
	Name    string
	Fields  []*PbField
	Imports []string
}

type PbField struct {
	Type    string
	Name    string
	Number  int
	Rule    ProtoFieldRule
	Comment string
	Default string
	Import  string
}

func NewPbMessage(table *octopus.Table, option *Option) *PbMessage {
	imports := util.NewStringSet()
	fields := make([]*PbField, 0)
	for i, column := range table.Columns {
		field := NewPbField(i+1, column)
		fields = append(fields, field)
		if field.Import != "" {
			imports.Add(field.Import)
		}
	}

	tableName := table.Name
	for _, prefix := range option.RemovePrefixes {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	messageName := strcase.ToCamel(tableName)

	if option.PrefixMapper != nil {
		if prefix := option.PrefixMapper.GetPrefix(table.Group); prefix != "" {
			messageName = prefix + messageName
		}
	}

	return &PbMessage{
		Name:    messageName,
		Fields:  fields,
		Imports: imports.Slice(),
	}
}

func NewPbField(number int, column *octopus.Column) *PbField {
	var fieldType string
	var imp string

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
		fieldType = "string"
	case octopus.ColTypeBoolean:
		fieldType = "bool"
	case octopus.ColTypeInt64:
		fieldType = "int64"
	case octopus.ColTypeInt8:
		fallthrough
	case octopus.ColTypeInt16:
		fallthrough
	case octopus.ColTypeInt24:
		fallthrough
	case octopus.ColTypeInt32:
		fieldType = "int32"
	case octopus.ColTypeDecimal:
		fieldType = "double"
	case octopus.ColTypeFloat:
		fieldType = "float"
	case octopus.ColTypeDouble:
		fieldType = "double"
	case octopus.ColTypeDateTime:
		fieldType = "google.protobuf.Timestamp"
		imp = "google/protobuf/timestamp.proto"
	default:
		if columnType == "bit" {
			if column.Size == 1 {
				fieldType = "bool"
				break
			}
		}
		fieldType = column.Type
		log.Printf("unsupported columnType: %s", columnType)
	}

	fieldName, _ := util.ToLowerCamel(column.Name)

	return &PbField{
		Type:    fieldType,
		Name:    fieldName,
		Number:  number,
		Rule:    Optional,
		Comment: column.Description,
		Default: column.DefaultValue,
		Import:  imp,
	}
}
