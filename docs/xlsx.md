# xlsx

## Import

```bash
$ oct import xlsx --help
```

```
OPTIONS:
   --input FILE, -i FILE   import xlsx from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  write octopus schema to FILE [$OCTOPUS_OUTPUT]
```

Import `*.xlsx` file:

```bash
$ oct import xlsx \
    -i user.xlsx \
    -o user.json
```

## Export

```bash
$ oct export xlsx --help
```

```
OPTIONS:
   --input FILE, -i FILE               read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE              export xlsx to FILE [$OCTOPUS_OUTPUT]
   --useNotNullColumn value, -n value  use not null column [$OCTOPUS_USE_NOT_NULL_COLUMN]
```

Export `*.xlsx` file:

```bash
$ oct export xlsx \
    -i user.json \
    -o user.xlsx
```
