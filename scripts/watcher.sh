#!/bin/bash

PID=".pid"

set -e

wait_for_changes() {
    echo 'Waiting for changes...'
    fswatch -1 -e ".*" -i "\\.go$" --recursive ./cmd/ ./config/ ./pkg/ ./internal/
}

start_server() {
    go run cmd/api/main.go & echo $! > $PID
}

reload_server() {
    echo 'Reloading server...'
    kill -INT $(pgrep -P `cat $PID`) || true
    start_server
}



if [[ $OSTYPE == 'darwin'* ]]; then
    # Exit on ctrl-c (without this, ctrl-c would go to fswatch, causing it to
    # reload instead of exit):
    trap 'exit 0' SIGINT

    start_server

    while true; do
        wait_for_changes
        reload_server
    done
else
    start_server
fi
