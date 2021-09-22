# Octopus file format

[한국어](kr/octopus-format.md)

|   Name    |        Type         | Description       |
| :-------: | :-----------------: | ----------------- |
| `author`  |      `string`       | DB schema author  |
|  `name`   |      `string`       | DB schema name    |
| `version` |      `string`       | DB schema version |
| `tables`  | [Table](#table)`[]` | DB table list     |

## Table

Database Table definition.

|    Name     |         Type          | Description                                     |
| :---------: | :-------------------: | ----------------------------------------------- |
|   `name`    |       `string`        | Table name                                      |
|  `columns`  | [Column](#column)`[]` | Column definition list                          |
|   `desc`    |       `string`        | Table description                               |
|   `group`   |       `string`        | Table logical group name                        |
| `className` |       `string`        | Class name to generate. For ORM code generation |
|  `indices`  |  [Index](#index)`[]`  | Index definition list                           |

## Column

Database Column definition.

|    Name    |            Type             | Description                                                          | Default |
| :--------: | :-------------------------: | :------------------------------------------------------------------- | :-----: |
|   `name`   |          `string`           | column name                                                          |         |
|   `type`   |          `string`           | column type. See [DataTypes](#datatypes).                            |         |
|   `desc`   |          `string`           | column description                                                   |         |
|   `size`   |            `int`            | column size                                                          |         |
|  `scale`   |            `int`            | column scale                                                         |         |
| `notnull`  |          `boolean`          | not null column                                                      | `false` |
|    `pk`    |          `boolean`          | primary key column                                                   | `false` |
|  `unique`  |          `boolean`          | unique key column                                                    | `false` |
| `autoinc`  |          `boolean`          | Auto incremental                                                     | `false` |
| `default`  |    `string` / `function`    | default value or function                                            |         |
| `onupdate` |    `string` / `function`    | function or value for `ON UPDATE` (mysql)                            |         |
|  `values`  |         `string[]`          | Permitted values.<br />Equivalent to `enum` and `set` types in mysql |         |
|   `ref`    | [Reference](#reference)`[]` | column reference (relation)                                          |         |

### `function` type

`function` Types are indicated by the `fn::` prefix.

For example, with the following mysql DDL:

```sql
CREATE TABLE t1 (
  ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

`ts` column can be defined as:

```json
{
  "name": "ts",
  "type": "timestamp",
  "default": "fn::CURRENT_TIMESTAMP",
  "onupdate": "fn::CURRENT_TIMESTAMP"
}
```

## Reference

Related columns like foreign key:

|      Name       |   Type   | Description                                                                               |
| :-------------: | :------: | :---------------------------------------------------------------------------------------- |
|     `table`     | `string` | target table name                                                                         |
|    `column`     | `string` | target column name                                                                        |
| -`relationship` | `string` | <ul><li>`1:1`: one to one</li><li>`1:n`: one to many</li><li>`n:1`: many to one</li></ul> |

## Index

Database Index definition.

|   Name    |    Type    | Description            |
| :-------: | :--------: | :--------------------- |
|  `name`   |  `string`  | Index name             |
| `columns` | `string[]` | Index column name list |

## DataTypes

Octopus data type:

|     Name     | Description                                                                                                                      | MySQL mapping |
| :----------: | -------------------------------------------------------------------------------------------------------------------------------- | :-----------: |
|   `string`   | String                                                                                                                           |   `varchar`   |
|  `tinyint`   | 1 byte integer                                                                                                                   |   `tinyint`   |
|  `smallint`  | 2 bytes integer                                                                                                                  |  `smallint`   |
| `mediumint`  | 3 bytes integer                                                                                                                  |  `mediumint`  |
|    `int`     | 4 bytes integer                                                                                                                  |     `int`     |
|  `integer`   | 4 bytes integer                                                                                                                  |     `int`     |
|   `bigint`   | 8 bytes integer                                                                                                                  |   `bigint`    |
|    `long`    | 8 bytes integer                                                                                                                  |   `bigint`    |
|  `numeric`   | Decimal data type<br />`size`: maximum number of digits<br />`scale`: number of digits to the right of the decimal point.        |   `decimal`   |
|    `real`    | Floating point data type<br />`size`: maximum number of digits<br />`scale`: number of digits to the right of the decimal point. |   `double`    |
| `timestamp`  | DateTime                                                                                                                         |  `datetime`   |
|  `tinyblob`  | 2<sup>8</sup> bytes binary large object                                                                                          |  `tinyblob`   |
|    `blob`    | 2<sup>16</sup> bytes binary large object                                                                                         |    `blob`     |
| `mediumblob` | 2<sup>24</sup> bytes binary large object                                                                                         | `mediumblob`  |
|  `longblob`  | 2<sup>32</sup> bytes binary large object                                                                                         |  `longblob`   |
|  `tinytext`  | 2<sup>8</sup> bytes text                                                                                                         |  `tinytext`   |
|    `text`    | 2<sup>16</sup> bytes text                                                                                                        |    `text`     |
| `mediumtext` | 2<sup>24</sup> bytes text                                                                                                        | `mediumtext`  |
|  `longtext`  | 2<sup>32</sup> bytes text                                                                                                        |  `longtext`   |

## Example

```json
{
  "version": "1.0.0",
  "tables": [
    {
      "name": "group",
      "columns": [
        {
          "name": "id",
          "type": "long",
          "desc": "unique id",
          "pk": true,
          "autoinc": true
        },
        {
          "name": "name",
          "type": "string",
          "desc": "group name",
          "size": 40,
          "unique": true
        }
      ],
      "desc": "Group table"
    },
    {
      "name": "user",
      "columns": [
        {
          "name": "id",
          "type": "long",
          "desc": "unique id",
          "pk": true,
          "autoinc": true
        },
        {
          "name": "name",
          "type": "string",
          "desc": "user login name",
          "size": 40,
          "unique": true
        },
        {
          "name": "group_id",
          "type": "long",
          "desc": "group ID",
          "ref": {
            "table": "group",
            "column": "id"
          }
        }
      ],
      "desc": "User table",
      "indices": [
        {
          "name": "group_id_index",
          "columns": ["group_id"]
        }
      ]
    }
  ]
}
```
