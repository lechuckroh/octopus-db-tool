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

	if err = c.fillMetaSheet(metaSheet, c.schema); err != nil {
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

		if err = c.fillGroupSheet(sheet, c.schema, group); err != nil {
			return err
		}
	}

	return file.Save(filename)
}

func (c *Exporter) fillMetaSheet(sheet *xlsx.Sheet, schema *octopus.Schema) error {
	_ = sheet.SetColWidth(0, 0, 10.5)
	_ = sheet.SetColWidth(1, 1, 10.5)
	font := c.defaultFont()
	style := c.newStyle(nil, nil, nil, font)

	row := sheet.AddRow()
	c.addCells(row, []string{xlsxMetaAuthor, schema.Author}, style)
	row = sheet.AddRow()
	c.addCells(row, []string{xlsxMetaName, schema.Name}, style)
	row = sheet.AddRow()
	c.addCells(row, []string{xlsxMetaVersion, schema.Version}, style)
	return nil
}

func (c *Exporter) newBorder(thickness, color string) *xlsx.Border {
	border := xlsx.NewBorder(thickness, thickness, thickness, thickness)
	if color != "" {
		border.LeftColor = color
		border.RightColor = color
		border.TopColor = color
		border.BottomColor = color
	}
	return border
}

func (c *Exporter) newSolidFill(color string) *xlsx.Fill {
	return xlsx.NewFill("solid", color, color)
}

func (c *Exporter) newAlignment(horizontal, vertical string) *xlsx.Alignment {
	return &xlsx.Alignment{
		Horizontal: horizontal,
		Vertical:   vertical,
	}
}

func (c *Exporter) defaultFont() *xlsx.Font {
	return xlsx.NewFont(xlsxDefaultFontSize, xlsxDefaultFontName)
}

func (c *Exporter) newStyle(
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

func (c *Exporter) fillGroupSheet(sheet *xlsx.Sheet, schema *octopus.Schema, group string) error {
	_ = sheet.SetColWidth(0, 0, 18)
	_ = sheet.SetColWidth(1, 1, 13.5)
	_ = sheet.SetColWidth(2, 2, 9.5)
	_ = sheet.SetColWidth(3, 3, 4.0)
	_ = sheet.SetColWidth(4, 4, util.IfThenElseFloat64(c.option.UseNotNullColumn, 6.0, 4.0))
	_ = sheet.SetColWidth(5, 5, 9.5)
	_ = sheet.SetColWidth(6, 6, 50)

	// alignment
	leftAlignment := c.newAlignment("default", "center")
	centerAlignment := c.newAlignment("center", "center")

	// border
	border := c.newBorder("thin", "")
	lightBorder := c.newBorder("thin", "00B2B2B2")

	// font
	boldFont := c.defaultFont()
	boldFont.Bold = true
	normalFont := c.defaultFont()
	refFont := xlsx.NewFont(8, xlsxDefaultFontName)
	refFont.Italic = true

	headerStyle := c.newStyle(c.newSolidFill("00CCFFCC"), border, centerAlignment, boldFont)
	tableStyle := c.newStyle(c.newSolidFill("00CCFFFF"), border, centerAlignment, boldFont)
	tableDescStyle := c.newStyle(c.newSolidFill("00FFFBCC"), lightBorder, leftAlignment, normalFont)
	normalStyle := c.newStyle(nil, lightBorder, leftAlignment, normalFont)
	boolStyle := c.newStyle(nil, lightBorder, centerAlignment, normalFont)
	referenceStyle := c.newStyle(nil, lightBorder, centerAlignment, refFont)

	// Header
	row := sheet.AddRow()
	nullHeaderText := util.IfThenElseString(c.option.UseNotNullColumn, headerNotNull, headerNullable)

	c.addCells(row, []string{
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
		c.addCell(row, table.Name, tableStyle)
		c.addCells(row, []string{"", "", "", "", ""}, nil)
		c.addCell(row, strings.TrimSpace(table.Description), tableDescStyle)

		// Columns
		for _, column := range table.Columns {
			row = sheet.AddRow()

			// table/Reference
			if ref := c.getColumnReference(column); ref != "" {
				c.addCell(row, ref, referenceStyle)
			} else {
				c.addCell(row, "", nil)
			}
			// column
			c.addCell(row, column.Name, normalStyle)
			// type
			c.addCell(row, c.formatType(column), normalStyle)
			// key
			c.addCell(row, util.BoolToString(column.PrimaryKey, "P", util.BoolToString(column.UniqueKey, "U", "")), boolStyle)
			// nullable
			if c.option.UseNotNullColumn {
				c.addCell(row, util.BoolToString(column.NotNull, "O", ""), boolStyle)
			} else {
				c.addCell(row, util.BoolToString(!column.NotNull, "O", ""), boolStyle)
			}
			// attributes
			c.addCell(row, strings.Join(c.getColumnAttributes(column), ", "), normalStyle)
			// description
			c.addCell(row, strings.TrimSpace(column.Description), normalStyle)
		}

		// add empty row
		if i < tableCount-1 {
			sheet.AddRow()
		}
	}

	return nil
}

func (c *Exporter) addCell(row *xlsx.Row, value string, style *xlsx.Style) *xlsx.Cell {
	cell := row.AddCell()
	cell.Value = value
	if style != nil {
		cell.SetStyle(style)
	}
	return cell
}

func (c *Exporter) addCells(row *xlsx.Row, values []string, style *xlsx.Style) {
	for _, value := range values {
		c.addCell(row, value, style)
	}
}

func (c *Exporter) getColumnAttributes(column *octopus.Column) []string {
	attrs := make([]string, 0)

	if column.AutoIncremental {
		attrs = append(attrs, "autoInc")
	}
	if column.DefaultValue != "" {
		attrs = append(attrs, "default:"+column.DefaultValue)
	}

	return attrs
}

func (c *Exporter) getColumnReference(column *octopus.Column) string {
	ref := column.Ref
	if ref == nil {
		return ""
	}
	return fmt.Sprintf("%s.%s", ref.Table, ref.Column)
}

func (c *Exporter) formatType(column *octopus.Column) string {
	if column.Size == 0 {
		return column.Type
	}
	if column.Scale == 0 {
		return fmt.Sprintf("%s(%d)", column.Type, column.Size)
	}
	return fmt.Sprintf("%s(%d,%d)", column.Type, column.Size, column.Scale)
}
