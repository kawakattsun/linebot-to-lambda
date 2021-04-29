package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kawakattsun/linebot2lambda"
)

const version = "0.1.0"

var config *linebot2lambda.Config

func lambdaHandler(events linebot2lambda.Webhook) error {
	return linebot2lambda.HandleRequest(config, events)
}

func main() {
	fmt.Printf("Start lambda function. %s %s\n", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"), version)
	c, err := linebot2lambda.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	config = c
	lambda.Start(lambdaHandler)
}
