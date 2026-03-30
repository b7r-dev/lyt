# lyt - Makefile for development and CI

# Default target
.PHONY: all
all: lint test build

# Build the lyt binary
.PHONY: build
build:
	go build -o lyt

# Run tests with coverage
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests without coverage (faster for CI)
.PHONY: test-ci
test-ci:
	go test -v -race ./...

# Lint the code
.PHONY: lint
lint:
	golangci-lint run ./... || (go vet ./... && go fmt ./...)

# Lint with auto-fix when possible
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./... || true

# Format code
.PHONY: fmt
fmt:
	go fmt ./... && go vet ./...

# Show coverage report
.PHONY: coverage
coverage:
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run the full CI pipeline locally
.PHONY: ci
ci: lint test

# Watch mode (requires fswatch or entr)
.PHONY: watch
watch:
	@echo "Install fswatch or entr to use watch mode"
	@echo "  brew install fswatch"
	@echo "Then run: find . -name '*.go' | entr -c make test"

# Clean build artifacts
.PHONY: clean
clean:
	rm -f lyt coverage.out coverage.html
	rm -rf dist

# Install dependencies
.PHONY: deps
deps:
	go mod download && go mod tidy

# Install development tools
.PHONY: setup
setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run lyt build
.PHONY: gen
gen:
	./lyt build

# Run lyt serve
.PHONY: serve
serve:
	./lyt serve

# Build and serve (default development workflow)
.PHONY: dev
dev: build gen serve
