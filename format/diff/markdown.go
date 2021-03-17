package diff

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/util"
	"io"
	"reflect"
	"sort"
	"strings"
)

type MarkdownWriter struct {
	writer io.Writer
}

func (w *MarkdownWriter) Write(data []byte) {
	if _, err := w.writer.Write(data); err != nil {
		panic(err)
	}
}

func (w *MarkdownWriter) WriteLine(s string) {
	if _, err := w.writer.Write([]byte(s + "\n")); err != nil {
		panic(err)
	}
}

func (w *MarkdownWriter) WriteH1(s string) {
	w.WriteLine("# " + s)
}

func (w *MarkdownWriter) WriteH2(s string) {
	w.WriteLine("## " + s)
}

func (w *MarkdownWriter) WriteH3(s string) {
	w.WriteLine("### " + s)
}

func (w *MarkdownWriter) WriteCodes(lang string, lines []string) {
	w.WriteLine("```" + lang)
	for _, line := range lines {
		w.WriteLine(line)
	}
	w.WriteLine("```")
}

func (w *MarkdownWriter) WriteCode(lang string, data []byte) {
	w.WriteLine("```" + lang)
	w.Write(data)
	w.WriteLine("")
	w.WriteLine("```")
}

func (w *MarkdownWriter) WriteSQL(data []byte) {
	w.WriteCode("sql", data)
}

func (w *MarkdownWriter) WriteSQLs(lines []string) {
	w.WriteCodes("sql", lines)
}

func NewMarkdownWriter(writer io.Writer) *MarkdownWriter {
	return &MarkdownWriter{writer: writer}
}

type MarkdownChangeSetWriter struct {
	writer       *MarkdownWriter
	option       *Option
	queryBuilder *MysqlQueryBuilder
}

func NewMarkdownChangeSetWirter(w io.Writer, option *Option) *MarkdownChangeSetWriter {
	return &MarkdownChangeSetWriter{
		writer:       NewMarkdownWriter(w),
		option:       option,
		queryBuilder: NewMysqlQueryBuilder(option),
	}
}

func (w *MarkdownChangeSetWriter) Write(result *Result) error {
	writer := w.writer
	writer.WriteH1("DB Schema changes")
	writer.WriteLine("")
	writer.WriteLine(fmt.Sprintf("* from: `%s`", result.From.Version))
	writer.WriteLine(fmt.Sprintf("* to: `%s`", result.To.Version))

	changesByTable := util.NewMultiMap()

	for _, changeSet := range result.ChangeSets {
		for _, change := range changeSet.Changes {
			changesByTable.Put(change.DepTable().Name, change)
		}
	}

	keys := changesByTable.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i].(string), keys[j].(string)) < 0
	})

	for _, key := range keys {
		tableName := key.(string)

		writer.WriteLine("")
		writer.WriteH2(tableName)

		changes, _ := changesByTable.Get(key)
		for _, value := range changes {
			change := value.(Change)
			if line, err := w.toMarkdownLine(change); err != nil {
				return err
			} else {
				writer.WriteLine("* " + line)
			}
		}
	}

	return nil
}

func (w *MarkdownChangeSetWriter) toMarkdownLine(change Change) (string, error) {
	switch change.(type) {
	case *CreateTable:
		return w.toCreateTableMD(change.(*CreateTable))
	case *DropTable:
		return w.toDropTableMD(change.(*DropTable))
	case *RenameTable:
		return w.toRenameTableMD(change.(*RenameTable))
	case *UpdatePrimaryKey:
		return w.toUpdatePrimaryKeyMD(change.(*UpdatePrimaryKey))
	case *CreateUniqueConstraint:
		return w.toCreateUniqueConstraintMD(change.(*CreateUniqueConstraint))
	case *DropUniqueConstraint:
		return w.toDropUniqueConstraintMD(change.(*DropUniqueConstraint))
	case *SetTableComment:
		return w.toSetTableCommentMD(change.(*SetTableComment))
	case *AddColumn:
		return w.toAddColumnMD(change.(*AddColumn))
	case *DropColumn:
		return w.toDropColumnMD(change.(*DropColumn))
	case *SetColumnComment:
		return w.toSetColumnCommentMD(change.(*SetColumnComment))
	case *ChangeColumnType:
		return w.toChangeColumnTypeMD(change.(*ChangeColumnType))
	case *RenameColumn:
		return w.toRenameColumnMD(change.(*RenameColumn))
	case *SetNotNullConstraint:
		return w.toSetNotNullConstraintMD(change.(*SetNotNullConstraint))
	case *SetAutoIncrement:
		return w.toSetAutoIncrementMD(change.(*SetAutoIncrement))
	case *SetDefaultValue:
		return w.toSetDefaultValueMD(change.(*SetDefaultValue))
	default:
		return "", fmt.Errorf("unhandled change type: %v", reflect.TypeOf(change))
	}
}

