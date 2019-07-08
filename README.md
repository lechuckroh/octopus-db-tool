# octopus-db-tools
octopus-db-tools provides:
* DB schema format conversion
* Generate SQL to create tables.
* Generate entity source codes

## Requirements
### Local Build
* Golang 1.12 or higher
* make

### Docker build
* docker
* docker-compose

## Build
### Local Build
```bash
$ make vendor
$ make compile
```

### Docker Build
```bash
$ make compile-docker; make compile-rmi
```

## Run
```bash
# create example schema file
$ ./oct create sample.ojson

# convert format
$ ./oct convert <sourceFile> <targetFile> \
    --sourceFormat=<srcFormat> \
    --targetFormat=<targetFormat>

# generate files
$ ./oct generate <sourceFile> <targetFile> \
    --sourceFormat=<srcFormat> \
    --targetFormat=<targetFormat> \
    --package=<packageName> \
    --removePrefix=<prefixes>
```

### Convert
```bash
# starUML2 -> octopus
$ ./oct convert sample.mdj sample.ojson --sourceFormat=staruml2 --targetFormat=octopus

# octopus -> xlsx
$ ./oct convert sample.ojson sample.xlsx --sourceFormat=octopus --targetFormat=xlsx
```

### Generate
```bash
# octopus -> JPA-kotlin
# - package: com.foo
# - output directory: ./output
# - remove tableName prefix starting with 'db_' or 'mydb_'
$ ./oct generate sample.ojson ./output --sourceFormat=octopus --targetFormat=jpa-kotlin --package=com.foo --removePrefix=db_,mydb_
```
