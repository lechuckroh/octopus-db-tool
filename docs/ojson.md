# ojson (octopus-db-tools v1)

## Import

```shell
$ oct import ojson --help
```

```
OPTIONS:
   --input FILE, -i FILE   import octopus v1 schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  write octopus schema to FILE [$OCTOPUS_OUTPUT]
```

Import `*.ojson` file:

```shell
$ oct import ojson \
    --input database.ojson \
    --output databse.json
```
