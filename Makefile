SHELL := /bin/bash

.DEFAULT_GOAL := help

.PHONY: help
# Show available commands with descriptions
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: init
# Initialize project — install dependencies
init:
	@go install github.com/matryer/moq@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: vendor
# Install Golang dependencies to /vendor
vendor:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: build
# Build cmd/qu with Go to ./bin
build:
	@CGO_ENABLED=0 go build -o ./bin/ ./cmd/...

.PHONY: run
# Runs cmd/qu with Go
run:
	@go run ./cmd/main/...

.PHONY: test
# Runs tests in parallel without cache
test:
	@go clean -testcache && go test -race -parallel 10 ./...

.PHONY: fix
# Run linter with --fix option for Golang files
fix:
	@GOGC=95 golangci-lint run --fix --verbose

.PHONY: lint
# Run linter fo Golang files
lint:
	@GOGC=95 golangci-lint run --verbose
