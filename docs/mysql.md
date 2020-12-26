# MySQL

## Import

### Help
```bash
$ oct import mysql --help
```

```
OPTIONS:
   --input FILE, -i FILE   import mysql DDL from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  export octopus schema to FILE [$OCTOPUS_OUTPUT]
```

### Import mysql DDL

To import mysql DDL file:

```bash 
$ oct import mysql -i database.sql -o database.ojson 
```

## Export

### Help
```bash
$ oct export mysql --help
```

```
OPTIONS:
   --input FILE, -i FILE               load input octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              export mysql DDL to FILE [$OCTOPUS_OUTPUT]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
```

### Export mysql DDL

To export to mysql DDL file with the following options:
* export tables in `common` and `admin` groups
* set unique constraint name suffix: `_uq`

```bash 
$ oct export mysql \
    -i database.ojson \
    -o database.sql \
    -g common,admin \
    -u _uq 
```