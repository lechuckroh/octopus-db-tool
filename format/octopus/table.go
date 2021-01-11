package octopus

import (
	"encoding/json"
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io/ioutil"
	"sort"
)

type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type Table struct {
	Name        string    `json:"name,omitempty"`
	Columns     []*Column `json:"columns,omitempty"`
	Description string    `json:"desc,omitempty"`
	Group       string    `json:"group,omitempty"`
	ClassName   string    `json:"className,omitempty"`
	Indices     []*Index  `json:"indices,omitempty"`
}

func (t *Table) AddColumn(column *Column) {
	if column != nil {
		t.Columns = append(t.Columns, column)
	}
}

func (t *Table) ColumnNameMap() map[string]*Column {
	result := make(map[string]*Column)

	for _, column := range t.Columns {
		result[column.Name] = column
	}

	return result
}

func (t *Table) ColumnByName(name string) *Column {
	for _, column := range t.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func (t *Table) PrimaryKeyNameSet() *util.StringSet {
	result := util.NewStringSet()
	for _, column := range t.Columns {
		if column.PrimaryKey {
			result.Add(column.Name)
		}
	}
	return result
}

func (t *Table) UniqueKeyNameSet() *util.StringSet {
	result := util.NewStringSet()
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

// TablesByName returns Table map where key is tableName
func (s *Schema) TablesByName() map[string]*Table {
	result := make(map[string]*Table)

	for _, table := range s.Tables {
		result[table.Name] = table
	}

	return result
}

//TableByName finds Table by name
func (s *Schema) TableByName(name string) *Table {
	for _, table := range s.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func (s *Schema) FindReference(ref Reference) (*Table, *Column) {
	if table := s.TableByName(ref.Table); table != nil {
		if column := table.ColumnByName(ref.Column); column != nil {
			return table, column
		}
	}
	return nil, nil
}

func (s *Schema) Groups() []string {
	groupSet := util.NewStringSet()
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
	if err := s.Normalize(); err != nil {
		return nil, err
	}
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
func (s *Schema) Normalize() error {
	sort.Sort(TableSlice(s.Tables))

	for _, table := range s.Tables {
		for _, column := range table.Columns {
			column.NormalizeType()

			if err := column.Validate(true); err != nil {
				return fmt.Errorf("table: '%s': %w", table.Name, err)
			}
		}
	}
	return nil
}
