PREFIX  ?= $(HOME)/.local
VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-next

BIN     = gen
BINDIR  = $(PREFIX)/bin

DISTDIR = ./dist

LDFLAGS = "-s -w -X go.astrophena.name/gen/internal/version.Version=$(VERSION) -buildid="

.PHONY: build generate install clean test dist help

build: ## Build
	@ go build -o $(BIN) -trimpath -ldflags=$(LDFLAGS)

generate: ## Generate
	@ go generate ./...

install: build ## Install
	@ mkdir -m755 -p $(BINDIR) && \
		install -m755 $(BIN) $(BINDIR)

clean: ## Clean
	@ rm -rf $(BIN) $(DISTDIR)

test: ## Run tests
	@ go test ./...

dist: ## Build with GoReleaser
	@ goreleaser --snapshot --skip-publish

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
