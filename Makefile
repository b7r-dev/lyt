# lyt - Makefile for development and CI

# Default target
.PHONY: all
all: lint test build

# Build the lyt binary
.PHONY: build
build:
	cd engine && go build -o lyt

# Run tests with coverage
.PHONY: test
test:
	cd engine && go test -v -race -coverprofile=coverage.out ./...

# Run tests without coverage (faster for CI)
.PHONY: test-ci
test-ci:
	cd engine && go test -v -race ./...

# Lint the code
.PHONY: lint
lint:
	cd engine && golangci-lint run ./... || (go vet ./... && go fmt ./...)

# Lint with auto-fix when possible
.PHONY: lint-fix
lint-fix:
	cd engine && golangci-lint run --fix ./... || true

# Format code
.PHONY: fmt
fmt:
	cd engine && go fmt ./... && go vet ./...

# Show coverage report
.PHONY: coverage
coverage:
	cd engine && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at engine/coverage.html"

# Run the full CI pipeline locally
.PHONY: ci
ci: lint test

# Watch mode (requires fswatch or entr)
.PHONY: watch
watch:
	@echo "Install fswatch or entr to use watch mode"
	@echo "  brew install fswatch"
	@echo "Then run: find engine -name '*.go' | entr -c make test"

# Clean build artifacts
.PHONY: clean
clean:
	rm -f engine/lyt engine/coverage.out engine/coverage.html
	rm -rf engine/dist

# Install dependencies
.PHONY: deps
deps:
	cd engine && go mod download && go mod tidy

# Install development tools
.PHONY: setup
setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run lyt build
.PHONY: gen
gen:
	cd engine && ./lyt build

# Run lyt serve
.PHONY: serve
serve:
	cd engine && ./lyt serve

# Build and serve (default development workflow)
.PHONY: dev
dev: build gen serve
