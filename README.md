# shift
Shift data between databases

### Workflow

* Read config file
* Create `_shift_state` table in target database which will hold table names and offset positions for all synced tables
* Read `read_limit` rows of data
* Write last_id written into `_shift_state`
* Continue until we've read everything in the source table