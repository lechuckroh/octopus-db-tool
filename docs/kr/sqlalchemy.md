# SqlAlchemy

[English](../sqlalchemy.md)

SQLAlchemy 는 파이썬 ORM 라이브러리입니다.

- [Homepage](https://www.sqlalchemy.org/)

## 소스 생성

```shell
$ oct generate sqlalchemy --help
```

|            옵션            |           환경변수           | 설명                                                                                         |
| :------------------------: | :--------------------------: | :------------------------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                                                        |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | 출력할 파일명/디렉토리명. `*.py` 형식인 경우 한 파일에 모든 클래스를 생성합니다.             |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분                            |
|      `-p`, `--prefix`      |       `OCTOPUS_PREFIX`       | 클래스 이름 접두사.<br />형식: `<그룹1>:<접두사1>[,<그룹2>:<접두사2>]...`                    |
|   `-r`, `--removePrefix`   |   `OCTOPUS_REMOVE_PREFIX`    | 생성할 클래스명에서 제거할 접두사.<br />여러개의 접두사 지정시 `,`로 구분                    |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | 유니크 제약 이름 접미사                                                                      |
|      `-t`, `--useUTC`      |      `OCTOPUS_USE_UTC`       | 플래그 설정시 audit 컬럼(`created_at`, `updated_at`)들에 대해 UTC 사용.<br />기본값: `false` |

### 예제

```shell
$ oct generate sqlalchemy \
    --input examples/user.json \
    --output output/user.py
```

`*.py` 파일은 다음과 같이 생성됩니다:

```python
from sqlalchemy import BigInteger, Column, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy_repr import RepresentableBase

Base = declarative_base(cls=RepresentableBase)


class UserGroup(Base):
    __tablename__ = 'group'

    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String(40), unique=True, nullable=False)


class User(Base):
    __tablename__ = 'user'

    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String(40), unique=True, nullable=False)
    group_id = Column(BigInteger)
```
