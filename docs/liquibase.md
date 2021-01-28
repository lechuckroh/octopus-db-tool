# Liquibase

## Generate

```bash
$ oct generate liquibase --help
```

```
OPTIONS:
   --input FILE, -i FILE               read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              export liquibase changelogs to FILE [$OCTOPUS_OUTPUT]
   --diff value, -d value               [$OCTOPUS_DIFF]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --comments, -c                      set true to generate column comments (default: false) [$OCTOPUS_COMMENTS]
```

Generate all:
* output directory: `./output`
* unique constraint Name : tableName + `_uq`
* generate comments

```bash
$ oct generate liquibase \
    --input database.json \
    --output ./output \
    --uniqueNameSuffix _uq \
    --comments
```

Generate diff changelog:
* output directory: `./output`
* unique constraint Name : tableName + `_uq`
* from octopus: `v1.json`
* to octopus: `v2.json`
* generate comments

```bash
$ oct generate liquibase \
    --input v2.json \
    --output ./output \
    --diff v1.json \
    --uniqueNameSuffix _uq \
    --comments
```
