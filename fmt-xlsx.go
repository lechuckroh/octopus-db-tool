package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
)

const (
	xlsxSheetMeta       = "Meta"
	xlsxMetaVersion     = "version"
	xlsxDefaultFontName = "Verdana"
	xlsxDefaultFontSize = 10
)

type Xlsx struct {
	metaSheet     *xlsx.Sheet
	sheetsByGroup map[string]*xlsx.Sheet
}

func (f *Xlsx) FromFile(filename string) error {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}

	f.sheetsByGroup = make(map[string]*xlsx.Sheet)

	for _, sheet := range xlFile.Sheets {
		sheetName := sheet.Name
		if sheetName == xlsxSheetMeta {
			f.metaSheet = sheet
		} else {
			f.sheetsByGroup[sheet.Name] = sheet
		}
	}
	return nil
}

func (f *Xlsx) ToSchema() (*Schema, error) {
	version := ""
	tables := make([]*Table, 0)

	if f.metaSheet != nil {
		keyValues := f.readMetaSheet()
		version = keyValues[xlsxMetaVersion]
	}

	for groupName, sheet := range f.sheetsByGroup {
		groupTables, err := f.readGroupSheet(groupName, sheet)
		if err != nil {
			return nil, err
		}
		tables = append(tables, groupTables...)
	}

	return &Schema{
		Version: version,
		Tables:  tables,
	}, nil
}

func (f *Xlsx) readMetaSheet() map[string]string {
	result := map[string]string{}

	for _, row := range f.metaSheet.Rows {
		if keyCell := f.getCell(row, 0); keyCell != nil {
			key := keyCell.Value
			valueCell := f.getCell(row, 1)
			if valueCell == nil {
				result[key] = ""
			} else {
				result[key] = valueCell.Value
			}
		}
	}
	return result
}

