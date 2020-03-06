package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
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

type Table struct {
	Name        string    `json:"name,omitempty"`
	Columns     []*Column `json:"columns,omitempty"`
	Description string    `json:"desc,omitempty"`
	Group       string    `json:"group,omitempty"`
	ClassName   string    `json:"className,omitempty"`
}

func (t *Table) AddColumn(column *Column) {
	if column != nil {
		t.Columns = append(t.Columns, column)
	}
}

func (t *Table) ColumnByName() map[string]*Column {
	result := make(map[string]*Column)

	for _, column := range t.Columns {
		result[column.Name] = column
	}

	return result
}

func (t *Table) PrimaryKeyNameSet() *StringSet {
	result := NewStringSet()
	for _, column := range t.Columns {
		if column.PrimaryKey {
			result.Add(column.Name)
		}
	}
	return result
}

func (t *Table) UniqueKeyNameSet() *StringSet {
	result := NewStringSet()
	for _, column := range t.Columns {
		if column.UniqueKey {
			result.Add(column.Name)
		}
	}
	return result
}

type TableSlice []*Table

func (s TableSlice) Len() int { return len(s) }
func (s TableSlice) Less(i, j int) bool {
	if s[i].Group < s[j].Group {
		return true
	}
	if s[i].Group == s[j].Group {
		return s[i].Name < s[j].Name
	}
	return false
}

func (s TableSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Schema struct {
	Author  string   `json:"author,omitempty"`
	Name    string   `json:"name,omitempty"`
	Version string   `json:"version,omitempty"`
	Tables  []*Table `json:"tables,omitempty"`
}

func (s *Schema) TableByName() map[string]*Table {
	result := make(map[string]*Table)

	for _, table := range s.Tables {
		result[table.Name] = table
	}

	return result
}

func (s *Schema) Groups() []string {
	groupSet := NewStringSet()
	for _, table := range s.Tables {
		groupSet.Add(table.Group)
	}
	result := groupSet.Slice()
	sort.Strings(result)
	return result
}

func (s *Schema) ToSchema() (*Schema, error) {
	return s, nil
}

func (s *Schema) ToFile(filename string) error {
	data, err := s.ToJson()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (s *Schema) ToJson() ([]byte, error) {
	s.Normalize()
	return json.MarshalIndent(s, "", "  ")
}

func (s *Schema) FromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return s.FromJson(data)
}

func (s *Schema) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

// Normalize converts colume types to lowercase
func (s *Schema) Normalize() {
	sort.Sort(TableSlice(s.Tables))

	for _, table := range s.Tables {
		for _, column := range table.Columns {
			if colType, ok := normalizeColumnType(column); ok {
				column.Type = colType
			} else {
				log.Printf("unknown column type: '%s', table: %s, column: %s",
					column.Type, table.Name, column.Name)
			}
		}
	}
}

// normalizeColumnType converts column type to octopus generalized column type
func normalizeColumnType(col *Column) (string, bool) {
	colType := strings.ToLower(col.Type)

	if colType == "string" || colType == "varchar" || colType == "char" {
		return ColTypeString, true
	}
	if colType == "int" || colType == "integer" || colType == "smallint" {
		return ColTypeInt, true
	}
	if colType == "bigint" || colType == "long" {
		return ColTypeLong, true
	}
	if colType == "datetime" {
		return ColTypeDateTime, true
	}
	if colType == "bool" || colType == "boolean" {
		return ColTypeBoolean, true
	}
	if colType == "number" || colType == "double" || colType == "decimal" {
		return ColTypeDecimal, true
	}
	if colType == "float" {
		return ColTypeFloat, true
	}
	if colType == "text" {
		return ColTypeText, true
	}
	if colType == "date" {
		return ColTypeDate, true
	}
	if colType == "time" {
		return ColTypeTime, true
	}
	if colType == "blob" {
		return ColTypeBlob, true
	}

	return colType, false
}
