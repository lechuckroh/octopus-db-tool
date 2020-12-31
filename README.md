# octopus-db-tools

[![Build Status](https://drone.lechuckcgx.com/api/badges/lechuckroh/octopus-db-tool/status.svg?ref=refs/heads/develop)](https://drone.lechuckcgx.com/lechuckroh/octopus-db-tool)

octopus-db-tools provides:
* Import/Export ERD definitions.
* Generate various file formats.

## Supported Formats

### Import
* Excel (`*.xlsx`)
* MySQL DDL (`*.sql`)
* octopus-db-tools v1 (`*.ojson`)
* StarUML

### Export
* DBML
* Excel (`*.xlsx`)
* MySQL DDL (`*.sql`)

### Generate
* GORM source files (`*.go`)
* GraphQL (`*.graphql`)
* JPA Kotlin (`*.kt`)
* Liquibase (`*.yaml`)
* PlantUML (`*.wsd`, `*.pu`, `*.puml`, `*.plantuml`, `*.iuml`)
* ProtoBuf (`*.proto`)
* [Quick DBD](https://www.quickdatabasediagrams.com/)
* SQLAlchemy (`*.py`)

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
# show help
$ ./oct --help
```

* [initialize](docs/init.md)
* Formats  
    * [DBML](docs/dbml.md)
    * [Excel](docs/xlsx.md)
    * [GORM](docs/gorm.md)
    * [GraphQL](docs/graphql.md)  
    * [JPA](docs/jpa.md)  
    * [Liquibase](docs/liquibase.md)  
    * [MySQL](docs/mysql.md)
    * [octopus-db-tools v1](docs/ojson.md)
    * [ProtoBuf](docs/protobuf.md)
    * [Quick DBD](docs/quickdbd.md)
    * [SQLAlchemy](docs/sqlalchemy.md)
    * [StarUML](docs/staruml.md)
