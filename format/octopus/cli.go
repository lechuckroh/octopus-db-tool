package octopus

import "github.com/urfave/cli/v2"

const (
	FlagOutput = "output"
)

func InitAction(c *cli.Context) error {
	schema := Schema{
		Version: "1.0.0",
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
							Table:        "group",
							Column:       "id",
							Relationship: RefManyToOne,
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
	return schema.ToFile(c.String(FlagOutput))
}

var InitCliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        FlagOutput,
		Aliases:     []string{"o"},
		Usage:       "init octopus schema to `FILE`",
		EnvVars:     []string{"OCTOPUS_OUTPUT"},
		DefaultText: "db.json",
	},
}
