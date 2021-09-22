# GraphQL

## Generate

```shell
$ oct generate graphql --help
```

```
OPTIONS:
   --input FILE, -i FILE             read octopus schema from FILE [$OCTOPUS_INPUT]
   --output DIR, -o DIR              generate graphql filess to DIR [$OCTOPUS_OUTPUT]
   --graphqlPackage value, -p value  set target graphql package name [$OCTOPUS_GRAPHQL_PACKAGE]
```

Generate `*.graphql` files:

```shell
# example with all CLI options
$ oct generate graphql \
    --input database.json \
    --output databse.graphql \
    --graphqlPackage my.graphql
```

### Example

```shell
$ oct generate graphql \
    --input examples/user.json \
    --output output/graphql/
```

Generated graphql file:

```graphql
schema {
    query: Query
}

type Query {
  groups: [Group]
  users: [User]
}

type Group {
  id: ID!
  name: String
}

type User {
  id: ID!
  name: String
  groupId: Int
}
```
