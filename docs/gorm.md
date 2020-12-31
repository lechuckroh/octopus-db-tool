# GORM

## Generate

```bash
$ oct generate gorm --help
```

```
OPTIONS:
   --input FILE, -i FILE               read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              generate GORM source file(s) to FILE/`DIR` [$OCTOPUS_OUTPUT]
   --gormModel value, -m value         set embedded base model for GORM model [$OCTOPUS_GORM_MODEL]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --package value, -k value           set package name [$OCTOPUS_PACKAGE]
   --prefix value, -p value            set model struct name prefix [$OCTOPUS_PREFIX]
   --removePrefix value, -r value      set prefixes to remove from model struct name. set multiple values with comma separated. [$OCTOPUS_REMOVE_PREFIX]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
```

Generate `*.go` file:

```bash
$ oct generate gorm \
    -i database.json \
    -o databse.go \
    -p model \
    -r tbl_,table_ \
    -p foo:F,bar:B \
    -u _uq \
    -g foo,bar
```
