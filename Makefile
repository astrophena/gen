GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LDFLAGS = "-s -w -X go.astrophena.name/gen/version.Version=$(VERSION) -buildid="

VERSION ?= $(shell git describe --abbrev=0 --tags | cut -c 2-)-next

DISTDIR = ./dist

.PHONY: build clean dist test help

build: ## Build
	@ GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -trimpath -ldflags=$(LDFLAGS)

clean: ## Clean
	@ git clean -d -f -x

dist: ## Build with GoReleaser
	@ goreleaser --snapshot --skip-publish --rm-dist

test: ## Test
	@ go test ./...

help: ## Show help
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;32m%-30s\033[0m %s\n", $$1, $$2}'
