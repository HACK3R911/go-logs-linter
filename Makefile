.PHONY: build test lint clean

build:
	go build -o bin/loglint ./cmd/loglint

test:
	go test -v ./...

lint:
	golangci-lint run ./...

custom:
	golangci-lint custom

clean:
	rm -rf bin/
	go clean -cache
