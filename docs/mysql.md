# MySQL

## Import

### Help
```bash
$ oct import mysql --help
```

```
OPTIONS:
   --author value, -a value    import with author [$OCTOPUS_AUTHOR]
   --input FILE, -i FILE       import mysql DDL from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE      write octopus schema to FILE [$OCTOPUS_OUTPUT]
   --excludes value, -x value  tables to exclude. separated by comma [$OCTOPUS_EXCLUDES]
   --version value, -v value   import with version [$OCTOPUS_VERSION]
```

### Import mysql DDL

To import mysql DDL file:

```bash 
$ oct import mysql -i mysql-ddl.sql -o database.json 
```

`mysql-ddl.sql` file can be generated with the following command:
```bash
$ mysqldump -u {user} -p{password} -h {host} --no-data {database} > mysql-ddl.sql

# use this if you get error: Unknown table 'column_statistics' in information_schema (1109)
$ mysqldump -u {user} -p{password} -h {host} --no-data --column-statistics=0 {database} > mysql-ddl.sql
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
    --input database.json \
    --output database.sql \
    --groups common,admin \
    --uniqueNameSuffix _uq 
```