func (f *Xlsx) readGroupSheet(groupName string, sheet *xlsx.Sheet) ([]*Table, error) {
	tables := make([]*Table, 0)

	var lastTable *Table
	tableFinished := true
	for i, row := range sheet.Rows {
		// skip header row
		if i == 0 {
			continue
		}

		tableName := strings.TrimSpace(f.getCellValue(row, 0))
		columnName := strings.TrimSpace(f.getCellValue(row, 1))

		// finish table if
		// - column is empty
		if columnName == "" && !tableFinished {
			tables = append(tables, lastTable)
			tableFinished = true
			continue
		}

		typeValue := strings.TrimSpace(f.getCellValue(row, 2))
		keyValue := strings.TrimSpace(f.getCellValue(row, 3))
		nullableValue := strings.TrimSpace(f.getCellValue(row, 4))
		attrValue := strings.TrimSpace(f.getCellValue(row, 5))
		description := strings.TrimSpace(f.getCellValue(row, 6))

		// create new table
		if tableFinished {
			if tableName != "" {
				lastTable = &Table{
					Name:        tableName,
					Columns:     make([]*Column, 0),
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

		// add column
		defaultValue := ""
		attrSet := NewStringSet()
		for _, attr := range strings.Split(attrValue, ",") {
			attr = strings.TrimSpace(attr)

			if strings.HasPrefix(attr, "default") {
				tokens := strings.SplitN(attr, ":", 2)
				if len(tokens) == 2 {
					defaultValue = tokens[1]
					continue
				}
			}

			attrSet.Add(strings.ToLower(attr))
		}

		// reference
		var ref *Reference
		if tableName != "" {
			tokens := strings.Split(tableName, ".")
			if len(tokens) == 2 {
				ref = &Reference{
					Table:  tokens[0],
					Column: tokens[1],
				}
			}
		}

		lastTable.AddColumn(&Column{
			Name:            columnName,
			Type:            typeValue,
			Description:     description,
			Size:            0,
			Nullable:        nullableValue != "",
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

func (f *Xlsx) ToFile(schema *Schema, filename string) error {
	file := xlsx.NewFile()
	metaSheet, err := file.AddSheet(xlsxSheetMeta)
	if err != nil {
		return err
	}

	if err = f.fillMetaSheet(metaSheet, schema); err != nil {
		return err
	}

	for _, group := range schema.Groups() {
		groupName := group
		if groupName == "" {
			groupName = "Common"
		}
		sheet, err := file.AddSheet(groupName)
		if err != nil {
			return err
		}

		sheet.SheetViews = []xlsx.SheetView{
			{
				Pane: &xlsx.Pane{
					XSplit:      2,
					YSplit:      1,
					TopLeftCell: "C2",
					ActivePane:  "bottomRight",
					State:       "frozen",
				},
			},
		}

		if err = f.fillGroupSheet(sheet, schema, group); err != nil {
			return err
		}
	}

	return file.Save(filename)
}

func (f *Xlsx) fillMetaSheet(sheet *xlsx.Sheet, schema *Schema) error {
	_ = sheet.SetColWidth(0, 0, 10.5)
	_ = sheet.SetColWidth(1, 1, 10.5)
	font := f.defaultFont()
	style := f.newStyle(nil, nil, nil, font)

	row := sheet.AddRow()
	f.addCells(row, []string{xlsxMetaVersion, schema.Version}, style)
	return nil
}

func (f *Xlsx) newBorder(thickness, color string) *xlsx.Border {
	border := xlsx.NewBorder(thickness, thickness, thickness, thickness)
	if color != "" {
		border.LeftColor = color
		border.RightColor = color
		border.TopColor = color
		border.BottomColor = color
	}
	return border
}

func (f *Xlsx) newSolidFill(color string) *xlsx.Fill {
	return xlsx.NewFill("solid", color, color)
}

func (f *Xlsx) newAlignment(horizontal, vertical string) *xlsx.Alignment {
	return &xlsx.Alignment{
		Horizontal: horizontal,
		Vertical:   vertical,
	}
}

func (f *Xlsx) defaultFont() *xlsx.Font {
	return xlsx.NewFont(xlsxDefaultFontSize, xlsxDefaultFontName)
}

func (f *Xlsx) newStyle(
	fill *xlsx.Fill,
	border *xlsx.Border,
	alignment *xlsx.Alignment,
	font *xlsx.Font,
) *xlsx.Style {
	style := xlsx.NewStyle()
	if fill != nil {
		style.ApplyFill = true
		style.Fill = *fill
	}
	if border != nil {
		style.ApplyBorder = true
		style.Border = *border
	}
	if alignment != nil {
		style.ApplyAlignment = true
		style.Alignment = *alignment
	}
	if font != nil {
		style.ApplyFont = true
		style.Font = *font
	}

	return style
}

func (f *Xlsx) fillGroupSheet(sheet *xlsx.Sheet, schema *Schema, group string) error {
	_ = sheet.SetColWidth(0, 0, 15.5)
	_ = sheet.SetColWidth(1, 1, 13.5)
	_ = sheet.SetColWidth(2, 2, 9.5)
	_ = sheet.SetColWidth(3, 3, 4.0)
	_ = sheet.SetColWidth(4, 4, 4.0)
	_ = sheet.SetColWidth(5, 5, 9.5)
	_ = sheet.SetColWidth(6, 6, 50)


	// alignment
	leftAlignment := f.newAlignment("default", "center")
	centerAlignment := f.newAlignment("center", "center")

	// border
	border := f.newBorder("thin", "")
	lightBorder := f.newBorder("thin", "00B2B2B2")

	// font
	boldFont := f.defaultFont()
	boldFont.Bold = true
	normalFont := f.defaultFont()
	refFont := xlsx.NewFont(8, xlsxDefaultFontName)
	refFont.Italic = true

	headerStyle := f.newStyle(f.newSolidFill("00CCFFCC"), border, centerAlignment, boldFont)
	tableStyle := f.newStyle(f.newSolidFill("00CCFFFF"), border, centerAlignment, boldFont)
	tableDescStyle := f.newStyle(f.newSolidFill("00FFFBCC"), lightBorder, leftAlignment, normalFont)
	normalStyle := f.newStyle(nil, lightBorder, leftAlignment, normalFont)
	boolStyle := f.newStyle(nil, lightBorder, centerAlignment, normalFont)
	referenceStyle := f.newStyle(nil, lightBorder, centerAlignment, refFont)


	// Header
	row := sheet.AddRow()
	f.addCells(row, []string{"Table/Reference", "Column", "Type", "Key", "null", "Attributes", "Description"}, headerStyle)

	tableCount := len(schema.Tables)
	for i, table := range schema.Tables {
		if table.Group != group {
			continue
		}

		// Table
		row = sheet.AddRow()
		f.addCell(row, table.Name, tableStyle)
		f.addCells(row, []string{"", "", "", "", ""}, nil)
		f.addCell(row, strings.TrimSpace(table.Description), tableDescStyle)

		// Columns
		for _, column := range table.Columns {
			row = sheet.AddRow()

			if ref := f.getColumnReference(column); ref != "" {
				f.addCell(row, ref, referenceStyle)
			} else {
				f.addCell(row, "", nil)
			}

			f.addCells(row,
				[]string{
					column.Name,
					f.formatType(column),
				},
				normalStyle)
			f.addCells(row,
				[]string{
					BoolToString(column.PrimaryKey, "P", BoolToString(column.UniqueKey, "U", "")),
					BoolToString(column.Nullable, "O", ""),
				},
				boolStyle)
			f.addCells(row,
				[]string{
					strings.Join(f.getColumnAttributes(column), ", "),
					strings.TrimSpace(column.Description),
				},
				normalStyle)
		}

		// add empty row
		if i < tableCount-1 {
			sheet.AddRow()
		}
	}

	return nil
}

func (f *Xlsx) getCell(row *xlsx.Row, colIdx int) *xlsx.Cell {
	colCount := len(row.Cells)
	if colIdx < colCount {
		return row.Cells[colIdx]
	}
	return nil
}

func (f *Xlsx) getCellValue(row *xlsx.Row, colIdx int) string {
	cell := f.getCell(row, colIdx)
	if cell == nil {
		return ""
	} else {
		return cell.Value
	}
}

func (f *Xlsx) addCell(row *xlsx.Row, value string, style *xlsx.Style) *xlsx.Cell {
	cell := row.AddCell()
	cell.Value = value
	if style != nil {
		cell.SetStyle(style)
	}
	return cell
}

func (f *Xlsx) addCells(row *xlsx.Row, values []string, style *xlsx.Style) {
	for _, value := range values {
		f.addCell(row, value, style)
	}
}

func (f *Xlsx) getColumnAttributes(column *Column) []string {
	attrs := make([]string, 0)

	if column.AutoIncremental {
		attrs = append(attrs, "autoInc")
	}
	if column.DefaultValue != "" {
		attrs = append(attrs, "default:"+column.DefaultValue)
	}

	return attrs
}

func (f *Xlsx) getColumnReference(column *Column) string {
	ref := column.Ref
	if ref == nil {
		return ""
	}
	return fmt.Sprintf("%s.%s", ref.Table, ref.Column)
}

func (f *Xlsx) formatType(column *Column) string {
	if column.Size == 0 {
		return column.Type
	}
	return fmt.Sprintf("%s(%d)", column.Type, column.Size)
}
