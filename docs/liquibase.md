# Liquibase

## Generate

```shell
$ oct generate liquibase --help
```

```
OPTIONS:
   --input FILE, -i FILE               read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              export liquibase changelogs to FILE [$OCTOPUS_OUTPUT]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --comments, -c                      set true to generate column comments (default: false) [$OCTOPUS_COMMENTS]
```

Generate all:
* output directory: `./output`
* unique constraint Name : tableName + `_uq`
* generate comments

```shell
$ oct generate liquibase \
    --input database.json \
    --output ./output/changelogs.yml \
    --uniqueNameSuffix _uq \
    --comments
```

### Example

```shell
$ oct generate liquibase \
    --input examples/user.json \
    --output ./output/changelogs.yml
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
        - column:
            name: name
            type: varchar(40)
            constraints:
              unique: true
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
        - column:
            name: name
            type: varchar(40)
            constraints:
              unique: true
        - column:
            name: group_id
            type: bigint
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

```shell
$ oct diff liquibase --help
```

```
OPTIONS:
   --author value, -a value            diff author [$OCTOPUS_AUTHOR]
   --from value, -f value              octopus schema to compare 'from' [$OCTOPUS_FROM]
   --groups value, -g value            filter table groups to compare. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --output FILE, -o FILE              diff output FILE [$OCTOPUS_OUTPUT]
   --to value, -t value                octopus schema to compare 'to' [$OCTOPUS_TO]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --comments, -c                      set true to compare column comments (default: false) [$OCTOPUS_COMMENTS]
   --help, -h                          show help (default: false)
```

Generate diff changelog:
* output file: `diff.yaml`
* unique constraint Name : tableName + `_uq`
* from octopus: `v1.json`
* to octopus: `v2.json`
* diff author: `foo`  
* generate comments diff

```shell
$ oct diff liquibase \
    --from v1.json \
    --to v2.json \
    --output diff.yaml \
    --uniqueNameSuffix _uq \
    --author foo \
    --comments
```
