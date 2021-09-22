# SqlAlchemy

SQLAlchemy is the Python SQL toolkit and Object Relational Mapper that gives application developers the full power and flexibility of SQL.
* [Homepage](https://www.sqlalchemy.org/)

## Generate

```shell
$ oct generate sqlalchemy --help
```

```
OPTIONS:
   --input FILE, -i FILE               load input octopus schema from FILE [$OCTOPUS_INPUT]
   --output PATH, -o PATH              generate python files to PATH [$OCTOPUS_OUTPUT]
   --groups value, -g value            filter table groups to generate. set multiple values with comma separated. [$OCTOPUS_GROUPS]
   --prefix value, -p value            set entity class name prefix [$OCTOPUS_PREFIX]
   --removePrefix value, -r value      set prefixes to remove from entity class name. set multiple values with comma separated. [$OCTOPUS_REMOVE_PREFIX]
   --uniqueNameSuffix value, -u value  set unique constraint name suffix [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --useUTC value, -t value            use UTC for audit column default value [$OCTOPUS_USE_UTC]
```

To generate all entity classes to a single file, set output path to `*.py`:

```shell 
# example with all CLI options
$ oct generate sqlalchemy \
    --input examples/user.json \
    --output user.py \
    --groups common \
    --removePrefix tbl_,table_ \
    --uniqueNameSuffix _uq \
    --useUTC true
```

To generate entity classes to separate files, set output path to directory:

```shell 
# example with all CLI options
$ oct generate sqlalchemy \
    --input examples/user.json \
    --output output/ \
    --groups common \
    --removePrefix tbl_,table_ \
    --uniqueNameSuffix _uq \
    --useUTC true
```

### Example

```shell
$ oct generate sqlalchemy \
    --input examples/user.json \
    --output output/user.py
```

Generated source file:

```python
from sqlalchemy import BigInteger, Column, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy_repr import RepresentableBase

Base = declarative_base(cls=RepresentableBase)


class Group(Base):
    __tablename__ = 'group'

    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String(40), unique=True)


class User(Base):
    __tablename__ = 'user'

    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String(40), unique=True)
    group_id = Column(BigInteger)
```
