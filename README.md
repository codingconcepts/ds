<p align="center">
  <img src="assets/cover.png" alt="drawing" width="350"/>
</p>

Shift data between databases

### Installation

Find the release that matches your architecture on the [releases](https://github.com/codingconcepts/shift/releases) page.

Download the tar, extract the executable, and move it into your PATH:

```
$ tar -xvf ds_[VERSION]_macOS.tar.gz
```

### Usage

```
$ ds

Shift data from one from database to another

Usage:
  ds [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  insert      Insert data from one database into another
  update      Bring the target database up-to-date with the source database
  version     Print ds version information

Flags:
  -c, --config string   absolute or relative path to the config file
  -h, --help            help for ds

Use "ds [command] --help" for more information about a command.
```

### Example

Create source database:
``` sh
make postgres
make postgres_create
make postgres_insert
```

Create target database:
``` sh
make cockroach
make cockroach_create
```

Show that data is currently out-of-sync:
``` sh
make verify
```

Bulk insert all rows from the source database and re-run verify:
```sh
ds insert --config examples/basic/config.yaml

make verify
```

Update data in the source database and re-run verify:
``` sh
make postgres_update

make verify
```

Bulk update all rows and re-run verify:
```sh
ds update --config examples/basic/config.yaml

make verify
```

### Todos

* Implement `delete` command to clear up missing rows