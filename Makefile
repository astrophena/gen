PREFIX  ?= $(HOME)/.local
VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-next

BIN     = gen
BINDIR  = $(PREFIX)/bin

LDFLAGS = "-s -w -X astrophena.me/gen/internal/buildinfo.Version=$(VERSION) -buildid="

.PHONY: build generate install clean test help

build: ## Build
	@ go build -o $(BIN) -trimpath -ldflags=$(LDFLAGS)

generate: ## Generate
	@ go generate ./...

install: build ## Install
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(BIN) $(BINDIR)

clean: ## Clean
	@ rm -f $(BIN)

test: ## Run tests
	@ go test ./...

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
