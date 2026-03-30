# lyt Makefile
# Usage: make install (requires sudo for system-wide install)

.PHONY: all build install clean test

# Binary name
BINARY=lyt

# Installation paths
PREFIX?=/usr/local
BIN_DIR=$(PREFIX)/bin
MAN_DIR=$(PREFIX)/share/man

# Detect OS for installation
UNAME_S:=$(shell uname -s)

all: build

build:
	go build -o $(BINARY) .

# Install lyt system-wide
# Usage: sudo make install
install: build
	@echo "Installing lyt to $(BIN_DIR)..."
ifeq ($(UNAME_S),Darwin)
	install -d $(BIN_DIR)
	install $(BINARY) $(BIN_DIR)/$(BINARY)
endif
ifeq ($(UNAME_S),Linux)
	install -d $(BIN_DIR)
	install -m 755 $(BINARY) $(BIN_DIR)/$(BINARY)
endif
	@echo "lyt installed to $(BIN_DIR)/$(BINARY)"
	@echo "Run '$(BINARY) --help' to get started"

# Uninstall lyt
uninstall:
	@echo "Removing lyt from $(BIN_DIR)..."
	rm -f $(BIN_DIR)/$(BINARY)
	@echo "lyt removed"

# Development: run tests
test:
	go test ./...

# Development: run lyt from source
run:
	go run .

# Development: build and preview
serve: build
	go run . serve -w
