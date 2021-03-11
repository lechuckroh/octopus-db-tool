package mysql

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	_ "github.com/pingcap/parser/test_driver"
	"io"
	"io/ioutil"
	"strings"
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
		return c.ImportSql(string(bytes))
	}
}

func (c *Importer) ImportFile(filename string) (*octopus.Schema, error) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return c.ImportSql(string(data))
	}
}

func (c *Importer) ImportSql(sql string) (*octopus.Schema, error) {
	p := parser.New()

	charset := ""
	collation := ""
	stmtNodes, _, err := p.Parse(sql, charset, collation)
	if err != nil {
		return nil, err
	}

	tableX := TableX{}
	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(&tableX)
	}

	return &octopus.Schema{
		Tables: tableX.tables,
	}, nil
}

// TableX is TableExtractor
type TableX struct {
	tables []*octopus.Table
}

func (x *TableX) Enter(in ast.Node) (ast.Node, bool) {
	if createTableStmt, ok := in.(*ast.CreateTableStmt); ok {
		tableName := createTableStmt.Table.Name.String()

		// constraints
		pkSet := util.NewStringSet()
		uniqSet := util.NewStringSet()
		var indices []*octopus.Index
		for _, cst := range createTableStmt.Constraints {
			switch cst.Tp {
			case ast.ConstraintPrimaryKey:
				for _, key := range cst.Keys {
					pkSet.Add(key.Column.Name.String())
				}
			case ast.ConstraintKey:
				break
			case ast.ConstraintIndex:
				var idxCols []string
				for _, key := range cst.Keys {
					idxCols = append(idxCols, key.Column.Name.String())
				}
				indices = append(indices, &octopus.Index{Name: cst.Name, Columns: idxCols})
				break
			case ast.ConstraintUniq:
				for _, key := range cst.Keys {
					uniqSet.Add(key.Column.Name.String())
				}
			case ast.ConstraintUniqKey:
				break
			case ast.ConstraintUniqIndex:
				break
			case ast.ConstraintForeignKey:
				break
			case ast.ConstraintFulltext:
				break
			case ast.ConstraintCheck:
				break
			}
		}

		// columns
		var columns []*octopus.Column
		for _, colDef := range createTableStmt.Cols {
			column := x.column(colDef, pkSet, uniqSet)
			columns = append(columns, column)
		}

		table := &octopus.Table{
			Name:    tableName,
			Columns: columns,
			Indices: indices,
		}
		x.tables = append(x.tables, table)
		return in, true
	}
	return in, false
}

