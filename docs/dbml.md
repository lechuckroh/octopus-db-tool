# DBML

[한국어](kr/dbml.md)

## Export

```shell
$ oct export dbml --help
```

|      Option      |  Env. Variable   | Description                                                                 |
| :--------------: | :--------------: | :-------------------------------------------------------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | Octopus schema file to read                                                 |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | DBML file to write                                                          |
| `-g`, `--groups` | `OCTOPUS_GROUPS` | Table groups to export.<br />Set multiple groups with comma(`,`) separated. |

### Example

```shell
$ oct export dbml \
    --input examples/user.json \
    --output output/user.dbml
```

Exported dbml file:

```
Table group {
  id int64 [pk, not null, note: "unique id"]
  name varchar(40) [unique, not null, note: "group name"]
}

Table user {
  id int64 [pk, not null, note: "unique id"]
  name varchar(40) [unique, not null, note: "user login name"]
  group_id int64 [ref: > group.id, note: "group ID"]
}
```
