package diff

import (
	"io"
)

type FlywayWriter struct {
	writer io.Writer
}

func (w *FlywayWriter) Write(data []byte) {
	if _, err := w.writer.Write(data); err != nil {
		panic(err)
	}
}

func (w *FlywayWriter) WriteLine(s string) {
	if _, err := w.writer.Write([]byte(s + "\n")); err != nil {
		panic(err)
	}
}

func NewFlywayWriter(writer io.Writer) *FlywayWriter {
	return &FlywayWriter{writer: writer}
}

type FlywayChangeSetWriter struct {
	writer       *FlywayWriter
	option       *Option
	queryBuilder *MysqlQueryBuilder
}

func NewFlywayChangeSetWirter(w io.Writer, option *Option) *FlywayChangeSetWriter {
	// TODO: add option to write multiple files by changeSet
	return &FlywayChangeSetWriter{
		writer:       NewFlywayWriter(w),
		option:       option,
		queryBuilder: NewMysqlQueryBuilder(option),
	}
}

func (w *FlywayChangeSetWriter) Write(result *Result) error {
	for _, changeSet := range result.ChangeSets {
		for _, change := range changeSet.Changes {
			if sqls, err := w.queryBuilder.ToSQL(change); err != nil {
				return err
			} else {
				for _, sql := range sqls {
					w.writer.WriteLine(sql)
				}
			}
		}
		w.writer.WriteLine("")
	}
	return nil
}
