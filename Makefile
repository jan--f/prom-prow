.PHONY: build test test-coverage lint clean

build:
	go build -o bin/bot ./cmd/bot

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.out coverage.html

docker-build:
	docker build -t prometheus-prow-bot:latest .

docker-run:
	docker run --rm prometheus-prow-bot:latest

.DEFAULT_GOAL := build
