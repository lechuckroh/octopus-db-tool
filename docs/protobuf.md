# ProtoBuf

## Generate

```bash
$ oct generate pb --help
```

```
OPTIONS:
   --input FILE, -i FILE               read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              generate protobuf definition to FILE [$OCTOPUS_OUTPUT]
   --goPackage value                   set go package name [$OCTOPUS_GO_PACKAGE]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --package value, -p value           set package name [$OCTOPUS_PACKAGE]
   --prefix value, -f value            set proto message name prefix [$OCTOPUS_PREFIX]
   --removePrefix value, -d value      set prefixes to remove from message name. set multiple values with comma separated. [$OCTOPUS_REMOVE_PREFIX]
   --relationTagDecr relationTagStart  set relation tags decremental from relationTagStart (default: false) [$OCTOPUS_RELATION_TAG_DECR]
   --relationTagStart value, -s value  set relation tags start index. set -1 to start from last of fields. [$OCTOPUS_RELATION_TAG_START]
```

Generate `*.proto` file:

```bash
$ oct generate pb \
    --input database.json \
    --output ./output \
    --goPackage foo/proto \
    --groups foo,bar \
    --package com.foo \
    --prefix foo:F,bar:B \
    --removePrefix tbl_,table_ \
    --relationTagStart 30
```
