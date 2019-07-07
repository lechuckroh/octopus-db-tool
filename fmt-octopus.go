package main

import (
	"encoding/json"
	"io/ioutil"
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
	AutoIncremental bool       `json:"ai,omitempty"`
	DefaultValue    string     `json:"default,omitempty"`
	Ref             *Reference `json:"ref,omitempty"`
}

type Table struct {
	Name        string    `json:"name,omitempty"`
	Columns     []*Column `json:"columns,omitempty"`
	Description string    `json:"desc,omitempty"`
	Group       string    `json:"group,omitempty"`
}

type Schema struct {
	Version string   `json:"version,omitempty"`
	Tables  []*Table `json:"tables,omitempty"`
}

func (s *Schema) ToFile(filename string) error {
	data, err := s.ToJson()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (s *Schema) ToJson() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (s *Schema) FromJson(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}
