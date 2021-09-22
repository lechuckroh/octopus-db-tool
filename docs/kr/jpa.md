# JPA

[English](../jpa.md)

## 코틀린 소스 생성

```shell
$ oct generate kt --help
```

|            옵션            |           환경변수           | 설명                                                                                           |
| :------------------------: | :--------------------------: | :--------------------------------------------------------------------------------------------- |
|      `-i`, `--input`       |       `OCTOPUS_INPUT`        | 입력으로 사용할 octopus 스키마 파일명                                                          |
|      `-o`, `--output`      |       `OCTOPUS_OUTPUT`       | 출력할 디렉토리명                                                                              |
|    `-a`, `--annotation`    |     `OCTOPUS_ANNOTATION`     | 커스텀 엔티티 클래스 애노테이션.<br />형식: `<그룹1>:<애노테이션1>[,<그룹2>:<애노테이션2>]...` |
|      `-g`, `--groups`      |       `OCTOPUS_GROUPS`       | 생성할 대상 테이블 그룹명.<br />여러개의 그룹을 지정시 `,`로 구분                              |
|     `-e`, `--idEntity`     |     `OCTOPUS_ID_ENTITY`      | `id` 필드만 가진 인터페이스 명                                                                 |
|     `-p`, `--package`      |      `OCTOPUS_PACKAGE`       | 엔티티 클래스 패키지 명                                                                        |
|      `-f`, `--prefix`      |       `OCTOPUS_PREFIX`       | 클래스 이름 접두사.<br />형식: `<그룹1>:<접두사1>[,<그룹2>:<접두사2>]...`                     |
|     `-l`, `--relation`     |      `OCTOPUS_RELATION`      | 가상 연관관계를 나타내는 애노테이션 타입.<br>사용 가능한 값: `VRelation`                       |
|   `-d`, `--removePrefix`   |   `OCTOPUS_REMOVE_PREFIX`    | 생성할 클래스명에서 제거할 접두사.<br />여러개의 접두사 지정시 `,`로 구분                      |
|   `-r`, `--reposPackage`   |   `OCTOPUS_REPOS_PACKAGE`    | Repository 클래스 패키지 명. 지정한 경우에만 생성                                              |
| `-q`, `--uniqueNameSuffix` | `OCTOPUS_UNIQUE_NAME_SUFFIX` | 유니크 제약 이름 접미사                                                                        |
|      `-u`, `--useUTC`      |      `OCTOPUS_USE_UTC`       | 플래그 설정시 audit 컬럼(`created_at`, `updated_at`)들에 대해 UTC 사용.<br />기본값: `false`   |

### 예제

```shell
$ oct generate kt \
    --input examples/user.json \
    --output output/ \
    --package octopus.entity \
    --reposPackage octopus.repos \
    --useUTC
```

`*.kt` 파일은 다음과 같이 생성됩니다:

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