func (w *MarkdownChangeSetWriter) toCreateTableMD(c *CreateTable) (string, error) {
	return "create table", nil
}

func (w *MarkdownChangeSetWriter) toDropTableMD(c *DropTable) (string, error) {
	return "drop table", nil
}

func (w *MarkdownChangeSetWriter) toRenameTableMD(c *RenameTable) (string, error) {
	return fmt.Sprintf("rename table `%s` → `%s`", c.OldTable.Name, c.NewTable.Name), nil
}

func (w *MarkdownChangeSetWriter) toUpdatePrimaryKeyMD(c *UpdatePrimaryKey) (string, error) {
	oldPKColumnNames := strings.Join(c.OldTable.PrimaryKeyNameSet().Slice(), ", ")
	newPKColumnNames := strings.Join(c.NewTable.PrimaryKeyNameSet().Slice(), ", ")
	return fmt.Sprintf("update PK: [%s] → [%s]", oldPKColumnNames, newPKColumnNames), nil
}

func (w *MarkdownChangeSetWriter) toCreateUniqueConstraintMD(c *CreateUniqueConstraint) (string, error) {
	newUniqueColumnNames := strings.Join(c.Table.UniqueKeyNameSet().Slice(), ", ")
	return fmt.Sprintf("create unique constraint: [%s]", newUniqueColumnNames), nil
}

func (w *MarkdownChangeSetWriter) toDropUniqueConstraintMD(c *DropUniqueConstraint) (string, error) {
	oldUniqueColumnNames := strings.Join(c.Table.UniqueKeyNameSet().Slice(), ", ")
	return fmt.Sprintf("drop unique constraint: [%s]", oldUniqueColumnNames), nil
}

func (w *MarkdownChangeSetWriter) toSetTableCommentMD(c *SetTableComment) (string, error) {
	return fmt.Sprintf("update table comment: `%s`", c.Table.Description), nil
}

func (w *MarkdownChangeSetWriter) toAddColumnMD(c *AddColumn) (string, error) {
	return fmt.Sprintf("add column: `%s`", c.Column.Name), nil
}

func (w *MarkdownChangeSetWriter) toDropColumnMD(c *DropColumn) (string, error) {
	return fmt.Sprintf("drop column: `%s`", c.ColumnName), nil
}

func (w *MarkdownChangeSetWriter) toSetColumnCommentMD(c *SetColumnComment) (string, error) {
	return fmt.Sprintf("update column comment: `%s` :`%s`", c.Column.Name, c.Column.Description), nil
}

func (w *MarkdownChangeSetWriter) toChangeColumnTypeMD(c *ChangeColumnType) (string, error) {
	return fmt.Sprintf("change column type: `%s` :`%s` → `%s`",
		c.NewColumn.Name, c.OldColumn.Format(), c.NewColumn.Format()), nil
}

func (w *MarkdownChangeSetWriter) toRenameColumnMD(c *RenameColumn) (string, error) {
	return fmt.Sprintf("rename column: `%s` → `%s`", c.OldColumn.Name, c.NewColumn.Name), nil
}

func (w *MarkdownChangeSetWriter) toSetNotNullConstraintMD(c *SetNotNullConstraint) (string, error) {
	if c.Column.NotNull {
		return fmt.Sprintf("set not null: `%s`", c.Column.Name), nil
	} else {
		return fmt.Sprintf("set nullable: `%s`", c.Column.Name), nil
	}
}

func (w *MarkdownChangeSetWriter) toSetAutoIncrementMD(c *SetAutoIncrement) (string, error) {
	if c.Column.AutoIncremental {
		return fmt.Sprintf("set auto incremental: `%s`", c.Column.Name), nil
	} else {
		return fmt.Sprintf("remove auto incremental: `%s`", c.Column.Name), nil
	}
}

func (w *MarkdownChangeSetWriter) toSetDefaultValueMD(c *SetDefaultValue) (string, error) {
	defaultValue := c.Column.DefaultValue
	if defaultValue == "" {
		return fmt.Sprintf("remove default value: `%s`", c.Column.Name), nil
	} else {
		return fmt.Sprintf("update default value: `%s` : `%s`", c.Column.Name, defaultValue), nil
	}
}
