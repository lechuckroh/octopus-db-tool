package xlsx

import (
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
	tables := make([]*octopus.Table, 0)

	if c.metaSheet != nil {
		keyValues := c.readMetaSheet()
		author = keyValues[xlsxMetaAuthor]
		name = keyValues[xlsxMetaName]
		version = keyValues[xlsxMetaVersion]
	}

	for groupName, sheet := range c.sheetsByGroup {
		groupTables, err := c.readGroupSheet(groupName, sheet)
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
		if keyCell := c.getCell(row, 0); keyCell != nil {
			key := keyCell.Value
			valueCell := c.getCell(row, 1)
			if valueCell == nil {
				result[key] = ""
			} else {
				result[key] = valueCell.Value
			}
		}
	}
	return result
}

func (c *Importer) readGroupSheet(groupName string, sheet *xlsx.Sheet) ([]*octopus.Table, error) {
	tables := make([]*octopus.Table, 0)

	var lastTable *octopus.Table
	tableFinished := true
	useNotNullColumn := false
	for i, row := range sheet.Rows {
		// skip header row
		if i == 0 {
			if strings.TrimSpace(c.getCellValue(row, 4)) == headerNotNull {
				useNotNullColumn = true
			}
			continue
		}

		tableName := strings.TrimSpace(c.getCellValue(row, 0))
		columnName := strings.TrimSpace(c.getCellValue(row, 1))

		// finish table if
		// - column is empty
		if columnName == "" && !tableFinished {
			tables = append(tables, lastTable)
			tableFinished = true
			continue
		}

		typeValue := strings.TrimSpace(c.getCellValue(row, 2))
		keyValue := strings.TrimSpace(c.getCellValue(row, 3))
		nullableValue := strings.TrimSpace(c.getCellValue(row, 4))
		attrValue := strings.TrimSpace(c.getCellValue(row, 5))
		description := strings.TrimSpace(c.getCellValue(row, 6))

		// create new table
		if tableFinished {
			if tableName != "" {
				lastTable = &octopus.Table{
					Name:        tableName,
					Columns:     make([]*octopus.Column, 0),
					Description: description,
					Group:       groupName,
					ClassName:   typeValue,
				}
				tableFinished = false
			}
			continue
		}
		if lastTable == nil {
			continue
		}

		// column type
		colType, colSize, colScale := util.ParseType(typeValue)

		// add column
		defaultValue := ""
		attrSet := util.NewStringSet()
		for _, attr := range strings.Split(attrValue, ",") {
			attr = strings.TrimSpace(attr)

			if strings.HasPrefix(attr, "default") {
				tokens := strings.SplitN(attr, ":", 2)
				if len(tokens) == 2 {
					defaultValue = c.fixDefaultValue(colType, tokens[1])
					continue
				}
			}

			attrSet.Add(strings.ToLower(attr))
		}

		// reference
		var ref *octopus.Reference
		if tableName != "" {
			tokens := strings.Split(tableName, ".")
			if len(tokens) == 2 {
				ref = &octopus.Reference{
					Table:  tokens[0],
					Column: tokens[1],
				}
			}
		}

		lastTable.AddColumn(&octopus.Column{
			Name:            columnName,
			Type:            colType,
			Description:     description,
			Size:            colSize,
			Scale:           colScale,
			NotNull:         util.TernaryBool(useNotNullColumn, nullableValue != "", nullableValue == ""),
			PrimaryKey:      keyValue == "P",
			UniqueKey:       keyValue == "U",
			AutoIncremental: attrSet.ContainsAny([]string{"ai", "autoinc", "auto_inc", "auto_incremental"}),
			DefaultValue:    defaultValue,
			Ref:             ref,
		})
	}

	if !tableFinished && lastTable != nil {
		tables = append(tables, lastTable)
	}

	return tables, nil
}

func (c *Importer) fixDefaultValue(colType string, defaultValue string) string {
	if util.IsBooleanType(colType) {
		return util.TernaryString(defaultValue == "true" || defaultValue == "1", "true", "false")
	}
	return defaultValue
}

func (c *Importer) getCell(row *xlsx.Row, colIdx int) *xlsx.Cell {
	colCount := len(row.Cells)
	if colIdx < colCount {
		return row.Cells[colIdx]
	}
	return nil
}

func (c *Importer) getCellValue(row *xlsx.Row, colIdx int) string {
	cell := c.getCell(row, colIdx)
	if cell == nil {
		return ""
	} else {
		return cell.Value
	}
}
