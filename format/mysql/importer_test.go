package mysql

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestMysql_ToSchema(t *testing.T) {
	Convey("ToSchema", t, func() {
		sql := strings.Join([]string{
			"CREATE TABLE IF NOT EXISTS `Table` (",
			"id bigint NOT NULL AUTO_INCREMENT,",
			"name varchar(20) NOT NULL DEFAULT 'noname',",
			"bool1 bit(1) NOT NULL DEFAULT 1,",
			"created_at datetime,",
			"PRIMARY KEY (`id`),",
			"UNIQUE KEY `Table_UNIQUE` (`name`)",
			");",
		}, "\n")

		mysql := Importer{}

		// read sql
		schema, err := mysql.ImportBytes([]byte(sql))
		if err != nil {
			t.Error(err)
		}

		expected := &octopus.Schema{
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
							Name:         "bool1",
							Type:         octopus.ColTypeBoolean,
							Size:         1,
							Nullable:     false,
							DefaultValue: "1",
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

		So(schema, ShouldResemble, expected)
	})
}
