# xlsx

[한국어](kr/xlsx.md)

## Import

```shell
$ oct import xlsx --help
```

|      Option      |  Env. Variable   | Description                |
| :--------------: | :--------------: | :------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | Excel file to import       |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | Target octopus schema file |

### Example

Import `*.xlsx` file:

```shell
$ oct import xlsx \
    --input user.xlsx \
    --output user.json
```

## Export

```shell
$ oct export xlsx --help
```

|      Option       |       Env. Variable       | Description                                                              |
| :---------------: | :-----------------------: | :----------------------------------------------------------------------- |
|  `-i`, `--input`  |      `OCTOPUS_INPUT`      | Input octopus schema file                                                |
| `-o`, `--output`  |     `OCTOPUS_OUTPUT`      | Output excel file                                                        |
| `--useNullColumn` | `OCTOPUS_USE_NULL_COLUMN` | Flag to use nullable column. Use `not null` if `false`. Default: `false` |

### Example

```shell
$ oct export xlsx \
    --input examples/user.json \
    --output output/user.xlsx
```

Generated `*.xlsx` file:

| table/ref. |  column  |    type     | key | not null |   attributes    |   description   |
| :--------: | :------: | :---------: | :-: | :------: | :-------------: | :-------------: |
|   group    |          |    table    |     |          | class=UserGroup |   Group table   |
|            |    id    |    int64    |  P  |    O     |     autoinc     |    unique id    |
|            |   name   | varchar(40) |  U  |    O     |                 |   group name    |
|            |          |             |     |          |                 |                 |
|    user    |          |    table    |     |          |                 |   User table    |
|            |    id    |    int64    |  P  |    O     |     autoInc     |    unique id    |
|            |   name   | varchar(40) |  U  |    O     |                 | user login name |
| >group.id  | group_id |    int64    |     |          |                 |    group ID     |

## Excel Sheets

### `Meta` Sheet

`Meta` sheet contains meta data of the schema.

|    key    |        value         |
| :-------: | :------------------: |
| `author`  | Author of the schema |
|  `name`   |     schema name      |
| `version` |    schema version    |

### Group Sheet

Each sheet name represents table group name. `Common` is used if not set.

## Excel Table format

### `table/ref.` column

- 1st row: table name
- after rows: set column reference(s) for columns.
  - pattern: `{relationship}{table}.{column}`
  - relationship
    - `>`: many to one
    - `<`: one to many
    - `-`: one to one
- specify index name if `key` column value is `I`.
  - same index names will be combined into a single index.

### `column` column

Set column name.

### `type` column

The 1st row of the table block should be `table` type.

Possible types:

- `binary`
- `bit`
- `blob16`
- `blob24`
- `blob32`
- `blob8`
- `boolean`
- `char`
- `date`
- `datetime`
- `decimal`
- `double`
- `enum`
- `float`
- `geometry`
- `int16`
- `int24`
- `int32`
- `int64`
- `int8`
- `json`
- `point`
- `set`
- `text16`
- `text24`
- `text32`
- `text8`
- `time`
- `varbinary`
- `varchar`
- `year`

To specify column size:

- `varchar(40)` can hold up to `40` characters.
- `decimal(5,2)`
  - `5`: precision
  - `2`: scale

### `key` column

The value can be either:

- `P`: primary key column
- `U`: unique constraint column
- `I`: part of index column

### `not null` / `nullable` column

- Use `O` to set true.
- The column header should be either:
  - `not null`: checked column cannot be null.
  - `nullable`: checked column is nullable.

### `attributes` column

Enumerate attributes separated by comma(`,`).

- `autoinc`: Auto Incremental column
- `default={value}`: Set default value
  - Use `fn::` prefix to use function.
  - The function should have no parameter.
  - Use function only without `()`.
  - `default=fn::CURRENT_TIMESTAMP` will be converted to `DEFAULT CURRENT_TIMESTAMP()` (in case of mysql).
- `class={value}`: set class name to generate.
  - This attribute is valid only if `type`=`table`.

### `description` column

Set column description.

## Example

| table/ref. |   column   |    type     | key | not null |                          attributes                          | description |
| :--------: | :--------: | :---------: | :-: | :------: | :----------------------------------------------------------: | :---------: |
|   group    |            |    table    |     |          |                       class=UserGroup                        | User Group  |
|            |     id     |    int64    |  P  |    O     |                           autoinc                            |             |
|            |    name    | varchar(20) |  U  |    O     |                                                              | group name  |
|            | created_at |  datetime   |     |    O     |                default=fn::CURRENT_TIMESTAMP                 |             |
|            | updated_at |  datetime   |     |    O     | default=fn::CURRENT_TIMESTAMP,onUpdate=fn::CURRENT_TIMESTAMP |             |
|            |            |             |     |          |                                                              |             |
|    user    |            |    table    |     |          |                                                              |    User     |
|            |     id     |    int64    |  P  |    O     |                           autoInc                            |   user id   |
|            |    name    | varchar(40) |     |    O     |                                                              |  user name  |
| >group.id  |  group_id  |    long     |     |    O     |                                                              |  group id   |
|  user_idx  |    name    |             |  I  |          |                                                              |             |
|  user_idx  |  group_id  |             |  I  |          |                                                              |             |

The above definition is equivalent to the following mysql DDL:

```sql
CREATE TABLE IF NOT EXISTS group (
  id bigint NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL COMMENT 'group name',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_uq` (`name`)
);

CREATE TABLE IF NOT EXISTS user (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'user id',
  name varchar(40) NOT NULL COMMENT 'user name',
  group_id bigint NOT NULL COMMENT 'group id',
  PRIMARY KEY (`id`),
  INDEX `user_idx` (`name`, `group_id`)
);
```
