# ojson (octopus-db-tools v1)

[English](../ojson.md)

## 임포트

octopus v1 스키마 파일을 임포트합니다.

```shell
$ oct import ojson --help
```

|       옵션       |     환경변수     | 설명                            |
| :--------------: | :--------------: | :------------------------------ |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 임포트할 Octopus v1 스키파 파일 |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 저장할 octopus 스키마 파일      |

### 예제

`*.ojson` 파일을 임포트합니다:

```shell
$ oct import ojson \
    --input database.ojson \
    --output databse.json
```
