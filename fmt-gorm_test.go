package main

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var gormTplTestSchema = &Schema{
	Tables: []*Table{
		{
			Name: "user",
			Columns: []*Column{
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
		output := &Output{
			Options: map[string]string{
				FlagPackage:          "lechuck",
				FlagUniqueNameSuffix: "_uq",
			},
		}
		prefixMapper := newPrefixMapper("common:C")
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

		gorm := NewGormTpl(gormTplTestSchema, output, nil, prefixMapper)

		for _, table := range gormTplTestSchema.Tables {
			gormStruct := NewGormStruct(table, output, prefixMapper)
			buf := new(bytes.Buffer)
			if err := gorm.GenerateStruct(buf, gormStruct); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			So(actual, ShouldResemble, expected)
		}
	})
}
