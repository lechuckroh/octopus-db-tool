# ojson (octopus-db-tools v1)

[한국어](kr/ojson.md)

## Import

Import octopus v1 schema file.

```shell
$ oct import ojson --help
```

|      Option      |  Env. Variable   | Description                      |
| :--------------: | :--------------: | :------------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | Octopus v1 schema file to import |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | Target octopus schema file       |

### Example

Import `*.ojson` file:

```shell
$ oct import ojson \
    --input database.ojson \
    --output databse.json
```
