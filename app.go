package main

import (
	"errors"
	"io/ioutil"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Create(filename string) error {
	schema := &Schema{
		Version: "0.0.1",
		Tables: []*Table{
			{
				Name:        "user",
				Description: "User table",
				Tags:        []string{"tag"},
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
							Table:    "group",
							Column:   "id",
							Nullable: true,
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
	bytes, err := schema.toJson()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func (app *App) Convert(source, sourceFormat, target, targetFormat string) error {
	// TODO
	return errors.New("not implemented")
}

func (app *App) Generate(source, target, format string) error {
	// TODO
	return errors.New("not implemented")
}
