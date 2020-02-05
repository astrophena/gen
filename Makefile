# Install binary by default to $HOME/.local/bin by systemd conventions.
# See https://www.freedesktop.org/software/systemd/man/file-hierarchy.html#Home%20Directory
PREFIX  ?= $(HOME)/.local
BINDIR  ?= $(PREFIX)/bin
# Default version format for snapshot releases. Example: 0.1.2-c91f0d3
VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-$(shell git rev-parse --short HEAD)

APP     = gen
LDFLAGS = "-s -w -X main.Version=$(VERSION) -buildid="

.PHONY: build install clean fmt help

build: ## Compile.
	@ go build -o $(APP) -trimpath -ldflags=$(LDFLAGS)

install: build ## Install to $PREFIX.
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(APP) $(BINDIR)

clean: ## Remove all generated files.
	@ rm -f $(APP)

fmt: ## Reformat sources.
	@ goimports -w .

help: ## Display this help.
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
