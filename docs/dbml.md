# DBML

## Export

```shell
$ oct export dbml --help
```

```
OPTIONS:
   --input FILE, -i FILE     read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE    export DBML to FILE [$OCTOPUS_OUTPUT]
   --groups value, -g value  filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
```

Export `*.dbml` file:

```shell
# example with all CLI options
$ oct export dbml \
    --input database.json \
    --output databse.dbml \
    --groups foo,bar
```

### Example

```shell
$ oct export dbml \
    --input examples/user.json \
    --output output/user.dbml
```

Exported dbml file:

```
Table group {
  id int64 [pk]
  name varchar(40) [unique]
}

Table user {
  id int64 [pk]
  name varchar(40) [unique]
  group_id int64 [ref: > group.id]
}
```
