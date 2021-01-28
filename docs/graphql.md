# GraphQL

## Generate

```bash
$ oct generate graphql --help
```

```
OPTIONS:
   --input FILE, -i FILE             read octopus schema from FILE [$OCTOPUS_INPUT]
   --output DIR, -o DIR              generate graphql filess to DIR [$OCTOPUS_OUTPUT]
   --graphqlPackage value, -p value  set target graphql package name [$OCTOPUS_GRAPHQL_PACKAGE]
```

Generate `*.graphql` files:

```bash
$ oct generate graphql \
    --input database.json \
    --output databse.graphql \
    --graphqlPackage my.graphql
```
