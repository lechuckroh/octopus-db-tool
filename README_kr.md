# octopus-db-tools

![Release](https://github.com/lechuckroh/octopus-db-tool/actions/workflows/release.yml/badge.svg)
![Test](https://github.com/lechuckroh/octopus-db-tool/actions/workflows/test.yml/badge.svg)

[English](README.md)

octopus-db-tools 기능:
* 다양한 ERD 정의 파일 형식 가져오기/내보내기
* 다양한 형식의 소스 파일 생성

## 목표

* 가능한 모든 종류의 DB 스키마 관련 형식을 지원하는 올인원 도구
* 버전 관리 및 변경 사항 diff, 머지를 쉽게 할 수 있도록 텍스트 파일 형식 사용
* octopus-db-tool 파일 포맷으로 DB 스키마를 정의하면, 다른 포맷의 DB 스키마 문서를 따로 관리하지 않아도 되는 SSOT(Single Source Of Truth)
* CI(Continuous Integration), CD(Continuous Deployment), IaC(Infrastructure as Code) 작업의 일부로 사용될 수 있는 단일 실행파일 커맨드라인 제공

## 지원하는 파일 형식

### 가져오기
* 엑셀 (`*.xlsx`)
* MySQL DDL (`*.sql`)
* octopus-db-tools v1 (`*.ojson`)
* StarUML

### 내보내기
* DBML
* 엑셀 (`*.xlsx`)
* MySQL DDL (`*.sql`)

### 파일 생성
* GORM 소스 파일 (`*.go`)
* GraphQL (`*.graphql`)
* JPA Kotlin (`*.kt`)
* Liquibase (`*.yaml`)
* PlantUML (`*.wsd`, `*.pu`, `*.puml`, `*.plantuml`, `*.iuml`)
* ProtoBuf (`*.proto`)
* [Quick DBD](https://www.quickdatabasediagrams.com/)
* SQLAlchemy (`*.py`)

## 빌드

### 로컬 빌드
필요 항목:
* Golang 1.17 +
* make

실행방법:
```shell
$ make vendor
$ make compile

# 특정 플랫폼용 실행파일 생성을 위한 크로스 컴파일
$ make compile-windows
$ make compile-linux
$ make compile-macos
```

### 도커 컨테이너 내부에서 빌드
```shell
$ make compile-docker; make compile-rmi
```

### 다운로드

[Releases](https://github.com/lechuckroh/octopus-db-tool/releases) 페이지에서 실행파일 다운로드할 수 있습니다.

## 실행방법

```shell
# 도움말 표시
$ ./oct --help
```

각각의 파일 형식별 페이지에서 커맨드라인 옵션을 확인할 수 있습니다:

* [파일 초기화](docs/init.md)
* 파일 형식별 커맨드
    * [DBML](docs/dbml.md)
    * [엑셀](docs/xlsx.md)
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


## 문서

* [octopus 파일 형식](docs/kr/octopus-format.md)
