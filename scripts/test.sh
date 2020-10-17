#!/usr/bin/env bash

set -e
echo "mode: atomic" > coverage-all.out

for d in $(go list ./... | grep -v -e internal/mock -e cmd/migration); do
    go test -race -p=1 -cover -covermode=atomic -coverprofile=coverage.out "$d"
    tail -n +2 coverage.out >> coverage-all.out
done
