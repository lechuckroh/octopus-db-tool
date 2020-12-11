package mysql

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/xwb1989/sqlparser"
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
		return c.ImportBytes(bytes)
	}
}

func (c *Importer) ImportFile(filename string) (*octopus.Schema, error) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return c.ImportBytes(data)
	}
}

func (c *Importer) ImportBytes(data []byte) (*octopus.Schema, error) {
	tokens := sqlparser.NewStringTokenizer(string(data))

	tables := make([]*octopus.Table, 0)
	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch stmt.(type) {
		case *sqlparser.DDL:
			ddl := stmt.(*sqlparser.DDL)
			columns := make([]*octopus.Column, 0)

			tableSpec := ddl.TableSpec

			if tableSpec != nil {
				pkSet := util.NewStringSet()
				uniqueSet := util.NewStringSet()
				for _, idx := range tableSpec.Indexes {
					info := idx.Info
					if info.Primary {
						for _, c := range idx.Columns {
							pkSet.Add(c.Column.String())
						}
					} else if info.Unique {
						for _, c := range idx.Columns {
							uniqueSet.Add(c.Column.String())
						}
					}
				}

				for _, col := range ddl.TableSpec.Columns {
					name := col.Name.String()
					nullable := !bool(col.Type.NotNull)
					defaultValue := util.SQLValToString(col.Type.Default, "")
					if nullable && defaultValue == "null" {
						defaultValue = ""
					}
					comment := util.SQLValToString(col.Type.Comment, "")
					columns = append(columns, &octopus.Column{
						Name:            name,
						Type:            c.fromColumnType(col.Type),
						Description:     comment,
						Size:            uint16(util.SQLValToInt(col.Type.Length, 0)),
						Scale:           uint16(util.SQLValToInt(col.Type.Scale, 0)),
						Nullable:        nullable,
						PrimaryKey:      pkSet.Contains(name),
						UniqueKey:       uniqueSet.Contains(name),
						AutoIncremental: bool(col.Type.Autoincrement),
						DefaultValue:    defaultValue,
					})
				}
				tables = append(tables, &octopus.Table{
					Name:    ddl.NewName.Name.String(),
					Columns: columns,
				})
			}
		}
	}

	schema := octopus.Schema{
		Tables: tables,
	}

	return &schema, nil
}

func (c *Importer) fromColumnType(colType sqlparser.ColumnType) string {
	switch colType.Type {
	case "bit":
		fallthrough
	case "bool":
		fallthrough
	case "boolean":
		return octopus.ColTypeBoolean
	case "tinyint":
		fallthrough
	case "smallint":
		fallthrough
	case "mediumint":
		fallthrough
	case "int":
		fallthrough
	case "integer":
		return octopus.ColTypeInt
	case "bigint":
		return octopus.ColTypeLong
	case "decimal":
		return octopus.ColTypeDecimal
	case "float":
		return octopus.ColTypeFloat
	case "double":
		return octopus.ColTypeDouble
	case "binary":
		fallthrough
	case "varbinary":
		fallthrough
	case "char":
		fallthrough
	case "varchar":
		return octopus.ColTypeString
	case "longtext":
		fallthrough
	case "mediumtext":
		fallthrough
	case "tinytext":
		fallthrough
	case "text":
		return octopus.ColTypeText
	case "enum":
		return colType.Type
	case "set":
		return colType.Type
	case "datetime":
		fallthrough
	case "timestamp":
		return octopus.ColTypeDateTime
	case "date":
		return octopus.ColTypeDate
	case "time":
		return octopus.ColTypeTime
	case "year":
		return colType.Type
	case "longblob":
		fallthrough
	case "mediumblob":
		fallthrough
	case "blob":
		return octopus.ColTypeBlob
	case "geometry":
		return colType.Type
	case "point":
		return colType.Type
	case "linestring":
		return colType.Type
	case "polygon":
		return colType.Type
	case "geometrycollection":
		return colType.Type
	case "multilinestring":
		return colType.Type
	case "multipoint":
		return colType.Type
	case "multipolygon":
		return colType.Type
	case "json":
		return colType.Type
	default:
		return colType.Type
	}
}
