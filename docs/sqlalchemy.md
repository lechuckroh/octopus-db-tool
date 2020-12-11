# SqlAlchemy

SQLAlchemy is the Python SQL toolkit and Object Relational Mapper that gives application developers the full power and flexibility of SQL.
* [Homepage](https://www.sqlalchemy.org/)

## Generate

### Help
```bash
$ oct generate sqlalchemy --help
```

```
OPTIONS:
   --input FILE, -i FILE               load input octopus schema from FILE [$OCTOPUS_INPUT]
   --output PATH, -o PATH              generate python files to PATH [$OCTOPUS_OUTPUT]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --prefix value, -p value            set entity class name prefix [$OCTOPUS_PREFIX]
   --removePrefix value, -r value      set prefixes to remove from entity class name. set multiple values with comma separated. [$OCTOPUS_REMOVE_PREFIX]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --useUTC value, -t value            use UTC for audit column default value [$OCTOPUS_USE_UTC]
```

### Generate single source file

To generate all entity classes to a single file, set output path to `*.py`:

```bash 
$ oct generate sqlalchemy \
    -i user.ojson 
    -o user.py \
    --uniqueNameSuffix=_uq \
    --groups=common \
    --useUTC=true
```

### Generate multiple source files

To generate entity classes to separate files, set output path to directory:

```bash 
$ oct generate sqlalchemy \
    -i user.ojson 
    -o output/ \
    --uniqueNameSuffix=_uq \
    --groups=common \
    --useUTC=true
```
