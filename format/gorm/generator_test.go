package gorm

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"strings"
	"testing"
)

var table1 = &octopus.Table{
	Name: "tbl_user1",
	Columns: []*octopus.Column{
		{
			Name:            "id",
			Type:            octopus.ColTypeInt64,
			PrimaryKey:      true,
			AutoIncremental: true,
			NotNull:         true,
		},
		{
			Name:      "name",
			Type:      octopus.ColTypeVarchar,
			Size:      100,
			UniqueKey: true,
			NotNull:   true,
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
			Name:    "float1_notnull",
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
			Name:    "double1_notnull",
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
			Name: "time_null",
			Type: octopus.ColTypeTime,
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
			Name:    "bit1_notnull",
			Type:    octopus.ColTypeBit,
			Size:    1,
			NotNull: true,
		},
		{
			Name: "bit1_null",
			Type: octopus.ColTypeBit,
			Size: 1,
		},
		{
			Name: "bit2",
			Type: octopus.ColTypeBit,
			Size: 2,
		},
		{
			Name: "invalid",
			Type: "invalidType",
		},
		{
			Name:    "created_at",
			Type:    octopus.ColTypeDateTime,
			NotNull: true,
		},
		{
			Name:    "updated_at",
			Type:    octopus.ColTypeDateTime,
			NotNull: true,
		},
		{
			Name: "deleted_at",
			Type: octopus.ColTypeDateTime,
		},
	},
	Description: "",
	Group:       "common",
}

var table2 = &octopus.Table{
	Name: "user2",
	Columns: []*octopus.Column{
		{
			Name:            "id",
			Type:            octopus.ColTypeInt64,
			PrimaryKey:      true,
			AutoIncremental: true,
			NotNull:         true,
		},
		{
			Name:      "name",
			Type:      octopus.ColTypeVarchar,
			Size:      100,
			UniqueKey: true,
			NotNull:   true,
		},
		{
			Name:      "passport_no",
			Type:      octopus.ColTypeVarchar,
			Size:      20,
			UniqueKey: true,
			NotNull:   true,
		},
		{
			Name: "ch",
			Type: octopus.ColTypeChar,
			Size: 10,
		},
		{
			Name:    "dec",
			Type:    octopus.ColTypeDecimal,
			Size:    20,
			Scale:   5,
			NotNull: true,
		},
		{
			Name:    "time_notnull",
			Type:    octopus.ColTypeTime,
			NotNull: true,
		},
		{
			Name:    "created_at",
			Type:    octopus.ColTypeDateTime,
			NotNull: true,
		},
		{
			Name:    "updated_at",
			Type:    octopus.ColTypeDateTime,
			NotNull: true,
		},
	},
	Description: "",
	Group:       "common",
	Indices: []*octopus.Index{
		{
			Name:    "dec_idx",
			Columns: []string{"dec"},
		},
		{
			Name:    "idx2",
			Columns: []string{"dec", "name"},
		},
	},
}

var testSchema = &octopus.Schema{
	Tables: []*octopus.Table{table1, table2},
}

func TestGorm_Generate(t *testing.T) {
	Convey("Generate", t, func() {
		packages := []string{"lechuck", ""}

		for _, pkg := range packages {
			option := &Option{
				PrefixMapper:     common.NewPrefixMapper("common:C"),
				RemovePrefixes:   []string{"tbl_"},
				Package:          pkg,
				UniqueNameSuffix: "_uq",
			}

			expectedStrings := []string{
				"package " + util.IfThenElseString(pkg != "", pkg, "main"),
				"",
				"import (",
				"	\"github.com/jinzhu/gorm\"",
				"	\"github.com/shopspring/decimal\"",
				"	\"gopkg.in/guregu/null.v4\"",
				"	\"time\"",
				")",
				"",
				"type CUser1 struct {",
				"	gorm.Model",
				"	Name string `gorm:\"type:varchar(100);unique;not null\"`",
				"	_1Decimal decimal.Decimal `gorm:\"column:1_decimal;type:decimal(20,5);not null\"`",
				"	_2Decimal decimal.NullDecimal `gorm:\"column:2_decimal;type:decimal(20,5)\"`",
				"	Bool1 null.Bool",
				"	Bool2 bool `gorm:\"column:bool2;not null\"`",
				"	Int1Notnull int8 `gorm:\"column:int1_notnull;not null\"`",
				"	Int1Null null.Int `gorm:\"column:int1_null\"`",
				"	Int2Notnull int16 `gorm:\"column:int2_notnull;not null\"`",
				"	Int2Null null.Int `gorm:\"column:int2_null\"`",
				"	Int3Notnull int32 `gorm:\"column:int3_notnull;not null\"`",
				"	Int3Null null.Int `gorm:\"column:int3_null\"`",
				"	Int4Notnull int32 `gorm:\"column:int4_notnull;not null\"`",
				"	Int4Null null.Int `gorm:\"column:int4_null\"`",
				"	Int5Notnull int64 `gorm:\"column:int5_notnull;not null\"`",
				"	Int5Null null.Int `gorm:\"column:int5_null\"`",
				"	Float1Notnull float32 `gorm:\"column:float1_notnull;type:float(20,5);not null\"`",
				"	FloatNull null.Float `gorm:\"type:float(20,10)\"`",
				"	Double1Notnull float64 `gorm:\"column:double1_notnull;type:double(20,5);not null\"`",
				"	DoubleNull null.Float `gorm:\"type:double(20,10)\"`",
				"	TimeNull null.Time",
				"	Blob1 []byte `gorm:\"column:blob1\"`",
				"	Blob2 []byte `gorm:\"column:blob2\"`",
				"	Blob3 []byte `gorm:\"column:blob3\"`",
				"	Blob4 []byte `gorm:\"column:blob4\"`",
				"	Bit1Notnull byte `gorm:\"column:bit1_notnull;type:bit(1);not null\"`",
				"	Bit1Null *byte `gorm:\"column:bit1_null;type:bit(1)\"`",
				"	Bit2 *byte `gorm:\"column:bit2;type:bit(2)\"`",
				"	Invalid interface{}",
				"}",
				"",
				"func (c *CUser1) TableName() string { return \"tbl_user1\" }",
				"",
				"type CUser2 struct {",
				"	ID int64 `gorm:\"primary_key;auto_increment\"`",
				"	Name string `gorm:\"type:varchar(100);unique_index:user2_uq;index:idx2,priority:2;not null\"`",
				"	PassportNo string `gorm:\"type:varchar(20);unique_index:user2_uq;not null\"`",
				"	Ch null.String `gorm:\"type:char(10)\"`",
				"	Dec decimal.Decimal `gorm:\"type:decimal(20,5);index:dec_idx;index:idx2,priority:1;not null\"`",
				"	TimeNotnull time.Time `gorm:\"not null\"`",
				"	CreatedAt time.Time `gorm:\"not null\"`",
				"	UpdatedAt time.Time `gorm:\"not null\"`",
				"}",
				"",
				"func (c *CUser2) TableName() string { return \"user2\" }",
				"",
				"",
			}
			expected := strings.Join(expectedStrings, "\n")

			gen := Generator{schema: testSchema, option: option}

			buf := new(bytes.Buffer)
			if err := gen.Generate(buf); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			if diff := cmp.Diff(expected, actual); diff != "" {
				log.Println(diff)
			}
			So(actual, ShouldResemble, expected)
		}
	})
}
