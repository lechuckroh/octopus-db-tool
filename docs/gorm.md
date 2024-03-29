# GORM

[한국어](kr/gorm.md)

## Generate

```shell
$ oct generate gorm --help
```

|            Option            |         Env. Variable         | Description                                                                                                                       |
| :--------------------------: | :---------------------------: | :-------------------------------------------------------------------------------------------------------------------------------- |
|       `-i`, `--input`        |        `OCTOPUS_INPUT`        | Octopus schema file to read                                                                                                       |
|       `-o`, `--output`       |       `OCTOPUS_OUTPUT`        | Target file or directory                                                                                                          |
| `-a`, `--pointerAssociation` | `OCTOPUS_POINTER_ASSOCIATION` | Use pointer type on associated field if flag is set                                                                               |
|       `-e`, `--embed`        |        `OCTOPUS_EMBED`        | Define embedded struct.<br />Format: `<structName>:<columnName1>[,<columnName2>]...`                                              |
|       `-g`, `--groups`       |       `OCTOPUS_GROUPS`        | Table groups to generate.<br />Set multiple groups with comma(`,`) separated.                                                     |
|      `-k`, `--package`       |       `OCTOPUS_PACKAGE`       | Source package name                                                                                                               |
|       `-p`, `--prefix`       |       `OCTOPUS_PREFIX`        | Model struct name prefix.<br />Format: `<group1>:<prefix1>[,<group2>:<prefix2>]...`<br />Example: `group1:prefix1,group2:prefix2` |
|    `-r`, `--removePrefix`    |    `OCTOPUS_REMOVE_PREFIX`    | Prefixes to remove from model struct name.<br />Set multiple prefixes with comma(`,`) separated.                                  |
|  `-u`, `--uniqueNameSuffix`  | `OCTOPUS_UNIQUE_NAME_SUFFIX`  | Unique constraint name suffix                                                                                                     |

### `--embed` option

- Matches column name only.
- Default embedded struct is:
  - `gorm.Model`: `id`, `created_at`, `updated_at`, `deleted_at`
  - See [gorm.Model](https://gorm.io/docs/models.html#gorm-Model)
  - To disable default `gorm.Model`, use `--embed gorm.Model:` option.

## Example

### Generate

```shell
$ oct generate gorm \
    --input examples/user.json \
    --output output/user.go \
    --package model
```

Generated `*.go` file:

```go
package model

import (
	"gopkg.in/guregu/null.v4"
)

type UserGroup struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(40);unique;not null"`
}

func (c *UserGroup) TableName() string { return "group" }

type User struct {
	ID        int64  `gorm:"primary_key;auto_increment"`
	Name      string `gorm:"type:varchar(40);unique;not null"`
	GroupID   null.Int
	UserGroup UserGroup `gorm:"foreignKey:GroupID;references:ID"`
}

func (c *User) TableName() string { return "user" }
```

### Generate with custom embedded struct

Assume that you have defined a custom embedded struct `IdName`:

```go
type IdName struct {
    ID   int64  `gorm:"primary_key;auto_increment"`
    Name string `gorm:"type:varchar(40);unique;not null"`
}
```

```shell
$ oct generate gorm \
    --input examples/user.json \
    --output output/user.go \
    --package model \
    --embed IdName:id,name
```

Generated `*.go` file:

```go
package model

import (
	"gopkg.in/guregu/null.v4"
)

type UserGroup struct {
	IdName
}

func (c *UserGroup) TableName() string { return "group" }

type User struct {
	IdName
	GroupID   null.Int
	UserGroup UserGroup `gorm:"foreignKey:GroupID;references:ID"`
}

func (c *User) TableName() string { return "user" }
```
