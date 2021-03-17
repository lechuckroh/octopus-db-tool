package diff

import (
	"bytes"
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/mysql"
	"reflect"
	"strings"
)

type MysqlQueryBuilder struct {
	option        *Option
	mysqlExporter *mysql.Exporter
}

func NewMysqlQueryBuilder(option *Option) *MysqlQueryBuilder {
	return &MysqlQueryBuilder{
		option: option,
		mysqlExporter: mysql.NewExporter(nil, &mysql.ExportOption{
			TableFilter:      option.TableFilter,
			UniqueNameSuffix: option.UniqueNameSuffix,
		}),
	}
}

func (b *MysqlQueryBuilder) ToSQL(change interface{}) ([]string, error) {
	switch change.(type) {
	case *CreateTable:
		return b.toCreateTableSQL(change.(*CreateTable))
	case *DropTable:
		return b.toDropTableSQL(change.(*DropTable))
	case *RenameTable:
		return b.toRenameTableSQL(change.(*RenameTable))
	case *UpdatePrimaryKey:
		return b.toUpdatePrimaryKeySQL(change.(*UpdatePrimaryKey))
	case *CreateUniqueConstraint:
		return b.toCreateUniqueConstraintSQL(change.(*CreateUniqueConstraint))
	case *DropUniqueConstraint:
		return b.toDropUniqueConstraintSQL(change.(*DropUniqueConstraint))
	case *SetTableComment:
		return b.toSetTableCommentSQL(change.(*SetTableComment))
	case *AddColumn:
		return b.toAddColumnSQL(change.(*AddColumn))
	case *DropColumn:
		return b.toDropColumnSQL(change.(*DropColumn))
	case *SetColumnComment:
		return b.toSetColumnCommentSQL(change.(*SetColumnComment))
	case *ChangeColumnType:
		return b.toChangeColumnTypeSQL(change.(*ChangeColumnType))
	case *RenameColumn:
		return b.toRenameColumnSQL(change.(*RenameColumn))
	case *SetNotNullConstraint:
		return b.toSetNotNullConstraintSQL(change.(*SetNotNullConstraint))
	case *SetAutoIncrement:
		return b.toSetAutoIncrementSQL(change.(*SetAutoIncrement))
	case *SetDefaultValue:
		return b.toSetDefaultValueSQL(change.(*SetDefaultValue))
	default:
		return nil, fmt.Errorf("unhandled change type: %v", reflect.TypeOf(change))
	}
}

func (b *MysqlQueryBuilder) toCreateTableSQL(c *CreateTable) ([]string, error) {
	buf := new(bytes.Buffer)
	if err := b.mysqlExporter.ExportTable(buf, c.Table); err != nil {
		return nil, err
	}

	return []string{buf.String()}, nil
}

func (b *MysqlQueryBuilder) toDropTableSQL(c *DropTable) ([]string, error) {
	return []string{
		fmt.Sprintf("DROP TABLE %s;", c.Table.Name),
	}, nil
}

func (b *MysqlQueryBuilder) toRenameTableSQL(c *RenameTable) ([]string, error) {
	return []string{
		fmt.Sprintf("RENAME TABLE %s TO %s;", c.OldTable.Name, c.NewTable.Name),
	}, nil
}

func (b *MysqlQueryBuilder) toUpdatePrimaryKeySQL(c *UpdatePrimaryKey) ([]string, error) {
	table := c.NewTable
	tableName := table.Name
	pkColumnNames := table.PrimaryKeyNameSet().Slice()

	result := []string{
		fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY;", tableName),
	}

	if len(pkColumnNames) > 0 {
		result = append(result,
			fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (%s);", tableName, strings.Join(pkColumnNames, ",")))
	}
	return result, nil
}

func (b *MysqlQueryBuilder) toCreateUniqueConstraintSQL(c *CreateUniqueConstraint) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toDropUniqueConstraintSQL(c *DropUniqueConstraint) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toSetTableCommentSQL(c *SetTableComment) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toAddColumnSQL(c *AddColumn) ([]string, error) {
	tableName := c.Table.Name
	column := c.Column
	columnType := b.mysqlExporter.ToMysqlColumnType(column)
	colConstraints := b.mysqlExporter.ColumnConstraints(column)

	sql := fmt.Sprintf("ALTER TABLE %s MODIFY ADD %s %s", tableName, column.Name, columnType)
	if colConstraints != "" {
		sql += " " + colConstraints
	}
	return []string{sql + ";"}, nil
}

func (b *MysqlQueryBuilder) toDropColumnSQL(c *DropColumn) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toSetColumnCommentSQL(c *SetColumnComment) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toChangeColumnTypeSQL(c *ChangeColumnType) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toRenameColumnSQL(c *RenameColumn) ([]string, error) {
	// TODO
	return nil, nil
}

func (b *MysqlQueryBuilder) toSetNotNullConstraintSQL(c *SetNotNullConstraint) ([]string, error) {
	tableName := c.Table.Name
	column := c.Column
	columnType := b.mysqlExporter.ToMysqlColumnType(column)

	var sql string
	if column.NotNull {
		sql = fmt.Sprintf("ALTER TABLE %s MODIFY %s %s NOT NULL;", tableName, column.Name, columnType)
	} else {
		sql = fmt.Sprintf("ALTER TABLE %s MODIFY %s %s;", tableName, column.Name, columnType)
	}
	return []string{sql}, nil
}

func (b *MysqlQueryBuilder) toSetAutoIncrementSQL(c *SetAutoIncrement) ([]string, error) {
	tableName := c.Table.Name
	column := c.Column
	columnName := column.Name
	columnType := b.mysqlExporter.ToMysqlColumnType(column)

	var sql string
	if c.Column.AutoIncremental {
		sql = strings.Join([]string{
			fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s", tableName, columnName, columnType),
			b.mysqlExporter.ColumnConstraints(column),
		}, " ") + ";"
	} else {
		sql = fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s AUTO_INCREMENT PRIMARY KEY;",
			tableName, columnName, columnType)
	}
	return []string{sql}, nil
}

func (b *MysqlQueryBuilder) toSetDefaultValueSQL(c *SetDefaultValue) ([]string, error) {
	// TODO
	return nil, nil
}
