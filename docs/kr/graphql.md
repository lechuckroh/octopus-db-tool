# GraphQL

[English](../graphql.md)

## 소스 생성

```shell
$ oct generate graphql --help
```

|           옵션           |         환경변수          | 설명                                  |
| :----------------------: | :-----------------------: | :------------------------------------ |
|     `-i`, `--input`      |      `OCTOPUS_INPUT`      | 입력으로 사용할 octopus 스키마 파일명 |
|     `-o`, `--output`     |     `OCTOPUS_OUTPUT`      | 출력할 디렉토리명                     |
| `-p`, `--graphqlPackage` | `OCTOPUS_GRAPHQL_PACKAGE` | 생성할 graphql 패키지명               |

### 예제

```shell
$ oct generate graphql \
    --input examples/user.json \
    --output output/graphql
```

`*.graphql` 파일은 다음과 같이 생성됩니다:

```graphql
schema {
  query: Query
}

type Query {
  userGroups: [UserGroup]
  users: [User]
}

type UserGroup {
  id: ID!
  name: String!
}

type User {
  id: ID!
  name: String!
  groupId: Int
}
```
