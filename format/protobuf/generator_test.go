package protobuf

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"log"
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
					Name: "created_at",
					Type: octopus.ColTypeDateTime,
				},
				{
					Name: "updated_at",
					Type: octopus.ColTypeDateTime,
				},
			},
			Description: "",
			Group:       "group",
		},
	},
}

// data class
func TestProtobufTpl_Generate(t *testing.T) {
	Convey("Generate", t, func() {
		option := &Option{
			Package:          "com.lechuck.foo",
			GoPackage:        "lechuck/foo",
			PrefixMapper:     common.NewPrefixMapper("common:C"),
			RelationTagStart: 100,
		}
		expected := `syntax = "proto3";

package com.lechuck.foo;

option go_package = "lechuck/foo";

import "google/protobuf/timestamp.proto";

message CUser {
  int64 id = 1;
  string name = 2;
  int64 groupId = 3;
  double 1Decimal = 4;
  double 2Decimal = 5;
  bool bool1 = 6;
  bool bool2 = 7;
  int32 int1Notnull = 8;
  int32 int1Null = 9;
  int32 int2Notnull = 10;
  int32 int2Null = 11;
  int32 int3Notnull = 12;
  int32 int3Null = 13;
  int32 int4Notnull = 14;
  int32 int4Null = 15;
  int64 int5Notnull = 16;
  int64 int5Null = 17;
  float floatNotnull = 18;
  float floatNull = 19;
  double doubleNotnull = 20;
  double doubleNull = 21;
  google.protobuf.Timestamp createdAt = 22;
  google.protobuf.Timestamp updatedAt = 23;
  Group group = 100;
}

message Group {
  int64 id = 1;
  string name = 2;
  google.protobuf.Timestamp createdAt = 3;
  google.protobuf.Timestamp updatedAt = 4;
}
`

		buf := new(bytes.Buffer)
		gen := newGenerator(protobufTplTestSchema, option)
		if err := gen.Generate(buf); err != nil {
			t.Error(err)
		}
		actual := buf.String()
		if diff := cmp.Diff(expected, actual); diff != "" {
			log.Println(diff)
		}
		So(actual, ShouldEqual, expected)
	})
}

func TestProtobufTpl_NoPackage(t *testing.T) {
	Convey("No package", t, func() {
		option := &Option{
			PrefixMapper:     common.NewPrefixMapper(""),
			RelationTagStart: 100,
			TableFilter: func(table *octopus.Table) bool {
				return table.Group == "group"
			},
		}
		expected := `syntax = "proto3";



import "google/protobuf/timestamp.proto";

message Group {
  int64 id = 1;
  string name = 2;
  google.protobuf.Timestamp createdAt = 3;
  google.protobuf.Timestamp updatedAt = 4;
}
`

		buf := new(bytes.Buffer)
		gen := newGenerator(protobufTplTestSchema, option)
		if err := gen.Generate(buf); err != nil {
			t.Error(err)
		}
		actual := buf.String()
		if diff := cmp.Diff(expected, actual); diff != "" {
			log.Println(diff)
		}
		So(actual, ShouldEqual, expected)
	})
}
