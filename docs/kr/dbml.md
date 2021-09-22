# DBML

[English](../dbml.md)

## 내보내기

```shell
$ oct export dbml --help
```

|       옵션       |     환경변수     | 설명                                                                |
| :--------------: | :--------------: | :------------------------------------------------------------------ |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 입력으로 사용할 octopus 스키마 파일명                               |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 출력할 DBML 파일명                                                  |
| `-g`, `--groups` | `OCTOPUS_GROUPS` | 내보내기 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분 |

### 예제

```shell
$ oct export dbml \
    --input examples/user.json \
    --output output/user.dbml
```

DBML 파일은 다음과 같이 생성됩니다:

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
