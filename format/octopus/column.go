package octopus

import (
	"log"
	"strings"
)

type Reference struct {
	Table  string `json:"table,omitempty"`
	Column string `json:"column,omitempty"`
}

type Column struct {
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Description     string     `json:"desc,omitempty"`
	Size            uint16     `json:"size,omitempty"`
	Scale           uint16     `json:"scale,omitempty"`
	Nullable        bool       `json:"nullable,omitempty"`
	PrimaryKey      bool       `json:"pk,omitempty"`
	UniqueKey       bool       `json:"unique,omitempty"`
	AutoIncremental bool       `json:"autoinc,omitempty"`
	DefaultValue    string     `json:"default,omitempty"`
	Ref             *Reference `json:"ref,omitempty"`
}

func (c *Column) NormalizeType() {
	c.Type = c.toNormalizedType(strings.ToLower(c.Type))
}

func (c *Column) toNormalizedType(colType string) string {
	switch colType {
	case "string":
		return ColTypeVarchar
	case "tinyint":
		return ColTypeInt8
	case "smallint":
		return ColTypeInt16
	case "mediumint":
		return ColTypeInt24
	case "int":
		fallthrough
	case "integer":
		return ColTypeInt32
	case "bigint":
		fallthrough
	case "long":
		return ColTypeInt64
	case "numeric":
		return ColTypeDecimal
	case "real":
		return ColTypeDouble
	case "timestamp":
		return ColTypeDateTime
	case "tinyblob":
		return ColTypeBlob8
	case "blob":
		return ColTypeBlob16
	case "mediumblob":
		return ColTypeBlob24
	case "longblob":
		return ColTypeBlob32
	case "tinytext":
		return ColTypeText8
	case "text":
		return ColTypeText16
	case "mediumtext":
		return ColTypeText24
	case "longtext":
		return ColTypeText32
	}
	return colType
}

func (c *Column) IsRenamed(target *Column, excludeDescription bool) bool {
	return c.Type == target.Type &&
		(excludeDescription || (c.Description == target.Description)) &&
		c.Size == target.Size &&
		c.Scale == target.Scale &&
		c.Nullable == target.Nullable &&
		c.PrimaryKey == target.PrimaryKey &&
		c.UniqueKey == target.UniqueKey &&
		c.AutoIncremental == target.AutoIncremental &&
		c.DefaultValue == target.DefaultValue
}

func (c *Column) Validate(autoCorrect bool) error {
	if c.Name == "" {
		return &EmptyColumnNameError{Column: c}
	}

	if c.AutoIncremental && !IsColTypeAutoIncremental(c.Type) {
		if autoCorrect {
			log.Printf("column: '%s', type: '%s' cannnot be autoIncremental. autoIncremental disabled.", c.Name, c.Type)
			c.AutoIncremental = false
		} else {
			return &InvalidAutoIncrementalError{Column: c}
		}
	}

	if !IsValidColType(c.Type) {
		return &InvalidColumnTypeError{Column: c}
	}

	return nil
}
