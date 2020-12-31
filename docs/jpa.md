# JPA

## Generate Kotlin

```bash
$ oct generate kt --help
```

```
OPTIONS:
   --input FILE, -i FILE                 read octopus schema from FILE [$OCTOPUS_INPUT]
   --output DIR, -o DIR                  generate kotlin files to DIR [$OCTOPUS_OUTPUT]
   --annotation FORMAT, -a FORMAT        Custom Entity class annotation. FORMAT: '{group1}:{annotations1}[,{group2}:{annotations2}]' [$OCTOPUS_ANNOTATION]
   --groups GROUPS, -g GROUPS            Filter table groups to generate. GROUPS are separated by comma [$OCTOPUS_GROUPS]
   --idEntity NAME, -e NAME              Interface NAME with 'id' field [$OCTOPUS_ID_ENTITY]
   --package PACKAGE, -p PACKAGE         Entity class PACKAGE name [$OCTOPUS_PACKAGE]
   --prefix FORMAT, -f FORMAT            Class name prefix. FORMAT: '{group1}:{prefix1}[,{group2}:{prefix2}]' [$OCTOPUS_PREFIX]
   --relation ANNOTATION, -l ANNOTATION  Virtual relation ANNOTATION type. Available values: VRelation [$OCTOPUS_RELATION]
   --removePrefix PREFIXES, -d PREFIXES  Table PREFIXES to remove from class name. Multiple prefixes are separated by comma [$OCTOPUS_REMOVE_PREFIX]
   --reposPackage PACKAGE, -r PACKAGE    Repository class PACKAGE name. Generated if not empty. [$OCTOPUS_REPOS_PACKAGE]
   --uniqueNameSuffix SUFFIX, -q SUFFIX  Unique constraint name SUFFIX. [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --useUTC, -u                          Set to use UTC for audit columns ('created_at', 'updated_at'). (default: false) [$OCTOPUS_USE_UTC]
```

Generate `*.kt` files:

```bash
$ oct generate kt \
    -i database.json \
    -o ./output \
    -a foo:@Foo,foobar:@Foo,@Bar \
    -g foo,bar \
    -e IdEntity \
    -p com.foo.entity \
    -f foo:F,bar:B \
    -l VRelation \
    -d tbl_,table_ \
    -r com.foo.repos \
    -q _uq \
    -u    
```
