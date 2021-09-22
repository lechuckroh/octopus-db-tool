# SqlAlchemy

[한국어](kr/sqlalchemy.md)

SQLAlchemy is the Python SQL toolkit and Object Relational Mapper that gives application developers the full power and flexibility of SQL.

- [Homepage](https://www.sqlalchemy.org/)

## Generate

```shell
$ oct generate sqlalchemy --help
```

|           Option           |        Env. Variable         | Description                                                                               |
| :------------------------: | :--------------------------: | :---------------------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | Octopus schema file to read                                                               |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | Output filename/directory. Set `*.py` to generate in a single file.                       |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | Table groups to generate.<br />Set multiple groups with comma(`,`) separated.             |
|      `-p`, `--prefix`      |       `OCTOPUS_PREFIX`       | Entity class name prefix.<br />Format: `{group1}:{prefix1}[,{group2}:{prefix2}]...`       |
|   `-r`, `--removePrefix`   |   `OCTOPUS_REMOVE_PREFIX`    | Prefixes to remove from class name.<br />Set multiple prefixes with comma(`,`) separated. |
| `-u`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | Unique constraint name suffix                                                             |
|      `-t`, `--useUTC`      |      `OCTOPUS_USE_UTC`       | Set flag to use UTC for audit columns (`created_at`, `updated_at`).<br />Default: `false` |

### Example

```shell
$ oct generate sqlalchemy \
    --input examples/user.json \
    --output output/user.py
```

Generated `*.py` file:

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
