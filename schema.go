package main

import "encoding/json"

type Reference struct {
	Table    string `json:"table,omitempty"`
	Column   string `json:"column,omitempty"`
	Nullable bool   `json:"nullable,omitempty"`
	Multi    bool   `json:"multi,omitempty"`
}

type Column struct {
	Name            string     `json:"name,omitempty"`
	Type            string     `json:"type,omitempty"`
	Description     string     `json:"desc,omitempty"`
	Size            int8       `json:"size,omitempty"`
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
	Tags        []string  `json:"tags,omitempty"`
}

type Schema struct {
	Version string   `json:"version,omitempty"`
	Tables  []*Table `json:"tables,omitempty"`
}

func (s *Schema) toJson() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (s *Schema) fromJson(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}
