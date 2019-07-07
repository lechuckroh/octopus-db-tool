package main

import (
	"encoding/json"
	"io/ioutil"
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


// Normalize converts colume types to lowercase
func (s *Schema) Normalize() {
	for _, table := range s.Tables {
		for _, column := range table.Columns {
			column.Type = strings.ToLower(column.Type)
		}
	}
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
