#!/bin/bash

if [ $# -lt 1 ]; then
  echo "Please specify an argument."
  exit 1
fi

PROJECT_NANE=$1

GOOS=linux GOARCH=amd64 go build -o ${PROJECT_NANE} && \
zip ${PROJECT_NANE}.zip ${PROJECT_NANE} client_secret.json && \
aws lambda update-function-code --function-name ${PROJECT_NANE}-reply --zip-file fileb://./${PROJECT_NANE}.zip --no-dry-run