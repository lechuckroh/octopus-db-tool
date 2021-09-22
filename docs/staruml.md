# Star UML

[한국어](kr/staruml.md)

* StarUML [Homepage](https://staruml.io/)

## Import

```shell
$ oct import ojson --help
```

|      Option      |  Env. Variable   | Description                      |
| :--------------: | :--------------: | :------------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | StarUML file to import |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | Target octopus schema file       |

### Example

Import starUML file:

```shell
$ oct import staruml \
    --input user.uml \
    --output user.json 
```
