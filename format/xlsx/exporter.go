package xlsx

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/tealeg/xlsx"
	"strings"
)

type ExportOption struct {
	UseNotNullColumn bool
}

type Exporter struct {
	schema *octopus.Schema
	option *ExportOption
}

func (c *Exporter) Export(filename string) error {
	file := xlsx.NewFile()
	metaSheet, err := file.AddSheet(xlsxSheetMeta)
	if err != nil {
		return err
	}

	if err = fillMetaSheet(metaSheet, c.schema); err != nil {
		return err
	}

	for _, group := range c.schema.Groups() {
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

		if err = fillGroupSheet(sheet, c.schema, group, c.option.UseNotNullColumn); err != nil {
			return err
		}
	}

	return file.Save(filename)
}

func fillMetaSheet(sheet *xlsx.Sheet, schema *octopus.Schema) error {
	_ = sheet.SetColWidth(0, 0, 10.5)
	_ = sheet.SetColWidth(1, 1, 10.5)
	font := defaultFont()
	style := newStyle(nil, nil, nil, font)

	row := sheet.AddRow()
	addCells(row, []string{xlsxMetaAuthor, schema.Author}, style)
	row = sheet.AddRow()
	addCells(row, []string{xlsxMetaName, schema.Name}, style)
	row = sheet.AddRow()
	addCells(row, []string{xlsxMetaVersion, schema.Version}, style)
	return nil
}

func newBorder(thickness, color string) *xlsx.Border {
	border := xlsx.NewBorder(thickness, thickness, thickness, thickness)
	if color != "" {
		border.LeftColor = color
		border.RightColor = color
		border.TopColor = color
		border.BottomColor = color
	}
	return border
}

func newSolidFill(color string) *xlsx.Fill {
	return xlsx.NewFill("solid", color, color)
}

func newAlignment(horizontal, vertical string) *xlsx.Alignment {
	return &xlsx.Alignment{
		Horizontal: horizontal,
		Vertical:   vertical,
	}
}

func defaultFont() *xlsx.Font {
	return xlsx.NewFont(xlsxDefaultFontSize, xlsxDefaultFontName)
}

func newStyle(
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

func fillGroupSheet(
	sheet *xlsx.Sheet,
	schema *octopus.Schema,
	group string,
	useNotNullColumn bool) error {
	_ = sheet.SetColWidth(0, 0, 18)
	_ = sheet.SetColWidth(1, 1, 13.5)
	_ = sheet.SetColWidth(2, 2, 9.5)
	_ = sheet.SetColWidth(3, 3, 4.0)
	_ = sheet.SetColWidth(4, 4, util.IfThenElseFloat64(useNotNullColumn, 6.0, 4.0))
	_ = sheet.SetColWidth(5, 5, 9.5)
	_ = sheet.SetColWidth(6, 6, 50)

	// alignment
	leftAlignment := newAlignment("default", "center")
	centerAlignment := newAlignment("center", "center")

	// border
	border := newBorder("thin", "")
	lightBorder := newBorder("thin", "00B2B2B2")

	// font
	boldFont := defaultFont()
	boldFont.Bold = true
	normalFont := defaultFont()
	refFont := xlsx.NewFont(8, xlsxDefaultFontName)
	refFont.Italic = true

	headerStyle := newStyle(newSolidFill("00CCFFCC"), border, centerAlignment, boldFont)
	tableStyle := newStyle(newSolidFill("00CCFFFF"), border, centerAlignment, boldFont)
	tableDescStyle := newStyle(newSolidFill("00FFFBCC"), lightBorder, leftAlignment, normalFont)
	normalStyle := newStyle(nil, lightBorder, leftAlignment, normalFont)
	boolStyle := newStyle(nil, lightBorder, centerAlignment, normalFont)
	referenceStyle := newStyle(nil, lightBorder, centerAlignment, refFont)

	// Header
	row := sheet.AddRow()
	nullHeaderText := util.IfThenElseString(useNotNullColumn, headerNotNull, headerNullable)

	addCells(row, []string{
		"Table/Reference",
		"Column",
		"Type",
		"Key",
		nullHeaderText,
		"Attributes",
		"Description",
	}, headerStyle)

	tableCount := len(schema.Tables)
	for i, table := range schema.Tables {
		if table.Group != group {
			continue
		}

		// Table
		row = sheet.AddRow()
		addCell(row, table.Name, tableStyle)
		addCells(row, []string{"", "", "", "", ""}, nil)
		addCell(row, strings.TrimSpace(table.Description), tableDescStyle)

		// Columns
		for _, column := range table.Columns {
			row = sheet.AddRow()

			// table/Reference
			if ref := getColumnReference(column); ref != "" {
				addCell(row, ref, referenceStyle)
			} else {
				addCell(row, "", nil)
			}
			// column
			addCell(row, column.Name, normalStyle)
			// type
			addCell(row, formatType(column), normalStyle)
			// key
			addCell(row, util.BoolToString(column.PrimaryKey, "P", util.BoolToString(column.UniqueKey, "U", "")), boolStyle)
			// nullable
			if useNotNullColumn {
				addCell(row, util.BoolToString(column.NotNull, "O", ""), boolStyle)
			} else {
				addCell(row, util.BoolToString(!column.NotNull, "O", ""), boolStyle)
			}
			// attributes
			addCell(row, strings.Join(getColumnAttributes(column), ", "), normalStyle)
			// description
			addCell(row, strings.TrimSpace(column.Description), normalStyle)
		}

		// add empty row
		if i < tableCount-1 {
			sheet.AddRow()
		}
	}

	return nil
}

func addCell(row *xlsx.Row, value string, style *xlsx.Style) *xlsx.Cell {
	cell := row.AddCell()
	cell.Value = value
	if style != nil {
		cell.SetStyle(style)
	}
	return cell
}

func addCells(row *xlsx.Row, values []string, style *xlsx.Style) {
	for _, value := range values {
		addCell(row, value, style)
	}
}

func getColumnAttributes(column *octopus.Column) []string {
	var attrs []string

	if column.AutoIncremental {
		attrs = append(attrs, "autoInc")
	}
	if column.DefaultValue != "" {
		attrs = append(attrs, "default:"+column.DefaultValue)
	}

	return attrs
}

func getColumnReference(column *octopus.Column) string {
	if ref := column.Ref; ref != nil {
		return fmt.Sprintf("%s.%s", ref.Table, ref.Column)
	}
	return ""
}

func formatType(column *octopus.Column) string {
	if column.Size == 0 {
		return column.Type
	}
	if column.Scale == 0 {
		return fmt.Sprintf("%s(%d)", column.Type, column.Size)
	}
	return fmt.Sprintf("%s(%d,%d)", column.Type, column.Size, column.Scale)
}
