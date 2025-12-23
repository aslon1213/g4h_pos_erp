#!/bin/bash
TIMESTAMP=$(date +%F-%H%M)
DUMP_PATH="/tmp/mongobackup-$TIMESTAMP"

/usr/bin/mongodump --uri="<source_host>" --out="$DUMP_PATH"
/usr/bin/mongorestore --drop --uri="<target_host>"$DUMP_PATH"
rm -rf "$DUMP_PATH"
