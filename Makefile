# Install by default to $HOME/.local by systemd conventions.
# See https://www.freedesktop.org/software/systemd/man/file-hierarchy.html#Home%20Directory
PREFIX  ?= $(HOME)/.local
BINDIR  ?= $(PREFIX)/bin
# Default version format for snapshot releases.
# Example: 0.1.2-c91f0d3
VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-$(shell git rev-parse --short HEAD)

APP     = gen
LDFLAGS = "-s -w -X main.version=$(VERSION) -buildid="

.PHONY: build install clean test help

build: ## Compile
	@ go build -o $(APP) -trimpath -ldflags=$(LDFLAGS)

install: build ## Install
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(APP) $(BINDIR)

clean: ## Remove generated files
	@ rm -f $(APP)

test: ## Run tests
	@ go test ./...

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
