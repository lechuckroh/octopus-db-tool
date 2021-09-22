# JPA

[한국어](kr/jpa.md)

## Generate Kotlin

```shell
$ oct generate kt --help
```

|           Option           |        Env. Variable         | Description                                                                                         |
| :------------------------: | :--------------------------: | :-------------------------------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | Octopus schema file to read                                                                         |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | Target file or directory                                                                            |
|    `-a`, `--annotation`    |     `OCTOPUS_ANNOTATION`     | Custom Entity class annotation.<br />Format: `{group1}:{annotations1}[,{group2}:{annotations2}]...` |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | Table groups to generate.<br />Set multiple groups with comma(`,`) separated.                       |
|     `-e`, `--idEntity`     |     `OCTOPUS_ID_ENTITY`      | Interface NAME with `id` field                                                                      |
|     `-p`, `--package`      |      `OCTOPUS_PACKAGE`       | Entity class package name                                                                           |
|      `-f`, `--prefix`      |       `OCTOPUS_PREFIX`       | Class name prefix.<br />Format: `{group1}:{prefix1}[,{group2}:{prefix2}]...`                        |
|     `-l`, `--relation`     |      `OCTOPUS_RELATION`      | Virtual relation annotation type.<br>Available values: `VRelation`                                  |
|   `-d`, `--removePrefix`   |   `OCTOPUS_REMOVE_PREFIX`    | Prefixes to remove from class name.<br />Set multiple prefixes with comma(`,`) separated.           |
|   `-r`, `--reposPackage`   |   `OCTOPUS_REPOS_PACKAGE`    | Repository class package name. Skip if not set.                                                     |
| `-q`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | Unique constraint name suffix                                                                       |
|      `-u`, `--useUTC`      |      `OCTOPUS_USE_UTC`       | Set flag to use UTC for audit columns (`created_at`, `updated_at`).<br />Default: `false`           |

### Example

```shell
$ oct generate kt \
    --input examples/user.json \
    --output output/ \
    --package octopus.entity \
    --reposPackage octopus.repos \
    --useUTC
```

Generated `*.kt` files:

#### `UserGroup.kt`

```kotlin
package octopus.entity

import javax.persistence.*

@Entity
@Table(name="group", uniqueConstraints = [
    UniqueConstraint(name = "group", columnNames = ["name"])
])
data class UserGroup(
        @Id
        @GeneratedValue(strategy = GenerationType.AUTO)
        @Column(nullable = false)
        var id: Long = 0L,

        @Column(nullable = false, length = 40)
        var name: String = ""
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
        @Column(nullable = false)
        var id: Long = 0L,

        @Column(nullable = false, length = 40)
        var name: String = "",

        var groupId: Long?
)
```

#### `UserGroupRepository.kt`

```kotlin
package octopus.repos

import octopus.entity.*
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface UserGroupRepository : JpaRepository<UserGroup, Long>
```

#### `UserRepository.kt`

```kotlin
package octopus.repos

import octopus.entity.*
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface UserRepository : JpaRepository<User, Long>
```
