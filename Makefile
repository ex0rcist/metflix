API_DOCS = docs/api

AGENT_VERSION ?= 0.1.0
SERVER_VERSION ?= 0.1.0

BUILD_DATE ?= $(shell date +%F\ %H:%M:%S)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)

help: ## display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

build: agent server staticlint
.PHONY: build

agent: ## build agent
	go build \
		-ldflags "\
			-X 'main.buildVersion=$(AGENT_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go
.PHONY: agent

server: ## build server
	rm -rf $(API_DOCS)
	swag init -g ./internal/httpserver/router.go --output docs/api

	go build \
		-ldflags "\
			-X 'main.buildVersion=$(SERVER_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go
.PHONY: server

staticlint: ## build static lint
	go build -o cmd/$@/$@ cmd/$@/*.go
.PHONY: staticlint

clean: ## remove build artifacts
	rm -rf cmd/agent/agent cmd/server/server cmd/staticlint/staticlint
.PHONY: clean

unit-tests: ## run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
.PHONY: unit-tests

godoc: ### show public packages documentation using godoc
	@echo "Project documentation is available at:"
	@echo "http://127.0.0.1:3000/pkg/github.com/ex0rcist/metflix/pkg/\n"
	@godoc -http=:3000 -play
.PHONY: godoc