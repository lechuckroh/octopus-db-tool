package mysql

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"strings"
	"testing"
)

func TestMysqlExport_Export(t *testing.T) {
	schema := &octopus.Schema{
		Tables: []*octopus.Table{
			{
				Name:  "Table",
				Group: "group1",
				Indices: []*octopus.Index{
					{Name: "idx_name", Columns: []string{"name"}},
					{Name: "idx_ints", Columns: []string{"age", "int2", "int3"}},
				},
				Columns: []*octopus.Column{
					{
						Name:            "id",
						Type:            octopus.ColTypeInt64,
						PrimaryKey:      true,
						NotNull:         true,
						AutoIncremental: true,
					},
					{
						Name:         "name",
						Type:         octopus.ColTypeVarchar,
						Size:         20,
						NotNull:      true,
						UniqueKey:    true,
						DefaultValue: "noname",
						Description:  "group name",
					},
					{
						Name:    "company_id",
						Type:    octopus.ColTypeInt64,
						NotNull: true,
						Ref: &octopus.Reference{
							Table:        "company",
							Column:       "id",
							Relationship: octopus.RefManyToOne,
						},
					},
					{
						Name:      "postal_code",
						Type:      octopus.ColTypeChar,
						Size:      6,
						UniqueKey: true,
					},
					{
						Name:         "age",
						Type:         octopus.ColTypeInt8,
						Size:         0,
						NotNull:      true,
						DefaultValue: "0",
					},
					{
						Name: "int2",
						Type: octopus.ColTypeInt16,
					},
					{
						Name: "int3",
						Type: octopus.ColTypeInt24,
					},
					{
						Name: "int4",
						Type: octopus.ColTypeInt64,
					},
					{
						Name:  "decimal1",
						Type:  octopus.ColTypeDecimal,
						Size:  10,
						Scale: 3,
					},
					{
						Name:         "float32",
						Type:         octopus.ColTypeFloat,
						Size:         10,
						Scale:        3,
						NotNull:      true,
						DefaultValue: "0",
					},
					{
						Name:         "float64",
						Type:         octopus.ColTypeDouble,
						Size:         20,
						Scale:        5,
						NotNull:      true,
						DefaultValue: "0",
					},
					{
						Name:        "int_value",
						Type:        octopus.ColTypeInt32,
						Size:        0,
						Description: "int value",
					},
					{
						Name:         "bool1",
						Type:         octopus.ColTypeBoolean,
						NotNull:      true,
						DefaultValue: "1",
					},
					{
						Name:         "bool2",
						Type:         octopus.ColTypeBoolean,
						NotNull:      true,
						DefaultValue: "1",
					},
					{
						Name: "bit1",
						Type: octopus.ColTypeBit,
						Size: 3,
					},
					{
						Name: "bin1",
						Type: octopus.ColTypeBinary,
					},
					{
						Name: "bin2",
						Type: octopus.ColTypeVarbinary,
					},
					{
						Name: "text1",
						Type: octopus.ColTypeText8,
					},
					{
						Name: "text2",
						Type: octopus.ColTypeText16,
					},
					{
						Name: "text3",
						Type: octopus.ColTypeText24,
					},
					{
						Name: "text4",
						Type: octopus.ColTypeText32,
					},
					{
						Name: "blob1",
						Type: octopus.ColTypeBlob8,
					},
					{
						Name: "blob2",
						Type: octopus.ColTypeBlob16,
					},
					{
						Name: "blob3",
						Type: octopus.ColTypeBlob24,
					},
					{
						Name: "blob4",
						Type: octopus.ColTypeBlob32,
					},
					{
						Name: "geo",
						Type: octopus.ColTypeGeometry,
					},
					{
						Name:   "enum1",
						Type:   octopus.ColTypeEnum,
						Values: []string{"1", "2", "3"},
					},
					{
						Name: "point1",
						Type: octopus.ColTypePoint,
					},
					{
						Name:   "set1",
						Type:   octopus.ColTypeSet,
						Values: []string{"1", "2", "3"},
					},
					{
						Name: "json1",
						Type: octopus.ColTypeJSON,
					},
					{
						Name: "date1",
						Type: octopus.ColTypeDate,
					},
					{
						Name: "datetime1",
						Type: octopus.ColTypeDateTime,
					},
					{
						Name: "time1",
						Type: octopus.ColTypeTime,
					},
					{
						Name: "year1",
						Type: octopus.ColTypeYear,
					},
					{
						Name:     "on_update_text",
						Type:     octopus.ColTypeVarchar,
						Size:     30,
						OnUpdate: "updated",
					},
					{
						Name:     "on_update_int",
						Type:     octopus.ColTypeInt16,
						OnUpdate: "1",
					},
					{
						Name:         "created_at",
						Type:         octopus.ColTypeDateTime,
						NotNull:      true,
						DefaultValue: "fn::CURRENT_TIMESTAMP",
					},
					{
						Name:         "updated_at",
						Type:         octopus.ColTypeDateTime,
						NotNull:      true,
						DefaultValue: "fn::CURRENT_TIMESTAMP",
						OnUpdate:     "fn::CURRENT_TIMESTAMP",
					},
					{
						Name: "unknown",
						Type: "invalidtype",
					},
				},
			},
			{
				Name:  "Table2",
				Group: "group2",
				Columns: []*octopus.Column{
					{
						Name:            "id",
						Type:            octopus.ColTypeInt64,
						PrimaryKey:      true,
						NotNull:         true,
						AutoIncremental: true,
					},
				},
			},
		},
	}

	Convey("export", t, func() {
		// given:
		option := ExportOption{
			UniqueNameSuffix: "_uq",
			TableFilter:      octopus.GetTableFilterFn("group1"),
		}
		exporter := Exporter{
			schema: schema,
			option: &option,
		}
		expected := strings.Join([]string{
			"CREATE TABLE IF NOT EXISTS Table (",
			"  id bigint NOT NULL AUTO_INCREMENT,",
			"  name varchar(20) NOT NULL DEFAULT 'noname' COMMENT 'group name',",
			"  company_id bigint NOT NULL,",
			"  postal_code char(6),",
			"  age tinyint NOT NULL DEFAULT 0,",
			"  int2 smallint,",
			"  int3 mediumint,",
			"  int4 bigint,",
			"  decimal1 decimal(10,3),",
			"  float32 float(10,3) NOT NULL DEFAULT 0,",
			"  float64 double(20,5) NOT NULL DEFAULT 0,",
			"  int_value int COMMENT 'int value',",
			"  bool1 bit(1) NOT NULL DEFAULT 1,",
			"  bool2 bit(1) NOT NULL DEFAULT 1,",
			"  bit1 bit(3),",
			"  bin1 binary,",
			"  bin2 varbinary,",
			"  text1 tinytext,",
			"  text2 text,",
			"  text3 mediumtext,",
			"  text4 longtext,",
			"  blob1 tinyblob,",
			"  blob2 blob,",
			"  blob3 mediumblob,",
			"  blob4 longblob,",
			"  geo geometry,",
			"  enum1 enum('1', '2', '3'),",
			"  point1 point,",
			"  set1 set('1', '2', '3'),",
			"  json1 json,",
			"  date1 date,",
			"  datetime1 datetime,",
			"  time1 time,",
			"  year1 year,",
			"  on_update_text varchar(30) ON UPDATE 'updated',",
			"  on_update_int smallint ON UPDATE 1,",
			"  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),",
			"  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),",
			"  unknown invalidtype,",
			"  PRIMARY KEY (`id`),",
			"  UNIQUE KEY `Table_uq` (`name`, `postal_code`),",
			"  INDEX `idx_name` (`name`),",
			"  INDEX `idx_ints` (`age`, `int2`, `int3`)",
			");",
			"",
		}, "\n")

		// when:
		buf := new(bytes.Buffer)
		err := exporter.Export(buf)

		So(err, ShouldBeNil)
		if diff := cmp.Diff(expected, buf.String()); diff != "" {
			log.Println(diff)
		}
		So(buf.String(), ShouldEqual, expected)
	})
}
