#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o kawakatsu-line && \
zip kawakatsu-line.zip kawakatsu-line client_secret.json && \
aws lambda update-function-code --function-name kawakatsu-line-reply --zip-file fileb://./kawakatsu-line.zip --no-dry-run