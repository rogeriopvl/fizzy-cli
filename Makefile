.PHONY: help dev-tools run build build-all build-dev clean test install sync-api-spec

VERSION := $(shell grep '"version"' package.json | sed 's/.*"version": "\([^"]*\)".*/\1/')
LDFLAGS := -ldflags="-s -w -X 'github.com/rogeriopvl/fizzy-cli/cmd.Version=$(VERSION)'"
LDFLAGS_DEBUG := -ldflags="-X 'github.com/rogeriopvl/fizzy-cli/cmd.Version=$(VERSION)'"

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
	GOBIN=$$(go env GOBIN); \
	if [ -z "$$GOBIN" ]; then GOBIN=$$(go env GOPATH)/bin; fi; \
	go build $(LDFLAGS) -o $$GOBIN/fizzy .

build:
	go build $(LDFLAGS) -o bin/fizzy .

build-dev:
	go build $(LDFLAGS_DEBUG) -o bin/fizzy .

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

sync-api-spec:
	@rm -rf docs/api
	@mkdir -p docs
	@curl -sL https://github.com/basecamp/fizzy/archive/refs/heads/main.tar.gz | tar -xz -C docs --strip-components=2 fizzy-main/docs/api
	@echo "API spec synced to docs/api"
