# Quick DBD

## Export

```bash
$ oct export quickdbd --help
```

```
OPTIONS:
   --input FILE, -i FILE   read octopus schema from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  export quickdbd to FILE [$OCTOPUS_OUTPUT]
 ```

Export to quickdbd file:

```bash
$ oct export quickdbd \
    --input database.json \
    --output quickdbd.txt
```
