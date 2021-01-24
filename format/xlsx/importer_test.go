package xlsx

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

var testSchema = &octopus.Schema{
	Tables: []*octopus.Table{
		{
			Name: "user",
			Columns: []*octopus.Column{
				{
					Name:            "id",
					Type:            octopus.ColTypeInt64,
					PrimaryKey:      true,
					AutoIncremental: true,
				},
				{
					Name:        "name",
					Type:        octopus.ColTypeVarchar,
					Size:        100,
					UniqueKey:   true,
					Description: "user table",
				},
				{
					Name:    "group_id",
					Type:    octopus.ColTypeInt64,
					NotNull: true,
					Ref: &octopus.Reference{
						Table:        "group",
						Column:       "id",
						Relationship: octopus.RefManyToOne,
					},
				},
				{
					Name:    "1_decimal",
					Type:    octopus.ColTypeDecimal,
					Size:    20,
					Scale:   5,
					NotNull: true,
				},
				{
					Name:  "2_decimal",
					Type:  octopus.ColTypeDecimal,
					Size:  20,
					Scale: 5,
				},
				{
					Name: "bool_1",
					Type: octopus.ColTypeBoolean,
				},
				{
					Name:    "bool2",
					Type:    octopus.ColTypeBoolean,
					NotNull: true,
				},
				{
					Name:    "int1_notnull",
					Type:    octopus.ColTypeInt8,
					NotNull: true,
				},
				{
					Name: "int1_null",
					Type: octopus.ColTypeInt8,
				},
				{
					Name:    "int2_notnull",
					Type:    octopus.ColTypeInt16,
					NotNull: true,
				},
				{
					Name: "int2_null",
					Type: octopus.ColTypeInt16,
				},
				{
					Name:    "int3_notnull",
					Type:    octopus.ColTypeInt24,
					NotNull: true,
				},
				{
					Name: "int3_null",
					Type: octopus.ColTypeInt24,
				},
				{
					Name:    "int4_notnull",
					Type:    octopus.ColTypeInt32,
					NotNull: true,
				},
				{
					Name: "int4_null",
					Type: octopus.ColTypeInt32,
				},
				{
					Name:    "int5_notnull",
					Type:    octopus.ColTypeInt64,
					NotNull: true,
				},
				{
					Name: "int5_null",
					Type: octopus.ColTypeInt64,
				},
				{
					Name:    "float_notnull",
					Type:    octopus.ColTypeFloat,
					Size:    20,
					Scale:   5,
					NotNull: true,
				},
				{
					Name:  "float_null",
					Type:  octopus.ColTypeFloat,
					Size:  20,
					Scale: 10,
				},
				{
					Name:    "double_notnull",
					Type:    octopus.ColTypeDouble,
					Size:    20,
					Scale:   5,
					NotNull: true,
				},
				{
					Name:  "double_null",
					Type:  octopus.ColTypeDouble,
					Size:  20,
					Scale: 10,
				},
				{
					Name:         "created_at",
					Type:         octopus.ColTypeDateTime,
					DefaultValue: "fn::CURRENT_TIMESTAMP",
				},
				{
					Name:         "updated_at",
					Type:         octopus.ColTypeDateTime,
					DefaultValue: "fn::CURRENT_TIMESTAMP",
					OnUpdate:     "fn::CURRENT_TIMESTAMP",
				},
			},
			Description: "",
			Group:       "common",
		},
		{
			Name: "group",
			Columns: []*octopus.Column{
				{
					Name:            "id",
					Type:            octopus.ColTypeInt64,
					PrimaryKey:      true,
					AutoIncremental: true,
				},
				{
					Name:      "name",
					Type:      octopus.ColTypeVarchar,
					Size:      100,
					UniqueKey: true,
				},
				{
					Name:         "created_at",
					Type:         octopus.ColTypeDateTime,
					DefaultValue: "fn::CURRENT_TIMESTAMP",
				},
				{
					Name:         "updated_at",
					Type:         octopus.ColTypeDateTime,
					DefaultValue: "fn::CURRENT_TIMESTAMP",
					OnUpdate:     "fn::CURRENT_TIMESTAMP",
				},
			},
			Description: "",
			Group:       "common",
			Indices: []*octopus.Index{
				{
					Name:    "name_idx",
					Columns: []string{"name"},
				},
			},
		},
	},
}

func TestExportImport(t *testing.T) {
	Convey("export and import", t, func() {
		xlsFile, err := ioutil.TempFile("", "*.xlsx")
		filename := xlsFile.Name()
		So(err, ShouldBeNil)

		fmt.Println(filename)
		defer os.Remove(filename)

		// export xlsx
		exp := Exporter{
			schema: testSchema,
			option: &ExportOption{
				UseNullColumn: false,
			},
		}
		err = exp.Export(filename)
		So(err, ShouldBeNil)

		// import xlsx
		imp := Importer{}
		schema, err := imp.Import(filename)
		So(err, ShouldBeNil)

		// test schema
		if diff := cmp.Diff(testSchema, schema); diff != "" {
			So(diff, ShouldBeEmpty)
		}
	})
}
