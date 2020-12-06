package main

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMysqlExport_Export(t *testing.T) {
	schema := &Schema{
		Tables: []*Table{
			{
				Name: "Table",
				Columns: []*Column{
					{
						Name:            "id",
						Type:            ColTypeLong,
						PrimaryKey:      true,
						Nullable:        false,
						AutoIncremental: true,
					},
					{
						Name:         "name",
						Type:         ColTypeString,
						Size:         20,
						Nullable:     false,
						UniqueKey:    true,
						DefaultValue: "noname",
					},
					{
						Name:     "created_at",
						Type:     ColTypeDateTime,
						Nullable: true,
					},
				},
			},
		},
	}

	Convey("export", t, func() {
		// given:
		exporter := MysqlExport{
			schema: schema,
		}
		option := MysqlExportOption{
			UniqueNameSuffix: "_uq",
		}
		expected := `CREATE TABLE IF NOT EXISTS ` + "Table" + ` (
  id bigint NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL DEFAULT 'noname',
  created_at datetime,
  PRIMARY KEY (` + "`id`" + `),
  UNIQUE KEY ` + "`Table_uq` (`name`)" + `
);
`

		// when:
		buf := new(bytes.Buffer)
		err := exporter.Export(buf, &option)

		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, expected)
	})
}
