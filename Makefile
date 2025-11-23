.PHONY: help build run test clean install-deps migrate

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install Go dependencies"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"

install-deps:
	go mod download
	go mod tidy

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

dev:
	air
