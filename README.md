# octopus-db-tools

![Release](https://github.com/lechuckroh/octopus-db-tool/actions/workflows/release.yml/badge.svg)
![Test](https://github.com/lechuckroh/octopus-db-tool/actions/workflows/test.yml/badge.svg)

[한국어](README_kr.md)

octopus-db-tools provides:
* Import/Export various ERD definitions.
* Generate various file formats.

## Goals

* All-in-one tool to support every possible DB schema formats.
* DB schema is stored in text format for version control, diff and merge.
* Octopus-db-tool file format can be used as a SSOT(Single Source Of Truth).
* Single binary executable CLI which can be used as a part of CI(Continuous Integration), CD(Continuous Deployment) and IaC(Infrastructure as Code).

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
* Golang 1.17 +
* make

Run:
```shell
$ make vendor
$ make compile

# build os-specific binary
$ make compile-windows
$ make compile-linux
$ make compile-macos
```

### Docker Build
```shell
$ make compile-docker; make compile-rmi
```

### Downloads

See [Releases](https://github.com/lechuckroh/octopus-db-tool/releases) page.

## Run

```shell
# show help
$ ./oct --help
```

See the following pages for command line options.

* [initialize](docs/init.md)
* Commands by format  
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


## Documents

* [octopus file format](docs/octopus-format.md)
