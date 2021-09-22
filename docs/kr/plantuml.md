# plantuml

[English](../plantuml.md)

## 파일 생성

```shell
$ oct generate plantuml --help
```

|       옵션       |     환경변수     | 설명                                  |
| :--------------: | :--------------: | :------------------------------------ |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 입력으로 사용할 octopus 스키마 파일명 |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 출력할 PlantUML 파일명/디렉토리       |

지원하는 파일 확장자는 다음과 같습니다.
그 외의 경우 디렉토리명으로 인식하며, `output.plantuml` 파일명이 사용됩니다.

- `*.wsd`
- `*.pu`
- `*.puml`
- `*.plantuml`
- `*.iuml`

### 예제

```shell
$ oct generate plantuml \
    --input examples/user.json \
    --output output/user.puml
```

`*.puml` 파일은 다음과 같이 생성됩니다:

```
@startuml
entity group {
    id: int64 <<PK>>
    --
    name: varchar <<UQ>>
}
entity user {
    id: int64 <<PK>>
    --
    name: varchar <<UQ>>
    group_id: int64 <<FK>>
}
user }o-|| group
@enduml
```

![](../images/plantuml-user.png)
