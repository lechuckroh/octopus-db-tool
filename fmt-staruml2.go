package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type StarUML2 struct {
	root    *StarUML2Element
	mapById map[string]*StarUML2Element
}

type StarUML2Ref struct {
	Ref string `json:"$ref"`
}

type StarUML2Element struct {
	ElemType      string       `json:"_type"`
	ID            string       `json:"_id"`
	Parent        *StarUML2Ref `json:"_parent"`
	Name          string       `json:"name"`
	Length        interface{}  `json:"length"`
	Type          string       `json:"type"`
	PrimaryKey    bool         `json:"primaryKey"`
	Unique        bool         `json:"unique"`
	Nullable      bool         `json:"nullable"`
	Documentation string       `json:"documentation"`
	Head          *StarUML2Ref `json:"head"`
	End1          *StarUML2Ref `json:"end1"`
	End2          *StarUML2Ref `json:"end2"`
	Reference     *StarUML2Ref `json:"reference"`
	ReferenceTo   *StarUML2Ref `json:"referenceTo"`

	OwnedElements []*StarUML2Element `json:"ownedElements"`
	Columns       []*StarUML2Element `json:"columns"`
}

func (f *StarUML2) FromFile(filename string) error {
	f.root = nil

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return f.FromJson(data)
}

func (f *StarUML2) FromJson(data []byte) error {
	f.root = nil

	root := &StarUML2Element{}
	err := json.Unmarshal(data, root)
	if err != nil {
		return err
	}

	f.root = root

	// fill mapById
	f.mapById = make(map[string]*StarUML2Element)
	f.walkElement(f.root)

	return nil
}

func (f *StarUML2) ToSchema() (*Schema, error) {
	project := f.root

	erdDataModels := f.findByType(project.OwnedElements, "ERDDataModel")
	if len(erdDataModels) == 0 {
		return nil, errors.New("ERDDataModel not found")
	}

	erdEntities := f.findByType(erdDataModels[0].OwnedElements, "ERDEntity")
	tables := make([]*Table, 0)
	for _, erdEntity := range erdEntities {
		// Create Columns
		columns := make([]*Column, 0)
		for _, erdColumn := range f.findByType(erdEntity.Columns, "ERDColumn") {
			column := &Column{
				Name:        erdColumn.Name,
				Type:        erdColumn.Type,
				Size:        uint16(toInt(erdColumn.Length, 0)),
				Nullable:    erdColumn.Nullable,
				PrimaryKey:  erdColumn.PrimaryKey,
				UniqueKey:   erdColumn.Unique,
				Description: erdColumn.Documentation,
			}

			if erdColumn.ReferenceTo != nil {
				targetColumn := f.findById(erdColumn.ReferenceTo.Ref)
				targetTable := f.findById(targetColumn.Parent.Ref)
				column.Ref = &Reference{
					Table:  targetTable.Name,
					Column: targetColumn.Name,
				}
			}

			columns = append(columns, column)
		}

		// Create Table
		table := &Table{
			Name:    erdEntity.Name,
			Columns: columns,
		}
		tables = append(tables, table)
	}

	schema := Schema{
		Version: project.Name,
		Tables:  tables,
	}
	return &schema, nil
}

func (f *StarUML2) findByType(elems []*StarUML2Element, typ string) []*StarUML2Element {
	result := make([]*StarUML2Element, 0)

	for _, elem := range elems {
		if elem.ElemType == typ {
			result = append(result, elem)
		}
	}

	return result
}

func (f *StarUML2) walkElement(elem *StarUML2Element) {
	f.mapById[elem.ID] = elem

	f.walkElements(elem.OwnedElements)
	f.walkElements(elem.Columns)
}

func (f *StarUML2) walkElements(elements []*StarUML2Element) {
	if elements == nil {
		return
	}
	for _, elem := range elements {
		f.walkElement(elem)
	}
}

func (f *StarUML2) findById(id string) *StarUML2Element {
	return f.mapById[id]
}
