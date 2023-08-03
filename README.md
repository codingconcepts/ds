# dshift
Shift data between databases

### Installation

Find the release that matches your architecture on the [releases](https://github.com/codingconcepts/shift/releases) page.

Download the tar, extract the executable, and move it into your PATH:

```
$ tar -xvf dshift_[VERSION]_macOS.tar.gz
```

### Usage

```
$ dshift

Shift data from one from database to another

Usage:
  dshift [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  insert      Insert data from one database into another
  update      Bring the target database up-to-date with the source database
  version     Print dshift version information

Flags:
  -c, --config string   absolute or relative path to the config file
  -h, --help            help for dshift

Use "dshift [command] --help" for more information about a command.
```

### Example

Create source database:
``` sh
make postgres
make postgres_create
make postgres_insert
make postgres_wal_level
```

Create target database:
``` sh
make cockroach
make cockroach_create
```

Shift
```sh
make dshift
```

Verify
``` sh
make verify
```

### Workflow

* Read config file
* Create `_shift_state` table in target database which will hold table names and offset positions for all synced tables
* Read `read_limit` rows of data
* Write last_id written into `_shift_state`
* Continue until we've read everything in the source table

### Todo

* Create a hidden `_shift_digest` column and populate it with the digest of the row at time of insert:

``` sql
SELECT sha256(e::TEXT) FROM example e;

ALTER TABLE example
  ADD COLUMN _shift_digest STRING
  CREATE IF NOT EXISTS FAMILY _shift
  NOT VISIBLE;
```

* Add the following cobra commands:
  * **update** - bulk upload changed table data
  * **delete** - bulk upload deleted table data