package gorm

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
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
					Name:     "updated_at",
					Type:     "datetime",
					Nullable: true,
				},
			},
			Description: "",
			Group:       "common",
		},
	},
}

// data class
func TestGorm_GenerateStruct(t *testing.T) {
	Convey("GenerateStruct", t, func() {
		option := &Option{
			PrefixMapper:     common.NewPrefixMapper("common:C"),
			Package:          "lechuck",
			UniqueNameSuffix: "_uq",
		}
		expectedStrings := []string{
			"type CUser struct {",
			"	gorm.Model",
			"	Name string `gorm:\"type:varchar(100);unique;not null\"`",
			"	Dec decimal.Decimal `gorm:\"type:decimal(20,5);not null\"`",
			"}",
			"",
			"func (c *CUser) TableName() string { return \"user\" }",
			"",
			"",
		}
		expected := strings.Join(expectedStrings, "\n")

		gorm := Generator{schema: testSchema, option: option}

		for _, table := range testSchema.Tables {
			gormStruct := NewGoStruct(table, option)
			buf := new(bytes.Buffer)
			if err := gorm.GenerateStruct(buf, gormStruct); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			So(actual, ShouldResemble, expected)
		}
	})
}
