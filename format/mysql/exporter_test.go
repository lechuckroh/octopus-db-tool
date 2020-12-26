package mysql

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMysqlExport_Export(t *testing.T) {
	schema := &octopus.Schema{
		Tables: []*octopus.Table{
			{
				Name: "Table",
				Columns: []*octopus.Column{
					{
						Name:            "id",
						Type:            octopus.ColTypeLong,
						PrimaryKey:      true,
						Nullable:        false,
						AutoIncremental: true,
					},
					{
						Name:         "name",
						Type:         octopus.ColTypeString,
						Size:         20,
						Nullable:     false,
						UniqueKey:    true,
						DefaultValue: "noname",
					},
					{
						Name:     "created_at",
						Type:     octopus.ColTypeDateTime,
						Nullable: true,
					},
				},
			},
		},
	}

	Convey("export", t, func() {
		// given:
		option := ExportOption{
			UniqueNameSuffix: "_uq",
		}
		exporter := Exporter{
			schema: schema,
			option: &option,
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
		err := exporter.Export(buf)

		So(err, ShouldBeNil)
		So(buf.String(), ShouldEqual, expected)
	})
}