# Quick DBD

[English](../quickdbd.md)

## 내보내기

```shell
$ oct export quickdbd --help
```

|       옵션       |     환경변수     | 설명                                  |
| :--------------: | :--------------: | :------------------------------------ |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 입력으로 사용할 octopus 스키마 파일명 |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 생성할 quickDBD 파일명                |

### 예제

```shell
$ oct export quickdbd \
    --input examples/user.json \
    --output output/user.txt
```

`*.txt` 파일은 다음과 같이 생성됩니다:

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
