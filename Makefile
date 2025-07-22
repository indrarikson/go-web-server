# Makefile for Go Web Server Template

.PHONY: all generate build test vet clean tidy run dev deps security lint help

# Build configuration
BINARY_NAME=server
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/web

# Default target
all: build

## generate: Generate code (templ, sqlc)
generate:
	@echo ">> generating code"
	@go generate ./...

## build: Build the application
build: generate
	@echo ">> building binary"
	@mkdir -p bin
	@go build -ldflags="-s -w" -o $(BINARY_PATH) $(MAIN_PATH)

## test: Run tests
test:
	@echo ">> running tests"
	@go test -race -cover ./...

## vet: Run go vet
vet:
	@echo ">> vetting code"
	@go vet ./...

## tidy: Format and tidy code
tidy:
	@echo ">> tidying and formatting"
	@go mod tidy
	@go fmt ./...

## run: Build and run application
run: build
	@echo ">> starting server"
	@$(BINARY_PATH)

## dev: Development with hot reload (requires air)
dev:
	@if ! command -v air > /dev/null 2>&1; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

## clean: Clean build artifacts
clean:
	@echo ">> cleaning"
	@rm -rf bin/ tmp/

## deps: Install development tools
deps:
	@echo ">> installing tools"
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest

## security: Run security scan
security:
	@echo ">> running security scan"
	@if command -v govulncheck > /dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed. Run 'make deps' first."; \
		exit 1; \
	fi

## lint: Run all linters
lint: vet security

## help: Show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'