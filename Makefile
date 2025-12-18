.PHONY: help dev-tools run build build-all clean test install

VERSION := $(shell grep '"version"' package.json | sed 's/.*"version": "\([^"]*\)".*/\1/')
LDFLAGS := -ldflags="-X 'github.com/rogeriopvl/fizzy/cmd.Version=$(VERSION)'"

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
	go install $(LDFLAGS) .

build:
	go build $(LDFLAGS) -o bin/fizzy .

build-all: clean
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/fizzy-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/fizzy-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/fizzy-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/fizzy-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/fizzy-windows-amd64.exe .
	@echo "Binaries built successfully in bin/"

clean:
	rm -rf bin/

test:
	gotestsum -- -v ./...
