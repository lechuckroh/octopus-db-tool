package gorm

import (
	"github.com/google/go-cmp/cmp"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestParseEmbeddedModel(t *testing.T) {
	Convey("Test", t, func() {
		testData := []struct {
			Definition          string
			ExpectedName        string
			ExpectedColumnNames []string
		}{
			{
				Definition:          "gorm.Model:id,created_at,updated_at,deleted_at",
				ExpectedName:        "gorm.Model",
				ExpectedColumnNames: []string{"id", "created_at", "updated_at", "deleted_at"},
			},
			{
				Definition:          "MyName0_1:",
				ExpectedName:        "MyName0_1",
				ExpectedColumnNames: []string{},
			},
			{
				Definition:          "MyName0_1=1234",
				ExpectedName:        "",
				ExpectedColumnNames: nil,
			},
		}

		for _, test := range testData {
			name, columnNames := parseEmbeddedModelDefinition(test.Definition)
			if diff := cmp.Diff(test.ExpectedName, name); diff != "" {
				log.Println(diff)
			}
			if diff := cmp.Diff(test.ExpectedColumnNames, columnNames); diff != "" {
				log.Println(diff)
			}
			So(name, ShouldEqual, test.ExpectedName)
			So(columnNames, ShouldResemble, test.ExpectedColumnNames)
		}
	})
}
