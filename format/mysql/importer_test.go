package mysql

import (
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestMysqlImporter_ImportSql(t *testing.T) {
	Convey("import", t, func() {
		sql := strings.Join([]string{
			"CREATE TABLE IF NOT EXISTS `Table` (",
			"id bigint NOT NULL AUTO_INCREMENT,",
			"name varchar(20) NOT NULL DEFAULT 'noname',",
			"postal_code char(6),",
			"age tinyint NOT NULL DEFAULT 0,",
			"float32 float(10,3) NOT NULL DEFAULT 0,",
			"float64 double(20,5) NOT NULL DEFAULT 0,",
			"int_value int comment 'int value',",
			"bool1 bit(1) NOT NULL DEFAULT 1,",
			"bool2 tinyint(1) NOT NULL DEFAULT 1,",
			"text1 tinytext,",
			"text2 text,",
			"text3 mediumtext,",
			"text4 longtext,",
			"blob1 tinyblob,",
			"blob2 blob,",
			"blob3 mediumblob,",
			"blob4 longblob,",
			"created_at datetime,",
			"PRIMARY KEY (`id`),",
			"UNIQUE KEY `Table_UNIQUE` (`name`)",
			");",
		}, "\n")

		mysql := Importer{}

		// read sql
		schema, err := mysql.ImportSql(sql)
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
							Type:            octopus.ColTypeInt64,
							PrimaryKey:      true,
							Nullable:        false,
							AutoIncremental: true,
						},
						{
							Name:         "name",
							Type:         octopus.ColTypeVarchar,
							Size:         20,
							Nullable:     false,
							UniqueKey:    true,
							DefaultValue: "noname",
						},
						{
							Name:         "postal_code",
							Type:         octopus.ColTypeChar,
							Size:         6,
							Nullable:     true,
						},
						{
							Name:         "age",
							Type:         octopus.ColTypeInt8,
							Size:         0,
							Nullable:     false,
							DefaultValue: "0",
						},
						{
							Name:         "float32",
							Type:         octopus.ColTypeFloat,
							Size:         10,
							Scale:        3,
							Nullable:     false,
							DefaultValue: "0",
						},
						{
							Name:         "float64",
							Type:         octopus.ColTypeDouble,
							Size:         20,
							Scale:        5,
							Nullable:     false,
							DefaultValue: "0",
						},
						{
							Name:        "int_value",
							Type:        octopus.ColTypeInt32,
							Size:        0,
							Nullable:    true,
							Description: "int value",
						},
						{
							Name:         "bool1",
							Type:         octopus.ColTypeBoolean,
							Nullable:     false,
							DefaultValue: "1",
						},
						{
							Name:         "bool2",
							Type:         octopus.ColTypeBoolean,
							Nullable:     false,
							DefaultValue: "1",
						},
						{
							Name:     "text1",
							Type:     octopus.ColTypeText8,
							Nullable: true,
						},
						{
							Name:     "text2",
							Type:     octopus.ColTypeText16,
							Nullable: true,
						},
						{
							Name:     "text3",
							Type:     octopus.ColTypeText24,
							Nullable: true,
						},
						{
							Name:     "text4",
							Type:     octopus.ColTypeText32,
							Nullable: true,
						},
						{
							Name:     "blob1",
							Type:     octopus.ColTypeBlob8,
							Nullable: true,
						},
						{
							Name:     "blob2",
							Type:     octopus.ColTypeBlob16,
							Nullable: true,
						},
						{
							Name:     "blob3",
							Type:     octopus.ColTypeBlob24,
							Nullable: true,
						},
						{
							Name:     "blob4",
							Type:     octopus.ColTypeBlob32,
							Nullable: true,
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
