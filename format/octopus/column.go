package octopus

import (
	"errors"
	"fmt"
	"log"
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
		return errors.New("column name is empty")
	}
	if c.AutoIncremental {
		if c.Type != ColTypeInt && c.Type != ColTypeLong {
			if autoCorrect {
				log.Printf("column: '%s', type: '%s' cannnot be autoIncremental. autoIncremental disabled.", c.Name, c.Type)
				c.AutoIncremental = false
			} else {
				return errors.New(fmt.Sprintf("column: '%s', type: '%s' cannnot be autoIncremental", c.Name, c.Type))
			}
		}
	}
	return nil
}
