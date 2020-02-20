# Install by default to $HOME/.local
PREFIX  ?= $(HOME)/.local
# Default version format
# Example: 0.1.2-c91f0d3
VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-$(shell git rev-parse --short HEAD)

BIN     = gen
BINDIR  = $(PREFIX)/bin

LDFLAGS = "-s -w -X main.version=$(VERSION) -buildid="

.PHONY: build install clean test help

build: ## Build
	@ go build -o $(BIN) -trimpath -ldflags=$(LDFLAGS)

install: build ## Install
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(BIN) $(BINDIR)

clean: ## Remove generated files
	@ rm -f $(BIN)

test: ## Run tests
	@ go test ./...

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
