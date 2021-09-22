# Initialize

[English](../init.md)

```shell
$ oct init --help
```

|       옵션       |     환경변수     | 설명                                              |
| :--------------: | :--------------: | :------------------------------------------------ |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 생성할 octopus 스키마 파일.<br/>기본값: `db.json` |

## 파일 생성

```shell
$ oct init --output user.json
```

`*.json` 파일은 다음과 같이 생성됩니다:

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
      "desc": "User table"
    }
  ]
}
```
