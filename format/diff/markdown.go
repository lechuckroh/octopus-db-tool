package diff

import (
	"io"
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

	for _, changeSet := range result.ChangeSets {
		for _, change := range changeSet.Changes {
			if sqls, err := w.queryBuilder.ToSQL(change); err != nil {
				return err
			} else {
				writer.WriteSQLs(sqls)
			}
		}
	}
	return nil
}
