# Octopus 파일 형식

[English](../octopus-format.md)

|   이름    |        타입         | 설명             |
| :-------: | :-----------------: | ---------------- |
| `author`  |      `string`       | DB 스키마 작성자 |
|  `name`   |      `string`       | DB 스키마 이름   |
| `version` |      `string`       | DB 스키마 버전   |
| `tables`  | [Table](#table)`[]` | DB 테이블 목록   |

## Table

DB 테이블 정의.

|    이름     |         타입          | 설명                                     |
| :---------: | :-------------------: | ---------------------------------------- |
|   `name`    |       `string`        | 테이블 명                                |
|  `columns`  | [Column](#column)`[]` | 컬럼 목록                                |
|   `desc`    |       `string`        | 테이블 설명                              |
|   `group`   |       `string`        | 테이블을 논리적으로 구분하기 위한 그룹명 |
| `className` |       `string`        | 생성할 클래스 명. ORM 코드 생성시 사용   |
|  `indices`  |  [Index](#index)`[]`  | 인덱스 목록                              |

## Column

DB 컬럼 정의.

|    이름    |            타입             | 설명                                                       | 기본값  |
| :--------: | :-------------------------: | :--------------------------------------------------------- | :-----: |
|   `name`   |          `string`           | 컬럼 명                                                    |         |
|   `type`   |          `string`           | 컬럼 타입. [데이터타입](#datatypes) 참고                   |         |
|   `desc`   |          `string`           | 컬럼 설명                                                  |         |
|   `size`   |            `int`            | 컬럼 길이                                                  |         |
|  `scale`   |            `int`            | 컬럼 소수점 이하 자리수                                    |         |
| `notnull`  |          `boolean`          | not null 컬럼 여부                                         | `false` |
|    `pk`    |          `boolean`          | primary key 컬럼 여부                                      | `false` |
|  `unique`  |          `boolean`          | unique key 컬럼 여부                                       | `false` |
| `autoinc`  |          `boolean`          | Auto incremental 여부                                      | `false` |
| `default`  |    `string` / `function`    | 기본값 또는 기본값 생성 함수                               |         |
| `onupdate` |    `string` / `function`    | 행이 변경되는 경우 업데이트할 값 또는 함수                 |         |
|  `values`  |         `string[]`          | 사용가능한 값 목록.<br />mysql의 `enum`, `set` 타입과 동일 |         |
|   `ref`    | [Reference](#reference)`[]` | 참조하는 다른 테이블 컬럼                                  |         |

### `function` 타입

`function` 타입은 `fn::` 접두사를 사용합니다.

예를 들어, 다음과 같은 DDL을 정의한다고 가정합니다:

```sql
CREATE TABLE t1 (
  ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

`ts` 컬럼은 다음과 같이 정의할 수 있습니다:

```json
{
  "name": "ts",
  "type": "timestamp",
  "default": "fn::CURRENT_TIMESTAMP",
  "onupdate": "fn::CURRENT_TIMESTAMP"
}
```

## Reference

외래키와 같이 연관된 다른 테이블의 컬럼을 정의합니다:

|      이름       |   타입   | 설명                                                                              |
| :-------------: | :------: | :-------------------------------------------------------------------------------- |
|     `table`     | `string` | 연관 컬럼의 테이블 명                                                             |
|    `column`     | `string` | 연관 컬럼 명                                                                      |
| -`relationship` | `string` | <ul><li>`1:1`: 1:1 매핑</li><li>`1:n`: 1:N 매핑</li><li>`n:1`: N:1 매핑</li></ul> |

## Index

DB 인덱스 정의

|   이름    |    타입    | 설명                        |
| :-------: | :--------: | :-------------------------- |
|  `name`   |  `string`  | 인덱스 명                   |
| `columns` | `string[]` | 인덱스를 구성하는 컬럼 목록 |

## DataTypes

octopus 데이터 타입:

|     이름     | 설명                                                                          |  MySQL 매핑  |
| :----------: | ----------------------------------------------------------------------------- | :----------: |
|   `string`   | 문자열                                                                        |  `varchar`   |
|  `tinyint`   | 1 바이트 정수                                                                 |  `tinyint`   |
|  `smallint`  | 2 바이트 정수                                                                 |  `smallint`  |
| `mediumint`  | 3 바이트 정수                                                                 | `mediumint`  |
|    `int`     | 4 바이트 정수                                                                 |    `int`     |
|  `integer`   | 4 바이트 정수                                                                 |    `int`     |
|   `bigint`   | 8 바이트 정수                                                                 |   `bigint`   |
|    `long`    | 8 바이트 정수                                                                 |   `bigint`   |
|  `numeric`   | Decimal 형식<br />`size`: 전체 숫자 자리수<br />`scale`: 소수점 이하 자리수   |  `decimal`   |
|    `real`    | 실수 형식<br />`size`: 전체 유효 자리수<br />`scale`: 소수점 이하 유효 자리수 |   `double`   |
| `timestamp`  | 날짜/시간                                                                     |  `datetime`  |
|  `tinyblob`  | 2<sup>8</sup> 바이트 바이너리 객체                                            |  `tinyblob`  |
|    `blob`    | 2<sup>16</sup> 바이트 바이너리 객체                                           |    `blob`    |
| `mediumblob` | 2<sup>24</sup> 바이트 바이너리 객체                                           | `mediumblob` |
|  `longblob`  | 2<sup>32</sup> 바이트 바이너리 객체                                           |  `longblob`  |
|  `tinytext`  | 2<sup>8</sup> 바이트 텍스트                                                   |  `tinytext`  |
|    `text`    | 2<sup>16</sup> 바이트 텍스트                                                  |    `text`    |
| `mediumtext` | 2<sup>24</sup> 바이트 텍스트                                                  | `mediumtext` |
|  `longtext`  | 2<sup>32</sup> 바이트 텍스트                                                  |  `longtext`  |

## 예제

```json
{
  "version": "1.0.0",
  "tables": [
    {
      "name": "group",
      "columns": [
        {
          "name": "id",
          "type": "long",
          "desc": "unique id",
          "pk": true,
          "autoinc": true
        },
        {
          "name": "name",
          "type": "string",
          "desc": "group name",
          "size": 40,
          "unique": true
        }
      ],
      "desc": "Group table"
    },
    {
      "name": "user",
      "columns": [
        {
          "name": "id",
          "type": "long",
          "desc": "unique id",
          "pk": true,
          "autoinc": true
        },
        {
          "name": "name",
          "type": "string",
          "desc": "user login name",
          "size": 40,
          "unique": true
        },
        {
          "name": "group_id",
          "type": "long",
          "desc": "group ID",
          "ref": {
            "table": "group",
            "column": "id"
          }
        }
      ],
      "desc": "User table",
      "indices": [
        {
          "name": "group_id_index",
          "columns": ["group_id"]
        }
      ]
    }
  ]
}
```
