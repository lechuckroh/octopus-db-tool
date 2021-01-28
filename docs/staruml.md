# Star UML

* StarUML [Homepage](https://staruml.io/)

## Import

```bash
$ oct import staruml --help
```

```
OPTIONS:
   --input FILE, -i FILE   import input starUML from FILE [$OCTOPUS_INPUT]
   --output FILE, -o FILE  write octopus schema to FILE [$OCTOPUS_OUTPUT]
```

Import `*uml` file:

```bash
$ oct import staruml --input user.uml --output user.json 
```
