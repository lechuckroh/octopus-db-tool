# ProtoBuf

[한국어](kr/protobuf.md)

## Generate

```shell
$ oct generate pb --help
```

|         Option         |        Env. Variable         | Description                                                                                                                        |
| :--------------------: | :--------------------------: | :--------------------------------------------------------------------------------------------------------------------------------- |
|    `-i`, `--input`     |       `OCTOPUS_INPUT`        | Octopus schema file to read                                                                                                        |
|    `-o`, `--output`    |       `OCTOPUS_OUTPUT`       | Target file                                                                                                                        |
|     `--goPackage`      |     `OCTOPUS_GO_PACKAGE`     | Protobuf golang package name                                                                                                       |
|    `-g`, `--groups`    |       `OCTOPUS_GROUPS`       | Table groups to generate.<br />Set multiple groups with comma(`,`) separated.                                                      |
|   `-p`, `--package`    |      `OCTOPUS_PACKAGE`       | Protobuf package name                                                                                                              |
|    `-f`, `--prefix`    |       `OCTOPUS_PREFIX`       | Proto message name prefix.<br />Format: `<group1>:<prefix1>[,<group2>:<prefix2>]...`<br />Example: `group1:prefix1,group2:prefix2` |
| `-d`, `--removePrefix` |   `OCTOPUS_REMOVE_PREFIX`    | Prefixes to remove from proto message name.<br />Set multiple prefixes with comma(`,`) separated.                                  |
|  `--relationTagDecr`   | `OCTOPUS_RELATION_TAG_DECR`  | Relation tags decremental from relationTagStart. Default: `false`                                                                  |
|  `--relationTagStart`  | `OCTOPUS_RELATION_TAG_START` | Relation tags start index. Set `-1` to start from last of fields.                                                                  |

### Example

```shell
$ oct generate pb \
    --input examples/user.json \
    --output output/user.proto \
    --goPackage model \
    --package octopus
```

Generated `*.proto` file:

```protobuf
syntax = "proto3";

package octopus;

option go_package = "model";

message Group {
  int64 id = 1;
  string name = 2;
}

message User {
  int64 id = 1;
  string name = 2;
  int64 groupId = 3;
  Group group = 4;
}
```
