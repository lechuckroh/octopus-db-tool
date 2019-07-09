package main

import (
	"encoding/json"
	"io/ioutil"
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
	Nullable        bool       `json:"nullable,omitempty"`
	PrimaryKey      bool       `json:"pk,omitempty"`
	UniqueKey       bool       `json:"unique,omitempty"`
	AutoIncremental bool       `json:"autoinc,omitempty"`
	DefaultValue    string     `json:"default,omitempty"`
	Ref             *Reference `json:"ref,omitempty"`
}

type Table struct {
	Name        string    `json:"name,omitempty"`
	Columns     []*Column `json:"columns,omitempty"`
	Description string    `json:"desc,omitempty"`
	Group       string    `json:"group,omitempty"`
	ClassName   string    `json:"className,omitempty"`
}

type Schema struct {
	Version string   `json:"version,omitempty"`
	Tables  []*Table `json:"tables,omitempty"`
}

func (t *Table) AddColumn(column *Column) {
	if column != nil {
		t.Columns = append(t.Columns, column)
	}
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

// Normalize converts colume types to lowercase
func (s *Schema) Normalize() {
	sort.Sort(TableSlice(s.Tables))

	for _, table := range s.Tables {
		for _, column := range table.Columns {
			column.Type = strings.ToLower(column.Type)
		}
	}
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
