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
