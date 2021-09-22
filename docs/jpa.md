# JPA

## Generate Kotlin

```shell
$ oct generate kt --help
```

```
OPTIONS:
   --input FILE, -i FILE                 read octopus schema from FILE [$OCTOPUS_INPUT]
   --output DIR, -o DIR                  generate kotlin files to DIR [$OCTOPUS_OUTPUT]
   --annotation FORMAT, -a FORMAT        Custom Entity class annotation. FORMAT: '{group1}:{annotations1}[,{group2}:{annotations2}]' [$OCTOPUS_ANNOTATION]
   --groups GROUPS, -g GROUPS            Filter table groups to generate. GROUPS are separated by comma [$OCTOPUS_GROUPS]
   --idEntity NAME, -e NAME              Interface NAME with 'id' field [$OCTOPUS_ID_ENTITY]
   --package PACKAGE, -p PACKAGE         Entity class PACKAGE name [$OCTOPUS_PACKAGE]
   --prefix FORMAT, -f FORMAT            Class name prefix. FORMAT: '{group1}:{prefix1}[,{group2}:{prefix2}]' [$OCTOPUS_PREFIX]
   --relation ANNOTATION, -l ANNOTATION  Virtual relation ANNOTATION type. Available values: VRelation [$OCTOPUS_RELATION]
   --removePrefix PREFIXES, -d PREFIXES  Table PREFIXES to remove from class name. Multiple prefixes are separated by comma [$OCTOPUS_REMOVE_PREFIX]
   --reposPackage PACKAGE, -r PACKAGE    Repository class PACKAGE name. Generated if not empty. [$OCTOPUS_REPOS_PACKAGE]
   --uniqueNameSuffix SUFFIX, -q SUFFIX  Unique constraint name SUFFIX. [$OCTOPUS_UNIQUE_NAME_SUFFIX]
   --useUTC, -u                          Set to use UTC for audit columns ('created_at', 'updated_at'). (default: false) [$OCTOPUS_USE_UTC]
```

Generate `*.kt` files:

```shell
# example with all CLI options
$ oct generate kt \
    --input database.json \
    --output ./output \
    --annotation foo:@Foo,foobar:@Foo,@Bar \
    --groups foo,bar \
    --idEntity IdEntity \
    --package com.foo.entity \
    --prefix foo:F,bar:B \
    --relation VRelation \
    --removePrefix tbl_,table_ \
    --reposPackage com.foo.repos \
    --uniqueNameSuffix _uq \
    --useUTC    
```

### Example

```shell
$ oct generate kt \
    --input examples/user.json \
    --output output/ \
    --package octopus.entity \
    --reposPackage octopus.repos \
    --useUTC
```

Generated source files:

#### `Group.kt`

```kotlin
package octopus.entity

import javax.persistence.*

@Entity
@Table(name="group", uniqueConstraints = [
    UniqueConstraint(name = "group", columnNames = ["name"])
])
data class Group(
        @Id
        @GeneratedValue(strategy = GenerationType.AUTO)
        var id: Long?,

        @Column(length = 40)
        var name: String?
)
```

#### `User.kt`

```kotlin
package octopus.entity

import javax.persistence.*

@Entity
@Table(name="user", uniqueConstraints = [
    UniqueConstraint(name = "user", columnNames = ["name"])
])
data class User(
        @Id
        @GeneratedValue(strategy = GenerationType.AUTO)
        var id: Long?,

        @Column(length = 40)
        var name: String?,

        var groupId: Long?
)
```

#### `GroupRepository.kt`

```kotlin
package octopus.repos

import octopus.entity.*
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface GroupRepository : JpaRepository<Group, Long?>
```

#### `UserRepository.kt`

```kotlin
package octopus.repos

import octopus.entity.*
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface UserRepository : JpaRepository<User, Long?>
```
