#!/usr/bin/env bash

set -ex

# Download SwaggerUI if not exist
if [ ! -d "./swaggerui" ]; then
    # curl -sL -o swagger-ui-3.19.4.zip https://github.com/swagger-api/swagger-ui/archive/v3.19.4.zip
    # unzip -qq -j swagger-ui-3.19.4.zip "swagger-ui-3.19.4/dist/*" -d ./swaggerui/
    # rm -f swagger-ui-3.19.4.zip
    git clone git@github.com:M15t/swagger-ui.git ./swaggerui
fi

# Generate swagger.json file
swagger generate spec --scan-models -w ./cmd/api/ -o ./swaggerui/swagger.json

if [[ "$AWS_LAMBDA_FUNCTION_NAME" != "" && "$STAGE" != "" && "$STAGE" != "development" ]]; then
    HOST=$(aws ssm get-parameters --name "/$AWS_LAMBDA_FUNCTION_NAME/$STAGE/host" --with-decryption | jq -r '.Parameters[0].Value')
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    # Replace HOST by corresponding env var
    sed -i '' -e "s#%{HOST}#$HOST#g" ./swaggerui/swagger.json
    # Replace default URL with latest commit ID to avoid browser caching
    sed -i '' -e "s|url:.*|url: \"./swagger.json?$(git rev-parse --short HEAD)\",|" ./swaggerui/index.html
else
    # Replace HOST by corresponding env var
    sed -i -e "s#%{HOST}#$HOST#g" ./swaggerui/swagger.json
    # Replace default URL with latest commit ID to avoid browser caching
    sed -i -e "s|url:.*|url: \"./swagger.json?$(git rev-parse --short HEAD)\",|" ./swaggerui/index.html
fi


