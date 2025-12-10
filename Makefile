.PHONY: help dev-tools run build clean

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:' Makefile | sed 's/:.*//g' | sed 's/^/  /'

dev-tools:
	@echo "Installing development tools..."
	go install github.com/spf13/cobra@latest

run:
	go run .

build:
	go build -o bin/fizzy .

clean:
	rm -f bin/fizzy

test:
	go test -v ./...
