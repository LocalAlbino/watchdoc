all: fmt build

build:
	go build -o ./bin/watchdoc.exe

fmt:
	go fmt
