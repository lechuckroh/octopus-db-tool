# ProtoBuf

## Generate

```bash
$ oct generate pb --help
```

```
OPTIONS:
   --input FILE, -i FILE           read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE          generate protobuf definition to FILE [$OCTOPUS_OUTPUT]
   --goPackage value               set go package name [$OCTOPUS_GO_PACKAGE]
   --groups value, -g value        filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --package value, -p value       set package name [$OCTOPUS_PACKAGE]
   --prefix value, -f value        set proto message name prefix [$OCTOPUS_PREFIX]
   --removePrefix value, -d value  set prefixes to remove from message name. set multiple values with comma separated. [$OCTOPUS_REMOVE_PREFIX]
 ```

Generate `*.proto` file:

```bash
$ oct generate pb \
    -i database.json \
    -o ./output \
    --goPackage foo/proto \
    -g foo,bar \
    -p com.foo \
    -f foo:F,bar:B \
    -d tbl_,table_
```
