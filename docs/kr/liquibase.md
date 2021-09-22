# Liquibase

[English](../liquibase.md)

## ChangeLog 생성

```shell
$ oct generate liquibase --help
```

|            옵션            |           환경변수           | 설명                                                              |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                             |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | 생성할 파일명                                                     |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분 |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | 유니크 제약 이름 접미사                                           |
|     `-c`, `--comments`     |      `OCTOPUS_COMMENTS`      | 테이블/컬럼 설명을 같이 생성할지 여부. 기본값: `false`            |

### 예제

```shell
$ oct generate liquibase \
    --input examples/user.json \
    --output output/changelogs.yml \
    --comments
```

`*.yaml` 파일은 다음과 같이 생성됩니다:

```yaml
databaseChangeLog:
  - objectQuotingStrategy: QUOTE_ALL_OBJECTS
  - changeSet:
      id: 1-1
      changes:
        - createTable:
            tableName: group
            remarks: Group table
            columns:
              - column:
                  name: id
                  type: bigint
                  autoIncrement: true
                  constraints:
                    primaryKey: true
                  remarks: unique id
              - column:
                  name: name
                  type: varchar(40)
                  constraints:
                    nullable: false
                    unique: true
                  remarks: group name
  - changeSet:
      id: 1-2
      preConditions:
        dbms:
          type: derby, h2, mssql, mariadb, mysql, postgresql, sqlite
        onError: CONTINUE
        onFail: CONTINUE
      changes:
        - addUniqueConstraint:
            tableName: group
            columnNames: name
            constraintName: group
  - changeSet:
      id: 2-1
      changes:
        - createTable:
            tableName: user
            remarks: User table
            columns:
              - column:
                  name: id
                  type: bigint
                  autoIncrement: true
                  constraints:
                    primaryKey: true
                  remarks: unique id
              - column:
                  name: name
                  type: varchar(40)
                  constraints:
                    nullable: false
                    unique: true
                  remarks: user login name
              - column:
                  name: group_id
                  type: bigint
                  remarks: group ID
  - changeSet:
      id: 2-2
      preConditions:
        dbms:
          type: derby, h2, mssql, mariadb, mysql, postgresql, sqlite
        onError: CONTINUE
        onFail: CONTINUE
      changes:
        - addUniqueConstraint:
            tableName: user
            columnNames: name
            constraintName: user
```

## Diff ChangeLog 생성

2개의 스키마를 비교해서 변경된 diff changelog 파일을 생성합니다.

```shell
$ oct diff liquibase --help
```

|            옵션            |           환경변수           | 설명                                                              |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------- |
|      `-a`, `--author`      |       `OCTOPUS_AUTHOR`       | Diff 작성자                                                       |
|       `-f`, `--from`       |        `OCTOPUS_FROM`        | 비교할 변경 전 스키마 파일명                                      |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분 |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | 생성할 diff changelog 파일명                                      |
|        `-t`, `--to`        |         `OCTOPUS_TO`         | 비교할 변경 후 스키마 파일명                                      |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | 유니크 제약 이름 접미사                                           |
|     `-c`, `--comments`     |      `OCTOPUS_COMMENTS`      | 테이블/컬럼 설명을 같이 생성할지 여부. 기본값: `false`            |

### 예제

```shell
$ oct diff liquibase \
    --from examples/user.json \
    --to examples/user-v2.json \
    --output output/diff.yaml \
    --uniqueNameSuffix _uq \
    --author foo \
    --comments
```

`*.yaml` 파일은 다음과 같이 생성됩니다:

```yaml
databaseChangeLog:
  - objectQuotingStrategy: QUOTE_ALL_OBJECTS
  - changeSet:
      id: 1-1
      author: foo
      changes:
        - modifyDataType:
            tableName: group
            columnName: name
            newDataType: varchar(80)
  - changeSet:
      id: 2-1
      author: foo
      changes:
        - addColumn:
            tableName: user
            columns:
              - column:
                  name: email
                  type: varchar(255)
                  remarks: user email
                  afterColumn: name
```
