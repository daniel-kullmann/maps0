#!/bin/bash

set -e -u

DIR="$(dirname "$0")"

FILE_BASE="$DIR/../node_tile_server"

node "$FILE_BASE/server.js" &
PID="$!"
echo "node server started ($PID)"

while inotifywait -q -e modify -r "$FILE_BASE"; do
    echo "changes detected; restarting node server ($PID)"
    kill "$PID"
    sleep 1
    node "$FILE_BASE/server.js" &
    PID="$!"
    echo "new node server started ($PID)"
done

