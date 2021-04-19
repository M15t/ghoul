#!/usr/bin/env bash

echo "Please enter function name:"

read funcName

for a in $(aws lambda list-aliases --function-name $funcName | jq -r '.Aliases[].Name'); do
    if [[ $a != "dev" && $a != "demo" && $a != "staging" && $a != "production" ]]; then
        aws lambda delete-alias --function-name $funcName --name $a
        echo "Deleted" $a
    fi
done
