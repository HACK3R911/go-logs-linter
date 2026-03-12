.PHONY: build test lint clean

run:
	go run ./cmd/loglint ./testdata/bad/. &
	go run ./cmd/loglint ./testdata/good/.

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
