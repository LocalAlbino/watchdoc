all: fmt build

build:
	go build -o ./bin/watchdoc.exe

fmt:
	go fmt

clean:
	rm ./watchdoc.json

test:
	go test ./cmd/... -v
