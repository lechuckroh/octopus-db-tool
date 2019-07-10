# octopus-db-tools
octopus-db-tools provides:
* ERD definition format conversion
* Generate SQL to create tables.
* Generate entity source codes

## Supported Formats
|  format  |input|output|generate|extension|
|---------------------|:-:|:-:|:-:|:------:|
| `octopus`           | O | O |   |`ojson` |
| `xlsx`              | O | O |   |`xlsx`  |
| `staruml2`          | O |   |   |`mdj`   |
| [`dbdiagram.io`][1] |   | O |   |        |
| [`quickdbd`][2]     |   | O |   |        |
| `jpa-kotlin`        |   |   | O |`kt`    |
| `jpa-groovy`        |   |   |   |`groovy`|
| `jpa-java`          |   |   |   |`java`  |
| `liquibase`         |   |   | O |`yaml`  |
| `opti-studio`       |   |   |   |`xml`   |
| `plantuml`          |   |   |   |`plantuml`|
| `schema-converter`  |   |   |   |`schema`|
| `sql-h2`            |   |   |   |`sql`   |
| `sql-mysql`         |   |   |   |`sql`   |
| `sql-oracle`        |   |   |   |`sql`   |
| `sql-sqlite3`       |   |   |   |`sql`   |
| `sql-sqlserver`     |   |   |   |`sql`   |


[1]: https://dbdiagram.io/
[2]: https://www.quickdatabasediagrams.com/

## Build
### Local Build
Requirements:
* Golang 1.12 or higher
* make

Run:
```bash
$ make vendor
$ make compile

# build os-specific binary
$ make compile-windows
$ make compile-linux
$ make compile-macos
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
    --removePrefix=<prefixes> \
    --uniqueNameSuffix=<suffix>
```

You can omit `--sorceFormat`, `--targetFormat` if file format can be detected.

### Convert
```bash
# starUML2 -> octopus
$ ./oct convert sample.mdj sample.ojson

# octopus -> xlsx
$ ./oct convert sample.ojson sample.xlsx
```

### Generate
```bash
# octopus -> JPA-kotlin
# - package: com.foo
# - output directory: ./output
# - remove tableName prefix starting with 'db_' or 'mydb_'
# - unique constraint Name : tableName + '_uq'
$ ./oct generate sample.ojson ./output \
    --targetFormat=jpa-kotlin \
    --package=com.foo \
    --removePrefix=db_,mydb_ \
    --uniqueNameSuffix=_uq
```
