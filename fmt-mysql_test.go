package main

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestMysql_ToSchema(t *testing.T) {
	sql := strings.Join([]string{
		"CREATE TABLE IF NOT EXISTS `Table` (",
		"id bigint NOT NULL AUTO_INCREMENT,",
		"name varchar(20) NOT NULL DEFAULT 'noname',",
		"created_at datetime,",
		"PRIMARY KEY (`id`),",
		"UNIQUE KEY `Table_UNIQUE` (`name`)",
		");",
	}, "\n")

	mysql := Mysql{}

	// read sql
	if err := mysql.FromString([]byte(sql)); err != nil {
		t.Error(err)
	}

	// convert to schema
	schema, err := mysql.ToSchema()
	if err != nil {
		t.Error(err)
	}

	expected := &Schema{
		Tables: []*Table{
			{
				Name: "Table",
				Columns: []*Column{
					{
						Name:            "id",
						Type:            ColTypeLong,
						PrimaryKey:      true,
						Nullable:        false,
						AutoIncremental: true,
					},
					{
						Name:         "name",
						Type:         ColTypeString,
						Size:         20,
						Nullable:     false,
						UniqueKey:    true,
						DefaultValue: "noname",
					},
					{
						Name:     "created_at",
						Type:     ColTypeDateTime,
						Nullable: true,
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, schema); diff != "" {
		t.Errorf("TestMysql_ToSchema() mismatch (-expected +actual):\n%s", diff)
	}
}

func TestMysql_ToString(t *testing.T) {
	schema := Schema{
		Author:  "Author",
		Name:    "Name",
		Version: "Version",
		Tables: []*Table{
			{
				Name: "Table",
				Columns: []*Column{
					{
						Name:            "id",
						Type:            ColTypeLong,
						PrimaryKey:      true,
						Nullable:        false,
						AutoIncremental: true,
					},
					{
						Name:         "name",
						Type:         ColTypeString,
						Size:         20,
						Nullable:     false,
						UniqueKey:    true,
						DefaultValue: "noname",
					},
				},
				Description: "description",
				Group:       "group",
				ClassName:   "className",
			},
		},
	}
	expected := []string{
		"CREATE TABLE IF NOT EXISTS `Table` (",
		"`id` bigint NOT NULL AUTO_INCREMENT,",
		"`name` varchar(20) NOT NULL DEFAULT 'noname',",
		"PRIMARY KEY (`id`),",
		"UNIQUE KEY `Table_UNIQUE` (`name`)",
		");",
	}

	// convert to string
	mysql := Mysql{}
	result, err := mysql.ToString(&schema)
	if err != nil {
		t.Error(err)
	}

	actual := make([]string, 0)
	for _, line := range strings.Split(string(result), "\n") {
		actual = append(actual, strings.TrimSpace(line))
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("TestMysql_ToString() mismatch (-expected +actual):\n%s", diff)
	}
}
