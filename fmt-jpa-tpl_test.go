package main

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var jpaKotlinTplTestSchema = &Schema{
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
func TestJPAKotlinTpl_GenerateEntityClass_1(t *testing.T) {
	output := &Output{
		Options: map[string]string{
			FlagIdEntity:         "IdEntity",
			FlagPackage:          "com.lechuck",
			FlagUniqueNameSuffix: "_uq",
			FlagAnnotation:       "common:@Common,admin:@Admin",
		},
	}
	annoMapper := newAnnotationMapper(output.Options[FlagAnnotation])
	prefixMapper := newPrefixMapper("common:C")
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

	jpaKotlin := NewJPAKotlinTpl(jpaKotlinTplTestSchema, output, nil, annoMapper, prefixMapper, true)

	for i, table := range jpaKotlinTplTestSchema.Tables {
		class := NewKotlinClass(table, output, annoMapper, prefixMapper)
		buf := new(bytes.Buffer)
		if err := jpaKotlin.GenerateEntityClass(buf, class); err != nil {
			t.Error(err)
		}
		actual := buf.String()
		if diff := cmp.Diff(expected[i], actual); diff != "" {
			t.Errorf("mismatch [%d] (-expected +actual):\n%s", i, diff)
		}
	}
}

// not data class
func TestJPAKotlinTpl_GenerateEntityClass_2(t *testing.T) {
	output := &Output{
		Options: map[string]string{
			FlagPackage:          "com.lechuck",
			FlagUniqueNameSuffix: "_uq",
		},
	}
	annoMapper := newAnnotationMapper(output.Options[FlagAnnotation])
	prefixMapper := newPrefixMapper("common:C")
	expected := []string{
		`package com.lechuck

import org.hibernate.annotations.CreationTimestamp
import org.hibernate.annotations.UpdateTimestamp

import java.math.BigDecimal
import java.sql.Timestamp
import javax.persistence.*


@Entity
@Table(name="user", uniqueConstraints = [
    UniqueConstraint(name = "user_uq", columnNames = ["name"])
])
class CUser(
        @Id
        @GeneratedValue(strategy = GenerationType.AUTO)
        @Column(nullable = false)
        var id: Long = 0L,

        @Column(nullable = false, length = 100)
        var name: String = "",

        @Column(nullable = false, precision = 20, scale = 5)
        var dec: BigDecimal = BigDecimal.ZERO,

        @CreationTimestamp
        @Column(nullable = false, updatable = false)
        var createdAt: Timestamp = Timestamp(System.currentTimeMillis()),

        @UpdateTimestamp
        var updatedAt: Timestamp?

): AbstractJpaPersistable<Long>()
`,
	}

	jpaKotlin := NewJPAKotlinTpl(jpaKotlinTplTestSchema, output, nil, annoMapper, prefixMapper, false)

	for i, table := range jpaKotlinTplTestSchema.Tables {
		class := NewKotlinClass(table, output, annoMapper, prefixMapper)
		buf := new(bytes.Buffer)
		if err := jpaKotlin.GenerateEntityClass(buf, class); err != nil {
			t.Error(err)
		}
		actual := buf.String()
		if diff := cmp.Diff(expected[i], actual); diff != "" {
			t.Errorf("mismatch [%d] (-expected +actual):\n%s", i, diff)
		}
	}
}
