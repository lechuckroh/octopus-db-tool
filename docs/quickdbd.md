# Quick DBD

## Export

```shell
$ oct export quickdbd --help
```

```
OPTIONS:
   --input FILE, -i FILE   read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  export quickdbd to FILE [$OCTOPUS_OUTPUT]
 ```

Export to quickdbd file:

```shell
# example with all CLI options
$ oct export quickdbd \
    --input database.json \
    --output quickdbd.txt
```

### Example

```shell
$ oct export quickdbd \
    --input examples/user.json \
    --output output/user.txt
```

Exported file:

```
group # Group table
-----
id int64 PK AUTOINCREMENT NULLABLE
name varchar UNIQUE NULLABLE

user # User table
----
id int64 PK AUTOINCREMENT NULLABLE
name varchar UNIQUE NULLABLE
group_id int64 NULLABLE FK >- group.id
```
