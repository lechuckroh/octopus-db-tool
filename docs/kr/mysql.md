# MySQL

[English](../mysql.md)

## DDL 임포트

```shell
$ oct import mysql --help
```

|        옵션        |      환경변수      | 설명                                         |
| :----------------: | :----------------: | :------------------------------------------- |
|  `-a`, `--author`  |  `OCTOPUS_AUTHOR`  | octopus 스키마 파일에 설정할 작성자          |
|  `-i`, `--input`   |  `OCTOPUS_INPUT`   | 임포트할 mysql DDL 파일                      |
|  `-o`, `--output`  |  `OCTOPUS_OUTPUT`  | 저장할 octopus 스키마 파일                   |
| `-x`, `--excludes` | `OCTOPUS_EXCLUDES` | 임포트하지 않을 테이블 목록. `,`로 구분한다. |
| `-v`, `--version`  | `OCTOPUS_VERSION`  | octopus 스키파 파일에 설정할 버전            |

### 예제

다음과 같은 방법으로 Mysql DB 테이블을 octopus 스키마 파일로 임포트할 수 있습니다.

```shell
$ mysqldump -u {user} -p{password} -h {host} --no-data {database} > mysql-ddl.sql
$ oct import mysql --input mysql-ddl.sql --output database.json

# Unknown table 'column_statistics' in information_schema (1109) : 이 에러가 발생하는 경우 다음과 같이 실행하세요.
$ mysqldump -u {user} -p{password} -h {host} --no-data --column-statistics=0 {database} > mysql-ddl.sql
```

## DDL 내보내기

```shell
$ oct export mysql --help
```

|            옵션            |           환경변수           | 설명                                                              |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                             |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | 생성할 mysql DDL 파일명                                           |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분 |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | 유니크 제약 이름 접미사                                           |

### 예제

```shell
$ oct export mysql \
    --input examples/user.json \
    --output output/user.sql
```

`*.sql` 파일은 다음과 같이 생성됩니다:

```sql
CREATE TABLE IF NOT EXISTS group (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'unique id',
  name varchar(40) NOT NULL COMMENT 'group name',
  PRIMARY KEY (`id`),
  UNIQUE KEY `group` (`name`)
);
CREATE TABLE IF NOT EXISTS user (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'unique id',
  name varchar(40) NOT NULL COMMENT 'user login name',
  group_id bigint COMMENT 'group ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user` (`name`)
);
```
