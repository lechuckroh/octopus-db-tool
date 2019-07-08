package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
)

type Xlsx struct {
}

func (f *Xlsx) FromFile(filename string) error {
	return errors.New("not implemented")
}

func (f *Xlsx) ToSchema() (*Schema, error) {
	// TODO
	return nil, errors.New("not implemented")
}

func (f *Xlsx) ToFile(schema *Schema, filename string) error {
	file := xlsx.NewFile()
	metaSheet, err := file.AddSheet("Meta")
	if err != nil {
		return err
	}

	if err = f.fillMetaSheet(metaSheet, schema); err != nil {
		return err
	}

	for _, group := range schema.Groups() {
		groupName := group
		if groupName == "" {
			groupName = "NoName"
		}
		sheet, err := file.AddSheet(groupName)
		if err != nil {
			return err
		}

		if err = f.fillGroupSheet(sheet, schema, group); err != nil {
			return err
		}
	}

	return file.Save(filename)
}

func (f *Xlsx) fillMetaSheet(sheet *xlsx.Sheet, schema *Schema) error {
	row := sheet.AddRow()
	f.addCell(row, "version", nil)
	f.addCell(row, schema.Version, nil)
	return nil
}

func (f *Xlsx) fillGroupSheet(sheet *xlsx.Sheet, schema *Schema, group string) error {
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
	lightBorder := *xlsx.NewBorder("thin", "thin", "thin", "thin")
	lightBorderColor := "00B2B2B2"
	lightBorder.LeftColor = lightBorderColor
	lightBorder.RightColor = lightBorderColor
	lightBorder.TopColor = lightBorderColor
	lightBorder.BottomColor = lightBorderColor

	headerStyle := xlsx.NewStyle()
	headerStyle.ApplyFill = true
	headerStyle.Fill = *xlsx.NewFill("solid", "00CCFFCC", "00CCFFCC")
	headerStyle.ApplyBorder = true
	headerStyle.Border = border
	headerStyle.Font = *xlsx.NewFont(14, "Verdana")
	headerStyle.Font.Bold = true

	tableStyle := xlsx.NewStyle()
	tableStyle.ApplyFill = true
	tableStyle.Fill = *xlsx.NewFill("solid", "00CCFFFF", "00CCFFFF")
	tableStyle.ApplyBorder = true
	tableStyle.Border = border
	tableStyle.Font.Bold = true

	normalStyle := xlsx.NewStyle()
	normalStyle.ApplyBorder = true
	normalStyle.Border = lightBorder


	// Header
	row := sheet.AddRow()
	f.addCells(row, []string{"Table/Reference", "Column", "Type", "Attributes", "Description"}, headerStyle)

	for _, table := range schema.Tables {
		if table.Group != group {
			continue
		}

		// Table
		row = sheet.AddRow()
		f.addCell(row, table.Name, tableStyle)

		// Columns
		for _, column := range table.Columns {
			row = sheet.AddRow()
			f.addCell(row, f.getColumnReference(column), nil)
			f.addCells(row,
				[]string{
					column.Name,
					f.formatType(column),
					strings.Join(f.getColumnAttributes(column), ", "),
					strings.TrimSpace(column.Description),
				},
				normalStyle)
		}
	}

	return nil
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

	if column.PrimaryKey {
		attrs = append(attrs, "PK")
	}
	if column.UniqueKey {
		attrs = append(attrs, "UNIQUE")
	}
	if column.AutoIncremental {
		attrs = append(attrs, "ID")
	}
	if column.Nullable {
		attrs = append(attrs, "nullable")
	}
	if column.DefaultValue != "" {
		attrs = append(attrs, "default="+column.DefaultValue)
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
