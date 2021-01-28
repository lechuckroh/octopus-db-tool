# DBML

**not implemented**

## Export

```bash
$ oct export dbml --help
```

```
OPTIONS:
   --input FILE, -i FILE     read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE    export DBML to FILE [$OCTOPUS_OUTPUT]
   --groups value, -g value  filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
```

Export `*.dbml` file:

```bash
$ oct export dbml \
    --input database.json \
    --output databse.dbml \
    --groups foo,bar
```
