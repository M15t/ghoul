#!/bin/bash

set -e

ENV="${1:-dev}";
if [[ $# > 0 ]]; then shift; fi

ARGS="$@"
UPFILE="./up.${ENV}.json"

if [ ! -f "$UPFILE" ]; then
  echo "ERROR: File '$UPFILE' not found!" >&2
  exit 1
fi

yes | cp -f $UPFILE up.json
trap 'rm -f up.json' EXIT

echo "$ up ${ARGS}"

up "$@"
