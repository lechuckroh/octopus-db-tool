# Liquibase

[한국어](kr/liquibase.md)

## Generate

```shell
$ oct generate liquibase --help
```

|           Option           |        Env. Variable         | Description                                                                   |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | Octopus schema file to read                                                   |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | Target file or directory                                                      |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | Table groups to generate.<br />Set multiple groups with comma(`,`) separated. |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | Unique constraint name suffix                                                 |
|     `-c`, `--comments`     |      `OCTOPUS_COMMENTS`      | Set flag to generate column comments. Default: `false`                        |

### Example

```shell
$ oct generate liquibase \
    --input examples/user.json \
    --output output/changelogs.yml \
    --comments
```

Generated yaml file:

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

## Diff

Generate diff changelogs comparing 2 schema files.

```shell
$ oct diff liquibase --help
```

|           Option           |        Env. Variable         | Description                                                                  |
| :------------------------: | :--------------------------: | :--------------------------------------------------------------------------- |
|      `-a`, `--author`      |       `OCTOPUS_AUTHOR`       | Diff author                                                                  |
|       `-f`, `--from`       |        `OCTOPUS_FROM`        | Octopus schema to compare 'from'                                             |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | Table groups to compare.<br />Set multiple groups with comma(`,`) separated. |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | Diff output file                                                             |
|        `-t`, `--to`        |         `OCTOPUS_TO`         | Octopus schema to compare 'to'                                               |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | Unique constraint name suffix                                                |
|     `-c`, `--comments`     |      `OCTOPUS_COMMENTS`      | Set flag to generate column comments. Default: `false`                       |

### Example

```shell
$ oct diff liquibase \
    --from examples/user.json \
    --to examples/user-v2.json \
    --output output/diff.yaml \
    --uniqueNameSuffix _uq \
    --author foo \
    --comments
```

Generated yaml file:

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
