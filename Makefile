.PHONY: help dev-tools run build build-all clean test install

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:' Makefile | sed 's/:.*//g' | sed 's/^/  /'

dev-tools:
	@echo "Installing development tools..."
	go install gotest.tools/gotestsum@latest

run:
	go run .

install:
	go install .

build:
	go build -o bin/fizzy .

build-all: clean
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/fizzy-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/fizzy-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin/fizzy-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/fizzy-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/fizzy-windows-amd64.exe .
	@echo "Binaries built successfully in bin/"

clean:
	rm -rf bin/

test:
	gotestsum -- -v ./...
