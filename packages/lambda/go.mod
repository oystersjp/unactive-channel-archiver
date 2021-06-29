module archiver-lambda

go 1.16

require (
	archiver/archiver v0.0.0
	github.com/aws/aws-lambda-go v1.24.0
)

replace archiver/archiver => ../archiver
