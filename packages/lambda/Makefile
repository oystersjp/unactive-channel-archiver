archiver-lambda: main.go ../archiver/*.go
	go build .

exec: archiver-lambda
	./archiver-lambda
.PHONY: exec

run: main.go ../archiver/*.go
	go run .
.PHONY: run
