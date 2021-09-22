# Quick DBD

[한국어](kr/quickdbd.md)

## Export

```shell
$ oct export quickdbd --help
```

|      Option      |  Env. Variable   | Description          |
| :--------------: | :--------------: | :------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | Octopus schema file  |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | Output quickDBD file |

### Example

```shell
$ oct export quickdbd \
    --input examples/user.json \
    --output output/user.txt
```

Exported `*.txt` file:

```
group # Group table
-----
id int64 PK AUTOINCREMENT # unique id
name varchar UNIQUE # group name

user # User table
----
id int64 PK AUTOINCREMENT # unique id
name varchar UNIQUE # user login name
group_id int64 NULLABLE FK >- group.id # group ID
```
