#!/bin/bash

set -e -u -x

DIR="$(dirname "$0")"

FILE_BASE="$DIR/../maps_backend"

"$DIR/uwsgi-stop.sh" || true
sleep 2
"$DIR/uwsgi.sh"

while inotifywait -e modify -r --exclude ".*db.sqlite3.*" "$FILE_BASE"; do
    "$DIR/uwsgi-stop.sh" || true
    sleep 2
    "$DIR/uwsgi.sh"
done

