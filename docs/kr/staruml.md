# Star UML

[English](../staruml.md)

- StarUML [Homepage](https://staruml.io/)

## 임포트

```shell
$ oct import ojson --help
```

|       옵션       |     환경변수     | 설명                         |
| :--------------: | :--------------: | :--------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 임포트할 starUML 파일명      |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 저장할 octopus 스키마 파일명 |

### 예제

StarUML 파일을 임포트합니다.

```shell
$ oct import staruml \
    --input user.uml \
    --output user.json
```
