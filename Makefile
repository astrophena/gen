# See https://www.freedesktop.org/software/systemd/man/file-hierarchy.html#Home%20Directory
PREFIX  ?= $(HOME)/.local
BINDIR  ?= $(PREFIX)/bin
VERSION ?= $(shell git rev-parse --short HEAD)

APP     = gen
LDFLAGS = "-s -w -X main.Version=$(VERSION)"

.PHONY: build install clean fmt help

build: ## Compile.
	@ go build -o $(APP) -ldflags=$(LDFLAGS)

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
