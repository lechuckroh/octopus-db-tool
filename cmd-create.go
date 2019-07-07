package main

import "io/ioutil"

type CreateCmd struct {
}

func (cmd *CreateCmd) createDefaultSchema() *Schema {
	return &Schema{
		Version: "0.0.1",
		Tables: []*Table{
			{
				Name:        "user",
				Description: "User table",
				Columns: []*Column{
					{
						Name:            "id",
						Type:            "bigint",
						Description:     "unique id",
						PrimaryKey:      true,
						AutoIncremental: true,
					},
					{
						Name:        "name",
						Type:        "varchar",
						Size:        40,
						Description: "user login name",
						UniqueKey:   true,
					},
					{
						Name:        "group_id",
						Type:        "bigint",
						Description: "group ID",
						Ref: &Reference{
							Table:  "group",
							Column: "id",
						},
					},
				},
			},
			{
				Name:        "group",
				Description: "Group table",
				Columns: []*Column{
					{
						Name:            "id",
						Type:            "bigint",
						Description:     "unique id",
						PrimaryKey:      true,
						AutoIncremental: true,
					},
					{
						Name:        "name",
						Type:        "varchar",
						Size:        40,
						Description: "group name",
						UniqueKey:   true,
					},
				},
			},
		},
	}
}

func (cmd *CreateCmd) Create(target *Output) error {
	schema := cmd.createDefaultSchema()
	bytes, err := schema.ToJson()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(target.Filename, bytes, 0644)
}
