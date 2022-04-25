---
date: 2022-04-25T16:27:37+02:00
title: "backpack restore"
slug: backpack_restore
url: /commands/backpack_restore/
---
## backpack restore

Restore from uploaded file

```
backpack restore [flags]
```

### Options

```
  -a, --action string        Name of action to restore from.
  -c, --config string        Path to config file.
      --except stringArray   List of connections to ignore.
  -f, --file string          Name of file to restore from
      --force                Force backup even if prev_hash is the same
  -h, --help                 help for restore
      --no-backup            Doesn't create backup
      --no-encrypt           Doesn't encrypt files
      --only stringArray     List of connections to try.
```

### Options inherited from parent commands

```
      --debug     Enable debug mode. MAY PRINT SENSITIVE INFORMATION
      --verbose   Print more information.
```

### SEE ALSO

* [backpack](/commands/backpack/)	 - Easily backup and restore folders to and from different storages

