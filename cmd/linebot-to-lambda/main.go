package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kawakattsun/linebot2lambda"
)

var config *linebot2lambda.Config

func lambdaHandler(events linebot2lambda.Webhook) error {
	return linebot2lambda.HandleRequest(config, events)
}

func main() {
	c, err := linebot2lambda.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	config = c
	lambda.Start(lambdaHandler)
}
