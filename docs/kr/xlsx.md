# xlsx

[English](../xlsx.md)

## 임포트

```shell
$ oct import xlsx --help
```

|       옵션       |     환경변수     | 설명                       |
| :--------------: | :--------------: | :------------------------- |
| `-i`, `--input`  | `OCTOPUS_INPUT`  | 임포트할 엑셀 파일명       |
| `-o`, `--output` | `OCTOPUS_OUTPUT` | 저장할 octopus 스키마 파일 |

### 예제

```shell
$ oct import xlsx \
    --input user.xlsx \
    --output user.json
```

## 내보내기

```shell
$ oct export xlsx --help
```

|       옵션        |         환경변수          | 설명                                                                            |
| :---------------: | :-----------------------: | :------------------------------------------------------------------------------ |
|  `-i`, `--input`  |      `OCTOPUS_INPUT`      | 입력으로 사용할 octopus 스키마 파일명                                           |
| `-o`, `--output`  |     `OCTOPUS_OUTPUT`      | 출력할 엑셀 파일                                                                |
| `--useNullColumn` | `OCTOPUS_USE_NULL_COLUMN` | nullable 컬럼 사용 플래그. `false`인 경우 `not null` 컬럼 사용. 기본값: `false` |

### 예제

```shell
$ oct export xlsx \
    --input examples/user.json \
    --output output/user.xlsx
```

`*.xlsx` 파일은 다음과 같이 생성됩니다:

| table/ref. |  column  |    type     | key | not null |   attributes    |   description   |
| :--------: | :------: | :---------: | :-: | :------: | :-------------: | :-------------: |
|   group    |          |    table    |     |          | class=UserGroup |   Group table   |
|            |    id    |    int64    |  P  |    O     |     autoinc     |    unique id    |
|            |   name   | varchar(40) |  U  |    O     |                 |   group name    |
|            |          |             |     |          |                 |                 |
|    user    |          |    table    |     |          |                 |   User table    |
|            |    id    |    int64    |  P  |    O     |     autoInc     |    unique id    |
|            |   name   | varchar(40) |  U  |    O     |                 | user login name |
| >group.id  | group_id |    int64    |     |          |                 |    group ID     |

## 엑셀 시트

### `Meta` 시트

`Meta` 시트는 스키마 메타 데이터를 저장합니다.

|    키     |      값       |
| :-------: | :-----------: |
| `author`  | 스키마 작성자 |
|  `name`   |  스키마 이름  |
| `version` |  스키마 버전  |

### 그룹 시트

테이블 그룹 이름별로 시트가 생성됩니다. 그룹 미지정시 `Common` 시트가 생성됩니다.

## 엑셀 테이블 형식

### `table/ref.` 컬럼

- 첫번째 행: 테이블 명
- 이후 행: 컬럼 연관관계 설정.
  - 패턴: `{연관관계}{테이블}.{컬럼}`
  - 연관관계
    - `>`: 다대일 관계 (N:1)
    - `<`: 일대다 관계 (1:N)
    - `-`: 일대일 관계 (1:1)
- `key` 컬럼이 `I`인 경우 인덱스명으로 사용.
  - 동일한 인덱스명이 설정된 컬럼이 여러개 있는 경우, 해당 컬럼들이 하나의 인덱스를 구성합니다.

### `column` 컬럼

컬럼명

### `type` 컬럼

첫번째 행은 반드시 `table` 타입이어야 합니다.

사용 가능한 타입 목록:

- `binary`
- `bit`
- `blob16`
- `blob24`
- `blob32`
- `blob8`
- `boolean`
- `char`
- `date`
- `datetime`
- `decimal`
- `double`
- `enum`
- `float`
- `geometry`
- `int16`
- `int24`
- `int32`
- `int64`
- `int8`
- `json`
- `point`
- `set`
- `text16`
- `text24`
- `text32`
- `text8`
- `time`
- `varbinary`
- `varchar`
- `year`

컬럼 크기는 다음과 같이 설정할 수 있습니다:

- `varchar(40)`
- `decimal(5,2)`
  - `5`: 전체 자릿수
  - `2`: 소수점 이하 자릿수

### `key` 컬럼

키와 관련된 제약을 설정합니다:

- `P`: 기본 키 (Primary Key)
- `U`: 유니크 키 (Unique Key)
- `I`: 인덱스 (Index)

### `not null` / `nullable` 컬럼

- 사용여부는 `O` 를 입력합니다.
- 컬럼 헤더는 다음 중 하나를 사용합니다.
  - `not null`: 체크된 컬럼은 not null로 설정
  - `nullable`: 체크된 컬럼은 nullable로 설정

### `attributes` 컬럼

여러개의 속성은 `,`로 구분해 나열합니다.

- `autoinc`: Auto Incremental 컬럼
- `default={value}`: 컬럼의 기본값 설정
  - 함수를 사용하려면 함수명 앞에 `fn::` 접두사를 사용합니다.
  - 함수는 파라미터를 가질 수 없습니다.
  - 함수 지정시 `()` 없이 함수명만 설정합니다.
  - `default=fn::CURRENT_TIMESTAMP` 는 mysql 의 경우 `DEFAULT CURRENT_TIMESTAMP()` 로 변환됩니다.
- `class={value}`: 소스 생성시 사용할 클래스명을 지정합니다.
  - 이 속성은 `type` 컬럼의 값이 `table`인 경우에만 유효합니다.

### `description` 컬럼

컬럼 설명

## 예제

| table/ref. |   column   |    type     | key | not null |                          attributes                          | description |
| :--------: | :--------: | :---------: | :-: | :------: | :----------------------------------------------------------: | :---------: |
|   group    |            |    table    |     |          |                       class=UserGroup                        | User Group  |
|            |     id     |    int64    |  P  |    O     |                           autoinc                            |             |
|            |    name    | varchar(20) |  U  |    O     |                                                              | group name  |
|            | created_at |  datetime   |     |    O     |                default=fn::CURRENT_TIMESTAMP                 |             |
|            | updated_at |  datetime   |     |    O     | default=fn::CURRENT_TIMESTAMP,onUpdate=fn::CURRENT_TIMESTAMP |             |
|            |            |             |     |          |                                                              |             |
|    user    |            |    table    |     |          |                                                              |    User     |
|            |     id     |    int64    |  P  |    O     |                           autoInc                            |   user id   |
|            |    name    | varchar(40) |     |    O     |                                                              |  user name  |
| >group.id  |  group_id  |    long     |     |    O     |                                                              |  group id   |
|  user_idx  |    name    |             |  I  |          |                                                              |             |
|  user_idx  |  group_id  |             |  I  |          |                                                              |             |

위와 같이 엑셀 시트를 정의한 경우, mysql DDL 은 다음과 같이 생성됩니다:

```sql
CREATE TABLE IF NOT EXISTS group (
  id bigint NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL COMMENT 'group name',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_uq` (`name`)
);

CREATE TABLE IF NOT EXISTS user (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'user id',
  name varchar(40) NOT NULL COMMENT 'user name',
  group_id bigint NOT NULL COMMENT 'group id',
  PRIMARY KEY (`id`),
  INDEX `user_idx` (`name`, `group_id`)
);
```
