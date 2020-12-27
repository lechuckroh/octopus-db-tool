package jpa

import (
	"bytes"
	"github.com/lechuckroh/octopus-db-tools/format/common"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var jpaKotlinTplTestSchema = &octopus.Schema{
	Tables: []*octopus.Table{
		{
			Name: "user",
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
					Name:    "dec",
					Type:    octopus.ColTypeDecimal,
					Size:    20,
					Scale:   5,
					NotNull: true,
				},
				{
					Name:    "created_at",
					Type:    octopus.ColTypeDateTime,
					NotNull: true,
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
func TestJPAKotlin_GenerateEntityClass(t *testing.T) {
	Convey("GenerateEntityClass", t, func() {
		option := &KtOption{
			PrefixMapper:     common.NewPrefixMapper("common:C"),
			AnnoMapper:       common.NewAnnotationMapper("common:@Common,admin:@Admin"),
			IdEntity:         "IdEntity",
			Package:          "com.lechuck",
			UniqueNameSuffix: "_uq",
		}
		expected := []string{
			`package com.lechuck

import org.hibernate.annotations.CreationTimestamp
import org.hibernate.annotations.UpdateTimestamp

import java.math.BigDecimal
import java.sql.Timestamp
import javax.persistence.*

@Common
@Entity
@Table(name="user", uniqueConstraints = [
    UniqueConstraint(name = "user_uq", columnNames = ["name"])
])
data class CUser(
        @Id
        @GeneratedValue(strategy = GenerationType.AUTO)
        @Column(nullable = false)
        override var id: Long = 0L,

        @Column(nullable = false, length = 100)
        var name: String = "",

        @Column(nullable = false, precision = 20, scale = 5)
        var dec: BigDecimal = BigDecimal.ZERO,

        @CreationTimestamp
        @Column(nullable = false, updatable = false)
        var createdAt: Timestamp = Timestamp(System.currentTimeMillis()),

        @UpdateTimestamp
        var updatedAt: Timestamp?

): IdEntity<Long>
`,
		}

		jpaKotlin := NewKtGenerator(jpaKotlinTplTestSchema, option)

		for i, table := range jpaKotlinTplTestSchema.Tables {
			class := NewKotlinClass(table, option)
			buf := new(bytes.Buffer)
			if err := jpaKotlin.GenerateEntityClass(buf, class); err != nil {
				t.Error(err)
			}
			actual := buf.String()
			So(actual, ShouldResemble, expected[i])
		}
	})
}
