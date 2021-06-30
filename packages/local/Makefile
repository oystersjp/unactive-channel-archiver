archiver-local: main.go ../archiver/*.go
	go mod download
	go build .

exec: archiver-local
	./archiver-local
.PHONY: exec

run: main.go ../archiver/*.go
	go run .
.PHONY: run
