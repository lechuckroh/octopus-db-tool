package xlsx

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/tealeg/xlsx"
	"strings"
)

type ImportOption struct {
}

type Importer struct {
	metaSheet        *xlsx.Sheet
	sheetsByGroup    map[string]*xlsx.Sheet
	useNotNullColumn bool
}

func (c *Importer) Import(filename string) (*octopus.Schema, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	c.sheetsByGroup = make(map[string]*xlsx.Sheet)

	for _, sheet := range xlFile.Sheets {
		sheetName := sheet.Name
		if sheetName == xlsxSheetMeta {
			c.metaSheet = sheet
		} else {
			c.sheetsByGroup[sheet.Name] = sheet
		}
	}
	return c.toSchema()
}

func (c *Importer) toSchema() (*octopus.Schema, error) {
	author := ""
	name := ""
	version := ""
	var tables []*octopus.Table

	if c.metaSheet != nil {
		keyValues := c.readMetaSheet()
		author = keyValues[xlsxMetaAuthor]
		name = keyValues[xlsxMetaName]
		version = keyValues[xlsxMetaVersion]
	}

	for groupName, sheet := range c.sheetsByGroup {
		groupTables, err := readGroupSheet(groupName, sheet)
		if err != nil {
			return nil, err
		}
		tables = append(tables, groupTables...)
	}

	return &octopus.Schema{
		Author:  author,
		Name:    name,
		Version: version,
		Tables:  tables,
	}, nil
}

func (c *Importer) readMetaSheet() map[string]string {
	result := map[string]string{}

	for _, row := range c.metaSheet.Rows {
		if keyCell := getCell(row, 0); keyCell != nil {
			key := keyCell.Value
			valueCell := getCell(row, 1)
			if valueCell == nil {
				result[key] = ""
			} else {
				result[key] = valueCell.Value
			}
		}
	}
	return result
}

func readGroupSheet(groupName string, sheet *xlsx.Sheet) ([]*octopus.Table, error) {
	var tables []*octopus.Table

	var lastTable *octopus.Table
	useNotNullColumn := false
	for i, row := range sheet.Rows {
		// skip header row
		if i == 0 {
			if strings.TrimSpace(getCellValue(row, 4)) == headerNotNull {
				useNotNullColumn = true
			}
			continue
		}

		tableName := strings.TrimSpace(getCellValue(row, 0))
		columnName := strings.TrimSpace(getCellValue(row, 1))
		typeValue := strings.TrimSpace(getCellValue(row, 2))
		keyValue := strings.TrimSpace(getCellValue(row, 3))
		nullableValue := strings.TrimSpace(getCellValue(row, 4))
		attrValue := strings.TrimSpace(getCellValue(row, 5))
		description := strings.TrimSpace(getCellValue(row, 6))

		attrMap := parseAttributes(attrValue)

		// table name row
		if typeValue == typeTable {
			// finalize lastTable
			if lastTable != nil {
				tables = append(tables, lastTable)
			}

			if tableName == "" {
				return tables, fmt.Errorf("row[%d]: table name is empty", i)
			}

			// create new table
			lastTable = &octopus.Table{
				Name:        tableName,
				Columns:     make([]*octopus.Column, 0),
				Description: description,
				Group:       groupName,
				ClassName:   attrMap[attrClass],
			}
			continue
		}

		// skip if table is not started yet
		if lastTable == nil {
			continue
		}

		// skip if column name is empty
		if columnName == "" {
			continue
		}

		// column type
		colType, colSize, colScale := util.ParseType(typeValue)

		// index
		if keyValue == keyIndex {
			indexName := tableName
			if index := lastTable.IndexByName(indexName); index != nil {
				index.AddColumn(columnName)
			} else {
				lastTable.AddIndex(indexName, columnName)
			}
		} else {
			// reference
			ref := parseReference(tableName)

			lastTable.AddColumn(&octopus.Column{
				Name:            columnName,
				Type:            colType,
				Description:     description,
				Size:            colSize,
				Scale:           colScale,
				NotNull:         util.IfThenElseBool(useNotNullColumn, nullableValue != "", nullableValue == ""),
				PrimaryKey:      keyValue == keyPrimary,
				UniqueKey:       keyValue == keyUnique,
				AutoIncremental: attrMap[attrAutoInc] == "true",
				DefaultValue:    attrMap[attrDefault],
				OnUpdate:        attrMap[attrOnUpdate],
				Ref:             ref,
			})
		}
	}

	if lastTable != nil {
		tables = append(tables, lastTable)
	}

	return tables, nil
}

func parseAttributes(value string) map[string]string {
	result := make(map[string]string)
	for _, attr := range strings.Split(value, ",") {
		attr = strings.TrimSpace(attr)

		tokens := strings.SplitN(attr, "=", 2)
		key := tokens[0]

		if len(tokens) == 2 {
			value := tokens[1]
			result[key] = value
		} else {
			result[key] = "true"
		}
	}
	return result
}

func parseReference(s string) *octopus.Reference {
	if s != "" {
		tokens := strings.Split(s, ".")
		if len(tokens) == 2 {
			table := tokens[0]
			column := tokens[1]

			if table != "" && column != "" {
				var relationship string
				switch []rune(table)[0] {
				case '>':
					relationship = octopus.RefManyToOne
				case '<':
					relationship = octopus.RefOneToMany
				case '-':
					relationship = octopus.RefOneToOne
				default:
					relationship = octopus.RefManyToOne
				}
				return &octopus.Reference{
					Table:        table,
					Column:       column,
					Relationship: relationship,
				}
			}
		}
	}
	return nil
}

func fixColumnValue(colType string, value string) string {
	if util.IsBooleanType(colType) {
		return util.IfThenElseString(value == "true" || value == "1", "true", "false")
	}
	return value
}

func getCell(row *xlsx.Row, colIdx int) *xlsx.Cell {
	colCount := len(row.Cells)
	if colIdx < colCount {
		return row.Cells[colIdx]
	}
	return nil
}

func getCellValue(row *xlsx.Row, colIdx int) string {
	cell := getCell(row, colIdx)
	if cell == nil {
		return ""
	} else {
		return cell.Value
	}
}
