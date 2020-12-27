package octopus

import (
	"github.com/lechuckroh/octopus-db-tools/util"
	"strings"
)

const (
	ColTypeBinary    = "binary"
	ColTypeBit       = "bit"
	ColTypeBlob16    = "blob16"
	ColTypeBlob24    = "blob24"
	ColTypeBlob32    = "blob32"
	ColTypeBlob8     = "blob8"
	ColTypeBoolean   = "boolean"
	ColTypeChar      = "char"
	ColTypeDate      = "date"
	ColTypeDateTime  = "datetime"
	ColTypeDecimal   = "decimal"
	ColTypeDouble    = "double"
	ColTypeEnum      = "enum"
	ColTypeFloat     = "float"
	ColTypeGeometry  = "geometry"
	ColTypeInt16     = "int16"
	ColTypeInt24     = "int24"
	ColTypeInt32     = "int32"
	ColTypeInt64     = "int64"
	ColTypeInt8      = "int8"
	ColTypeJSON      = "json"
	ColTypePoint     = "point"
	ColTypeSet       = "set"
	ColTypeText16    = "text16"
	ColTypeText24    = "text24"
	ColTypeText32    = "text32"
	ColTypeText8     = "text8"
	ColTypeTime      = "time"
	ColTypeVarbinary = "varbinary"
	ColTypeVarchar   = "varchar"
	ColTypeYear      = "year"
)

var colTypeSet = util.NewStringSet(
	ColTypeBinary,
	ColTypeBit,
	ColTypeBlob16,
	ColTypeBlob24,
	ColTypeBlob32,
	ColTypeBlob8,
	ColTypeBoolean,
	ColTypeChar,
	ColTypeDate,
	ColTypeDateTime,
	ColTypeDecimal,
	ColTypeDouble,
	ColTypeEnum,
	ColTypeFloat,
	ColTypeGeometry,
	ColTypeInt16,
	ColTypeInt24,
	ColTypeInt32,
	ColTypeInt64,
	ColTypeInt8,
	ColTypeJSON,
	ColTypePoint,
	ColTypeSet,
	ColTypeText16,
	ColTypeText24,
	ColTypeText32,
	ColTypeText8,
	ColTypeTime,
	ColTypeVarbinary,
	ColTypeVarchar,
	ColTypeYear,
)

func containsColType(colType string, colTypes []string) bool {
	lowerType := strings.ToLower(colType)
	for _, t := range colTypes {
		if lowerType == t {
			return true
		}
	}
	return false
}

func IsValidColType(colType string) bool {
	return colTypeSet.Contains(colType)
}

// IsColTypeNumeric checks if numeric column type.
func IsColTypeNumeric(colType string) bool {
	return containsColType(colType,
		[]string{
			ColTypeDecimal,
			ColTypeDouble,
			ColTypeFloat,
			ColTypeInt16,
			ColTypeInt24,
			ColTypeInt32,
			ColTypeInt64,
			ColTypeInt8,
		})
}

// IsColTypeAutoIncremental checks if auto incremental column type.
func IsColTypeAutoIncremental(colType string) bool {
	return containsColType(colType,
		[]string{
			ColTypeInt16,
			ColTypeInt24,
			ColTypeInt32,
			ColTypeInt64,
			ColTypeInt8,
		})
}

// IsColTypeClob checks if clob column type.
func IsColTypeClob(colType string) bool {
	return containsColType(colType,
		[]string{
			ColTypeText8,
			ColTypeText16,
			ColTypeText24,
			ColTypeText32,
		})
}

// IsColTypeDecimal checks if decimal column type.
func IsColTypeDecimal(colType string) bool {
	return containsColType(colType,
		[]string{
			ColTypeFloat,
			ColTypeDouble,
			ColTypeDecimal,
		})
}

// IsColTypeClob checks if clob column type.
func IsColTypeString(colType string) bool {
	return containsColType(colType,
		[]string{
			ColTypeText8,
			ColTypeText16,
			ColTypeText24,
			ColTypeText32,
			ColTypeChar,
			ColTypeVarchar,
		})
}
