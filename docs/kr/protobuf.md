# ProtoBuf

[English](../protobuf.md)

## 파일 생성

```shell
$ oct generate pb --help
```

|          옵션          |           환경변수           | 설명                                                                                                                             |
| :--------------------: | :--------------------------: | :------------------------------------------------------------------------------------------------------------------------------- |
|    `-i`, `--input`     |       `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                                                                                            |
|    `-o`, `--output`    |       `OCTOPUS_OUTPUT`       | 출력할 파일명                                                                                                                    |
|     `--goPackage`      |     `OCTOPUS_GO_PACKAGE`     | Protobuf golang 패키지명                                                                                                         |
|    `-g`, `--groups`    |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분                                                                |
|   `-p`, `--package`    |      `OCTOPUS_PACKAGE`       | Protobuf 패키지명                                                                                                                |
|    `-f`, `--prefix`    |       `OCTOPUS_PREFIX`       | 생성할 proto 메시지명의 접두사.<br />형식: `<그룹1>:<접두사1>[,<그룹2>:<접두사2>]...`<br />예제: `group1:prefix1,group2:prefix2` |
| `-d`, `--removePrefix` |   `OCTOPUS_REMOVE_PREFIX`    | proto 메시지명에서 제거할 접두사.<br />여러개의 접두사를 지정시 `,`로 구분                                                       |
|  `--relationTagDecr`   | `OCTOPUS_RELATION_TAG_DECR`  | 연관관계를 나타내는 필드 태그 인덱스를 내림차순으로 사용할지 여부. 기본값: `false`                                               |
|  `--relationTagStart`  | `OCTOPUS_RELATION_TAG_START` | 연관관계를 나타내는 필드 태그 인덱스 시작값. `-1`로 설정할 경우 마지막 필드 다음 인덱스부터 시작합니다.                          |

### 예제

```shell
$ oct generate pb \
    --input examples/user.json \
    --output output/user.proto \
    --goPackage model \
    --package octopus
```

`*.proto` 파일은 다음과 같이 생성됩니다:

```protobuf
syntax = "proto3";

package octopus;

option go_package = "model";

message Group {
  int64 id = 1;
  string name = 2;
}

message User {
  int64 id = 1;
  string name = 2;
  int64 groupId = 3;
  Group group = 4;
}
```
