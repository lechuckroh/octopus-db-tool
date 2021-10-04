# GORM

[English](../gorm.md)

## 소스 생성

```shell
$ oct generate gorm --help
```

|             옵션             |           환경변수            | 설명                                                                                                                               |
| :--------------------------: | :---------------------------: | :--------------------------------------------------------------------------------------------------------------------------------- |
|       `-i`, `--input`        |        `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                                                                                              |
|       `-o`, `--output`       |       `OCTOPUS_OUTPUT`        | 출력할 파일명 또는 디렉토리명                                                                                                      |
| `-a`, `--pointerAssociation` | `OCTOPUS_POINTER_ASSOCIATION` | 플래그 설정시 연관관계로 설정된 필드 타입에 포인터 타입 사용                                                                       |
|       `-e`, `--embed`        |        `OCTOPUS_EMBED`        | 사용할 embedded struct 정의.<br />형식: `<structName>:<컬럼1>[,<컬럼2>]...`                                                        |
|       `-g`, `--groups`       |       `OCTOPUS_GROUPS`        | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분                                                                  |
|      `-k`, `--package`       |       `OCTOPUS_PACKAGE`       | 생성할 소스 파일의 패키지명                                                                                                        |
|       `-p`, `--prefix`       |       `OCTOPUS_PREFIX`        | 생성할 모델 struct 이름의 접두사.<br />형식: `<그룹1>:<접두사1>[,<그룹2>:<접두사2>]...`<br />예제: `group1:prefix1,group2:prefix2` |
|    `-r`, `--removePrefix`    |    `OCTOPUS_REMOVE_PREFIX`    | 모델 struct 이름에서 제거할 접두사.<br />여러개의 접두사를 지정시 `,`로 구분                                                       |
|  `-u`, `--uniqueNameSuffix`  | `OCTOPUS_UNIQUE_NAME_SUFFIX`  | 유니크 제약 이름에 사용할 접미사                                                                                                   |

### `--embed` 옵션

- embedded struct 사용 여부는 컬럼 이름만 사용해서 판단합니다.
- 기본적으로 사용되는 embedded struct:
  - `gorm.Model`: `id`, `created_at`, `updated_at`, `deleted_at`
  - [gorm.Model](https://gorm.io/docs/models.html#gorm-Model) 문서 참고
  - 기본 `gorm.Model`을 사용하지 않으려면, `--embed gorm.Model:` 옵션을 사용합니다.

## 예제

### 소스 생성

```shell
$ oct generate gorm \
    --input examples/user.json \
    --output output/user.go \
    --package model
```

`*.go` 파일은 다음과 같이 생성됩니다:

```go
package model

import (
	"gopkg.in/guregu/null.v4"
)

type UserGroup struct {
	ID int64 `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(40);unique;not null"`
}

func (c *UserGroup) TableName() string { return "group" }

type User struct {
	ID int64 `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(40);unique;not null"`
	GroupID null.Int
	UserGroup UserGroup `gorm:"foreignKey:GroupID;references:ID"`
}

func (c *User) TableName() string { return "user" }
```

### 커스텀 embedded struct 를 사용해 소스 생성

다음과 같이 `IdName` 이라는 커스텀 embedded struct를 정의했다고 가정합니다:

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

`*.go` 파일은 다음과 같이 생성됩니다:

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
  GroupID null.Int
  UserGroup UserGroup `gorm:"foreignKey:GroupID;references:ID"`
}

func (c *User) TableName() string { return "user" }
```
