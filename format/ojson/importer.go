package ojson

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"io"
	"io/ioutil"
)

type ImportOption struct {
}

type Importer struct {
	option *ImportOption
}

func (c *Importer) Import(reader io.Reader) (*octopus.Schema, error) {
	if bytes, err := ioutil.ReadAll(reader); err != nil {
		return nil, err
	} else {
		return c.ImportJSON(bytes)
	}
}

func (c *Importer) ImportFile(filename string) (*octopus.Schema, error) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return c.ImportJSON(data)
	}
}

func (c *Importer) ImportJSON(data []byte) (*octopus.Schema, error) {
	schema1 := Schema{}
	if err := schema1.FromJson(data); err != nil {
		return nil, err
	}
	schema1.Normalize()

	var tables []*octopus.Table
	for _, table1 := range schema1.Tables {
		tables = append(tables, c.fromOjsonTable(table1))
	}

	schema2 := octopus.Schema{
		Author:  schema1.Author,
		Name:    schema1.Name,
		Version: schema1.Version,
		Tables:  tables,
	}
	return &schema2, nil
}

func (c *Importer) fromOjsonTable(table1 *Table) *octopus.Table {
	var columns []*octopus.Column
	for _, column1 := range table1.Columns {
		columns = append(columns, c.fromOjsonColumn(column1))
	}

	return &octopus.Table{
		Name:        table1.Name,
		Columns:     columns,
		Description: table1.Description,
		Group:       table1.Group,
	}
}

func (c *Importer) fromOjsonColumn(column1 *Column) *octopus.Column {
	var ref *octopus.Reference
	if column1.Ref != nil {
		ref = &octopus.Reference{
			Table:  column1.Ref.Table,
			Column: column1.Ref.Column,
		}
	}
	return &octopus.Column{
		Name:            column1.Name,
		Type:            column1.Type,
		Description:     column1.Description,
		Size:            column1.Size,
		Scale:           column1.Scale,
		NotNull:         column1.NotNull,
		PrimaryKey:      column1.PrimaryKey,
		UniqueKey:       column1.UniqueKey,
		AutoIncremental: column1.AutoIncremental,
		DefaultValue:    column1.Description,
		Ref:             ref,
	}
}
