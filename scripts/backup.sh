#!/bin/bash
TIMESTAMP=$(date +%F-%H%M)
DUMP_PATH="/tmp/mongobackup-$TIMESTAMP"

/usr/bin/mongodump --uri="<host 1>" --out="$DUMP_PATH"
/usr/bin/mongorestore --drop --uri="<host 2>" "$DUMP_PATH"
rm -rf "$DUMP_PATH"
