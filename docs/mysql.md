# MySQL

[한국어](kr/mysql.md)

## Import

```shell
$ oct import mysql --help
```

|       Option       |   Env. Variable    | Description                                |
| :----------------: | :----------------: | :----------------------------------------- |
|  `-a`, `--author`  |  `OCTOPUS_AUTHOR`  | Import with author                         |
|  `-i`, `--input`   |  `OCTOPUS_INPUT`   | Input mysql DDL file                       |
|  `-o`, `--output`  |  `OCTOPUS_OUTPUT`  | Output octopus output file                 |
| `-x`, `--excludes` | `OCTOPUS_EXCLUDES` | Tables to exclude. Separated by comma(`,`) |
| `-v`, `--version`  | `OCTOPUS_VERSION`  | Import with version                        |

### Example

Import existing mysql DB:

```shell
$ mysqldump -u {user} -p{password} -h {host} --no-data {database} > mysql-ddl.sql
$ oct import mysql --input mysql-ddl.sql --output database.json

# use this if you get error: Unknown table 'column_statistics' in information_schema (1109)
$ mysqldump -u {user} -p{password} -h {host} --no-data --column-statistics=0 {database} > mysql-ddl.sql
```

## Export

```shell
$ oct export mysql --help
```

|           Option           |        Env. Variable         | Description                                                                   |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | Input octopus schema file                                                     |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | Output mysql DDL file                                                         |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | Table groups to generate.<br />Set multiple groups with comma(`,`) separated. |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | Unique constraint name suffix                                                 |

### Example

```shell
$ oct export mysql \
    --input examples/user.json \
    --output output/user.sql
```

Exported DDL file:

```sql
CREATE TABLE IF NOT EXISTS group (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'unique id',
  name varchar(40) NOT NULL COMMENT 'group name',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group` (`name`)
);
CREATE TABLE IF NOT EXISTS user (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'unique id',
  name varchar(40) NOT NULL COMMENT 'user login name',
  group_id bigint COMMENT 'group ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user` (`name`)
);
```
