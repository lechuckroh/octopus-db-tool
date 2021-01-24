package staruml

import (
	"encoding/json"
	"errors"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io/ioutil"
	"strings"
)

type ImportOption struct {
}

type Importer struct {
	root    *Element
	mapById map[string]*Element
}

func (c *Importer) ImportFile(filename string) (*octopus.Schema, error) {
	c.root = nil

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return c.ImportJson(data)
}

func (c *Importer) ImportJson(data []byte) (*octopus.Schema, error) {
	c.root = nil

	root := &Element{}
	err := json.Unmarshal(data, root)
	if err != nil {
		return nil, err
	}

	c.root = root

	// fill mapById
	c.mapById = make(map[string]*Element)
	c.walkElement(c.root)

	return c.toSchema()
}

func (c *Importer) toSchema() (*octopus.Schema, error) {
	project := c.root

	erdDataModels := c.findByType(project.OwnedElements, "ERDDataModel")
	if len(erdDataModels) == 0 {
		return nil, errors.New("ERDDataModel not found")
	}

	erdEntities := c.findByType(erdDataModels[0].OwnedElements, "ERDEntity")
	var tables []*octopus.Table
	for _, erdEntity := range erdEntities {
		// Create Columns
		var columns []*octopus.Column
		for _, erdColumn := range c.findByType(erdEntity.Columns, "ERDColumn") {
			colType, colSize, colScale := util.ParseType(erdColumn.Type)
			size := uint16(util.ToInt(erdColumn.Length, 0))
			if size == 0 {
				size = colSize
			}

			column := &octopus.Column{
				Name:            erdColumn.Name,
				Type:            colType,
				Size:            size,
				Scale:           colScale,
				NotNull:         !erdColumn.Nullable,
				PrimaryKey:      erdColumn.PrimaryKey,
				UniqueKey:       erdColumn.Unique,
				AutoIncremental: false,
				Description:     strings.TrimSpace(erdColumn.Documentation),
			}

			if erdColumn.ReferenceTo != nil {
				targetColumn := c.findById(erdColumn.ReferenceTo.Ref)
				targetTable := c.findById(targetColumn.Parent.Ref)
				column.Ref = &octopus.Reference{
					Table:  targetTable.Name,
					Column: targetColumn.Name,
				}
			}

			columns = append(columns, column)
		}

		// Create Table
		table := &octopus.Table{
			Name:    erdEntity.Name,
			Columns: columns,
		}
		tables = append(tables, table)
	}

	schema := octopus.Schema{
		Version: project.Name,
		Tables:  tables,
	}
	return &schema, nil
}

func (c *Importer) findByType(elems []*Element, typ string) []*Element {
	var result []*Element

	for _, elem := range elems {
		if elem.ElemType == typ {
			result = append(result, elem)
		}
	}

	return result
}

func (c *Importer) walkElement(elem *Element) {
	c.mapById[elem.ID] = elem

	c.walkElements(elem.OwnedElements)
	c.walkElements(elem.Columns)
}

func (c *Importer) walkElements(elements []*Element) {
	if elements == nil {
		return
	}
	for _, elem := range elements {
		c.walkElement(elem)
	}
}

func (c *Importer) findById(id string) *Element {
	return c.mapById[id]
}
