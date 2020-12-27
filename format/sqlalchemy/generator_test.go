package sqlalchemy

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var testSchema = &octopus.Schema{
	Tables: []*octopus.Table{
		{
			Name: "user",
			Columns: []*octopus.Column{
				{
					Name:            "id",
					Type:            "long",
					PrimaryKey:      true,
					AutoIncremental: true,
				},
				{
					Name:      "name",
					Type:      "string",
					Size:      100,
					UniqueKey: true,
				},
				{
					Name:  "dec",
					Type:  "decimal",
					Size:  20,
					Scale: 5,
				},
				{
					Name: "created_at",
					Type: "datetime",
				},
				{
					Name: "updated_at",
					Type: "datetime",
				},
			},
			Description: "",
			Group:       "common",
		},
	},
}

func TestGenerator_GenerateClass(t *testing.T) {
	// TODO: remove skip
	SkipConvey("GenerateClass", t, func() {
		option := &Option{
			PrefixMapper:     common.NewPrefixMapper("common:C"),
			TableFilter:      nil,
			RemovePrefixes:   []string{},
			UniqueNameSuffix: "_uq",
			UseUTC:           true,
		}
		// FIXME: update expected
		expected := []string{
			``,
		}

		g := Generator{
			schema: testSchema,
			option: option,
		}

		for i, table := range testSchema.Tables {
			class := NewSaClass(table, option)
			buf := new(bytes.Buffer)
			if err := g.GenerateClass(buf, class); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			So(actual, ShouldResemble, expected[i])
		}
	})
}
