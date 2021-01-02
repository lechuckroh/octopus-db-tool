package octopus

import (
	"fmt"
	"log"
	"strings"
)

type Reference struct {
	Table  string `json:"table,omitempty"`
	Column string `json:"column,omitempty"`
}

const FnPrefix = "fn::"
const FnPrefixLen = len(FnPrefix)

type Column struct {
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Description     string     `json:"description,omitempty"`
	Size            uint16     `json:"size,omitempty"`
	Scale           uint16     `json:"scale,omitempty"`
	NotNull         bool       `json:"notnull,omitempty"`
	PrimaryKey      bool       `json:"pk,omitempty"`
	UniqueKey       bool       `json:"unique,omitempty"`
	AutoIncremental bool       `json:"autoinc,omitempty"`
	DefaultValue    string     `json:"default,omitempty"`
	OnUpdate        string     `json:"onupdate,omitempty"`
	Values          []string   `json:"values,omitempty"`
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
		c.NotNull == target.NotNull &&
		c.PrimaryKey == target.PrimaryKey &&
		c.UniqueKey == target.UniqueKey &&
		c.AutoIncremental == target.AutoIncremental &&
		c.DefaultValue == target.DefaultValue &&
		c.OnUpdate == target.OnUpdate
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

	memberCount := len(c.Values)
	if c.Type == ColTypeSet {
		if memberCount <= 0 || memberCount > 64 {
			return &InvalidColumnValuesError{
				Column: c,
				Msg:    fmt.Sprintf("invalud set member count: %d (0 < count < 64)", memberCount),
			}
		}
	}
	if c.Type == ColTypeEnum {
		if memberCount <= 0 {
			return &InvalidColumnValuesError{Column: c, Msg: "empty enum values"}
		}
	}

	if !IsValidColType(c.Type) {
		return &InvalidColumnTypeError{Column: c}
	}

	return nil
}

func (c *Column) SetDefaultValue(value interface{}) {
	if value == nil {
		c.DefaultValue = ""
	} else {
		c.DefaultValue = fmt.Sprintf("%v", value)
	}
}

func (c *Column) SetDefaultValueFn(fnName string) {
	c.DefaultValue = FnPrefix + fnName
}

// GetDefaultValue returns defaultValue.
// bool is true if defaultValue is function call.
func (c *Column) GetDefaultValue() (string, bool) {
	return c.getValue(c.DefaultValue)
}

func (c *Column) SetOnUpdate(value interface{}) {
	if value == nil {
		c.OnUpdate = ""
	} else {
		c.OnUpdate = fmt.Sprintf("%v", value)
	}
}

func (c *Column) SetOnUpdateFn(fnName string) {
	c.OnUpdate = FnPrefix + fnName
}

// GetOnUpdate returns onUpdate value.
// bool is true if onUpdate is function call.
func (c *Column) GetOnUpdate() (string, bool) {
	return c.getValue(c.OnUpdate)
}

func (c *Column) getValue(value string) (string, bool) {
	if strings.HasPrefix(value, FnPrefix) {
		return value[FnPrefixLen:], true
	}
	return value, false
}
