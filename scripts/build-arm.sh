#!/usr/bin/env bash

set -e
now=$(date +'%Y-%m-%dT%T%z')
version=$(git rev-parse --short HEAD)
package="go-sample/pkg/server"

go build -a -ldflags "-X $package.version=$version -X $package.buildTime=$now" -o bootstrap cmd/api/main.go
