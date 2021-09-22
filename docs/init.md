# Initialize

[한국어](kr/init.md)

```shell
$ oct init --help
```

|      Option      |  Env. Variable   | Description                                           |
| :--------------: | :--------------: | :---------------------------------------------------- |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | Octopus schema file to create.<br/>Default: `db.json` |

## Create a new file

```shell
$ oct init --output user.json
```

Generated `*.json` file:

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
