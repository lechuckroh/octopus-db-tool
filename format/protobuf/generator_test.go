package protobuf

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var protobufTplTestSchema = &octopus.Schema{
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
					Name:      "name",
					Type:      octopus.ColTypeVarchar,
					Size:      100,
					UniqueKey: true,
				},
				{
					Name:  "dec",
					Type:  octopus.ColTypeDecimal,
					Size:  20,
					Scale: 5,
				},
				{
					Name: "created_at",
					Type: octopus.ColTypeDateTime,
				},
				{
					Name: "updated_at",
					Type: octopus.ColTypeDateTime,
				},
			},
			Description: "",
			Group:       "common",
		},
	},
}

// data class
func TestProtobufTpl_Generate(t *testing.T) {
	Convey("Generate", t, func() {
		option := &Option{
			Package:      "com.lechuck.foo",
			GoPackage:    "lechuck/foo",
			PrefixMapper: common.NewPrefixMapper("common:C"),
		}
		expected := []string{
			`syntax = "proto3";

package com.lechuck.hello;

option go_package = "proto/hello";

import "google/protobuf/timestamp.proto";

message CUser {
  int64 id = 1;
  string name = 2;
  double dec = 3;
  google.protobuf.Timestamp createdAt = 4;
  google.protobuf.Timestamp updatedAt = 5;
}
`,
		}

		protobuf := NewGenerator(protobufTplTestSchema, option)

		for i, table := range protobufTplTestSchema.Tables {
			messages := []*PbMessage{
				NewPbMessage(table, option),
			}
			buf := new(bytes.Buffer)
			if err := protobuf.generateProto(buf, messages, "com.lechuck.hello", "proto/hello"); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			So(actual, ShouldResemble, expected[i])
		}
	})
}
