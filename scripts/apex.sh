#!/bin/bash

set -e

ENV="${1:-dev}";
if [ ! -f "./project.$ENV.json" ]; then
  echo "ERROR: Invalid environment '$ENV'! File './project.$ENV.json' is not found." >&2
  exit 1
fi
if [[ $# > 0 ]]; then shift; fi

if [[ "$ENV" == "dev" ]]; then
  ALIAS="dev"
else
  ALIAS="staging"
fi

ARGS=()
while [[ $# -gt 0 ]]
do
  key="$1"
  ARGS+=("$1")

  case $key in
    deploy)
    DEPLOYING=true
    shift
    ;;
    -a|--alias)
    ALIAS="$2"
    ARGS+=("$2")
    shift
    shift
    ;;
    *)
    shift
    ;;
  esac
done

if [[ "$DEPLOYING" == "true" ]]; then
  COMMIT="$(git rev-parse --short HEAD)"
  ARGS+=("--set COMMIT=$COMMIT")
  ARGS+=("--set CONFIG_STAGE=$ALIAS")
fi

echo "$ apex -e $ENV ${ARGS[@]}"
# echo "DEPLOYING: $DEPLOYING"
# echo "ALIAS: $ALIAS"
echo ""
apex -e $ENV ${ARGS[@]}

if [[ "$DEPLOYING" == "true" ]]; then
  echo ""
  echo "$ apex alias -e $ENV -v $ALIAS commit-$COMMIT"
  echo ""
  apex alias -e $ENV -v $ALIAS commit-$COMMIT
fi