func (x *TableX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func (x *TableX) column(
	colDef *ast.ColumnDef,
	pkSet *util.StringSet,
	uniqSet *util.StringSet,
) *octopus.Column {
	name := colDef.Name.String()
	colType, colLength, colScale := x.columnType(colDef)
	column := octopus.Column{
		Name:       name,
		Type:       colType,
		Size:       colLength,
		Scale:      colScale,
		PrimaryKey: pkSet.Contains(name),
		UniqueKey:  uniqSet.Contains(name),
		NotNull:    false,
		Values:     x.columnValues(colType, colDef),
	}

	for _, colOption := range colDef.Options {
		switch colOption.Tp {
		case ast.ColumnOptionPrimaryKey:
			column.PrimaryKey = true
		case ast.ColumnOptionNotNull:
			column.NotNull = true
		case ast.ColumnOptionAutoIncrement:
			column.AutoIncremental = true
			column.NotNull = true
		case ast.ColumnOptionDefaultValue:
			switch colOption.Expr.(type) {
			case ast.ValueExpr:
				column.SetDefaultValue(colOption.Expr.(ast.ValueExpr).GetValue())
			case *ast.FuncCallExpr:
				column.SetDefaultValueFn(colOption.Expr.(*ast.FuncCallExpr).FnName.String())
			default:
				fmt.Printf("unhandled default value. column: %s, expr: %v", name, colOption.Expr)
			}
		case ast.ColumnOptionNull:
			column.NotNull = false
		case ast.ColumnOptionOnUpdate:
			switch colOption.Expr.(type) {
			case ast.ValueExpr:
				column.SetOnUpdate(colOption.Expr.(ast.ValueExpr).GetValue())
			case *ast.FuncCallExpr:
				column.SetOnUpdateFn(colOption.Expr.(*ast.FuncCallExpr).FnName.String())
			default:
				fmt.Printf("unhandled default value. column: %s, expr: %v", name, colOption.Expr)
			}
		case ast.ColumnOptionFulltext:
			break
		case ast.ColumnOptionComment:
			valueExpr := colOption.Expr.(ast.ValueExpr)
			column.Description = fmt.Sprintf("%v", valueExpr.GetValue())
		case ast.ColumnOptionGenerated:
			break
		case ast.ColumnOptionReference:
			break
		case ast.ColumnOptionCollate:
			break
		case ast.ColumnOptionCheck:
			break
		case ast.ColumnOptionColumnFormat:
			break
		case ast.ColumnOptionStorage:
			break
		case ast.ColumnOptionAutoRandom:
			break
		}
	}
	return &column
}

func (x *TableX) columnValues(colType string, colDef *ast.ColumnDef) []string {
	switch colType {
	case octopus.ColTypeEnum:
		return colDef.Tp.Elems
	case octopus.ColTypeSet:
		return colDef.Tp.Elems
	}
	return nil
}

func (x *TableX) columnType(colDef *ast.ColumnDef) (string, uint16, uint16) {
	var size, scale uint16
	if colDef.Tp.Flen > 0 {
		size = uint16(colDef.Tp.Flen)
	}
	if colDef.Tp.Decimal > 0 {
		scale = uint16(colDef.Tp.Decimal)
	}

	switch colDef.Tp.Tp {
	case mysql.TypeDecimal:
		return octopus.ColTypeDecimal, size, scale
	case mysql.TypeTiny:
		if colDef.Tp.Flen == 1 {
			return octopus.ColTypeBoolean, 0, 0
		} else {
			return octopus.ColTypeInt8, size, 0
		}
	case mysql.TypeShort:
		return octopus.ColTypeInt16, size, 0
	case mysql.TypeLong:
		return octopus.ColTypeInt32, size, 0
	case mysql.TypeFloat:
		return octopus.ColTypeFloat, size, scale
	case mysql.TypeDouble:
		return octopus.ColTypeDouble, size, scale
	case mysql.TypeNull:
		return colDef.Tp.InfoSchemaStr(), size, scale
	case mysql.TypeTimestamp:
		return octopus.ColTypeDateTime, 0, 0
	case mysql.TypeLonglong:
		return octopus.ColTypeInt64, size, 0
	case mysql.TypeInt24:
		return octopus.ColTypeInt24, size, 0
	case mysql.TypeDate:
		return octopus.ColTypeDate, 0, 0
	case mysql.TypeDuration:
		return octopus.ColTypeTime, 0, 0
	case mysql.TypeDatetime:
		return octopus.ColTypeDateTime, 0, 0
	case mysql.TypeYear:
		return octopus.ColTypeYear, 0, 0
	case mysql.TypeNewDate:
		return octopus.ColTypeDate, 0, 0
	case mysql.TypeVarchar:
		return octopus.ColTypeVarchar, size, 0
	case mysql.TypeBit:
		if colDef.Tp.Flen == 1 {
			return octopus.ColTypeBoolean, 0, 0
		} else {
			return octopus.ColTypeBit, size, 0
		}
	case mysql.TypeJSON:
		return octopus.ColTypeJSON, 0, 0
	case mysql.TypeNewDecimal:
		return octopus.ColTypeDecimal, size, scale
	case mysql.TypeEnum:
		return octopus.ColTypeEnum, 0, 0
	case mysql.TypeSet:
		return octopus.ColTypeSet, 0, 0
	case mysql.TypeTinyBlob:
		if strings.ToLower(colDef.Tp.String()) == "tinytext" {
			return octopus.ColTypeText8, 0, 0
		} else {
			return octopus.ColTypeBlob8, 0, 0
		}
	case mysql.TypeMediumBlob:
		if strings.ToLower(colDef.Tp.String()) == "mediumtext" {
			return octopus.ColTypeText24, 0, 0
		} else {
			return octopus.ColTypeBlob24, 0, 0
		}
	case mysql.TypeLongBlob:
		if strings.ToLower(colDef.Tp.String()) == "longtext" {
			return octopus.ColTypeText32, 0, 0
		} else {
			return octopus.ColTypeBlob32, 0, 0
		}
	case mysql.TypeBlob:
		if strings.ToLower(colDef.Tp.String()) == "text" {
			return octopus.ColTypeText16, 0, 0
		} else {
			return octopus.ColTypeBlob16, 0, 0
		}
	case mysql.TypeString:
		if colDef.Tp.Flag == mysql.BinaryFlag {
			return octopus.ColTypeBinary, size, 0
		}
		return octopus.ColTypeChar, size, 0
	case mysql.TypeGeometry:
		return octopus.ColTypeGeometry, 0, 0
	default:
		return colDef.Tp.InfoSchemaStr(), size, scale
	}
}
