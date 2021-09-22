# GraphQL

[한국어](kr/graphql.md)

## Generate

```shell
$ oct generate graphql --help
```

|          Option          |       Env. Variable       | Description                 |
| :----------------------: | :-----------------------: | :-------------------------- |
|     `-i`, `--input`      |      `OCTOPUS_INPUT`      | Octopus schema file to read |
|     `-o`, `--output`     |     `OCTOPUS_OUTPUT`      | Target directory            |
| `-p`, `--graphqlPackage` | `OCTOPUS_GRAPHQL_PACKAGE` | Target graphql package name |

### Example

```shell
$ oct generate graphql \
    --input examples/user.json \
    --output output/graphql
```

Generated `*.graphql` file:

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
