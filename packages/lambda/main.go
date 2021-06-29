package main

import (
	"archiver/archiver"

	"github.com/aws/aws-lambda-go/lambda"
)

func execSay() {
	archiver.Say()
}

func main() {
	lambda.Start(execSay)
}
