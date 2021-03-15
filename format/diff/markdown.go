package diff

import (
	"bytes"
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/mysql"
	"io"
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
	writer        *MarkdownWriter
	option        *Option
	mysqlExporter *mysql.Exporter
}

func NewMarkdownChangeSetWirter(w io.Writer, option *Option) *MarkdownChangeSetWriter {
	return &MarkdownChangeSetWriter{
		writer: NewMarkdownWriter(w),
		option: option,
		mysqlExporter: mysql.NewExporter(nil, &mysql.ExportOption{
			TableFilter:      option.TableFilter,
			UniqueNameSuffix: option.UniqueNameSuffix,
		}),
	}
}

func (w *MarkdownChangeSetWriter) Write(result *Result) error {
	writer := w.writer
	writer.WriteH1("DB Schema changes")

	for _, changeSet := range result.ChangeSets {
		for _, change := range changeSet.Changes {
			switch change.(type) {
			case *CreateTable:
				w.writeCreateTable(change.(*CreateTable))
			case *DropTable:
				w.writeDropTable(change.(*DropTable))
			case *RenameTable:
				w.writeRenameTable(change.(*RenameTable))
			case *UpdatePrimaryKey:
				w.writeUpdatePrimaryKey(change.(*UpdatePrimaryKey))
			case *AddColumn:
				w.writeAddColumn(change.(*AddColumn))
			case *SetNotNullConstraint:
				w.writeSetNotNullConstraint(change.(*SetNotNullConstraint))
			}
		}
	}
	return nil
}

func (w *MarkdownChangeSetWriter) writeCreateTable(change *CreateTable) {
	writer := w.writer
	writer.WriteH2(change.Table.Name)

	buf := new(bytes.Buffer)
	if err := w.mysqlExporter.ExportTable(buf, change.Table); err != nil {
		panic(err)
	}

	writer.WriteSQL(buf.Bytes())
}

func (w *MarkdownChangeSetWriter) writeDropTable(change *DropTable) {
	writer := w.writer
	writer.WriteH2(change.Table.Name)
	sql := fmt.Sprintf("drop table %s;", change.Table.Name)
	writer.WriteSQLs([]string{sql})
}

func (w *MarkdownChangeSetWriter) writeRenameTable(change *RenameTable) {
	writer := w.writer
	writer.WriteH2(change.NewTable.Name)
	sql := fmt.Sprintf("rename table %s to %s;", change.OldTable.Name, change.NewTable.Name)
	writer.WriteSQLs([]string{sql})
}

func (w *MarkdownChangeSetWriter) writeUpdatePrimaryKey(change *UpdatePrimaryKey) {
	writer := w.writer
	writer.WriteH2(change.NewTable.Name)
	sqls := []string{
		fmt.Sprintf("alter table %s drop primary key;", change.NewTable.Name),
		fmt.Sprintf("alter table %s add primary key (%s);", change.NewTable.Name, strings.Join(change.NewTable.PrimaryKeyNameSet().Slice(), ",")),
	}
	writer.WriteSQLs(sqls)
}

func (w *MarkdownChangeSetWriter) writeAddColumn(change *AddColumn) {
	tableName := change.Table.Name
	column := change.Column
	columnType := w.mysqlExporter.ToMysqlColumnType(column)

	writer := w.writer
	writer.WriteH2(tableName)

	// TODO: default value
	var sql string
	if column.NotNull {
		sql = fmt.Sprintf("alter table %s modify add %s %s not null;", tableName, column.Name, columnType)
	} else {
		sql = fmt.Sprintf("alter table %s modify add %s %s;", tableName, column.Name, columnType)
	}
	writer.WriteSQLs([]string{sql})
}

func (w *MarkdownChangeSetWriter) writeSetNotNullConstraint(change *SetNotNullConstraint) {
	tableName := change.Table.Name
	column := change.Column
	columnType := w.mysqlExporter.ToMysqlColumnType(column)

	writer := w.writer
	writer.WriteH2(tableName)

	var sql string
	if column.NotNull {
		sql = fmt.Sprintf("alter table %s modify %s %s not null;", tableName, column.Name, columnType)
	} else {
		sql = fmt.Sprintf("alter table %s modify %s %s;", tableName, column.Name, columnType)
	}
	writer.WriteSQLs([]string{sql})
}
