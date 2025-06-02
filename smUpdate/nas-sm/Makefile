SHELL := /usr/bin/env bash -o pipefail

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin


PROJECT    := wipro5gcore
#VERSION           ?= $(shell git describe --always --tags --dirty --match 'v*')
#COMMIT     ?= $(shell git rev-parse HEAD)
#BRANCH     ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE ?= $(shell date +%s)
BUILD_HOST ?= $(shell hostname)
BUILD_USER ?= $(shell id -un)

GRPCNASSMF_DOCKERFILE_DIR := docker
GRPCNASSMF_IMAGE_TAG := grpcnassmf
GRPCNASSMF_BIN_NAME := grpcnassmf
# PDUSMSP_DEV_DOCKERFILE_DIR := docker/dev/smf/pdusmsp
# PDUSMSP_PROD_DOCKERFILE_DIR := docker/prod/smf/pdusmsp
# PDUSMSP_DEV_IMAGE_TAG := pdusmsp_dev
# PDUSMSP_PROD_IMAGE_TAG := pdusmsp_prod
# CSP_DEV_DOCKERFILE_DIR := docker/dev/amf/csp
# CSP_PROD_DOCKERFILE_DIR := docker/prod/amf/csp
# CSP_DEV_IMAGE_TAG := csp_dev
# CSP_PROD_IMAGE_TAG := csp_prod
# PDUSMSP_BIN_NAME := pdusmsp
# CSP_BIN_NAME := csp
GOPKG := w5gc.io/wipro5gcore/v1
LDFLAGS = -w -s \
        -X $(GOPKG)/pkg/version.app=$(PROJECT) \
        -X $(GOPKG)/pkg/version.buildDate=$(BUILD_DATE) \
        -X $(GOPKG)/pkg/version.buildUser=$(BUILD_USER) \
        -X $(GOPKG)/pkg/version.buildHost=$(BUILD_HOST)
#       -X $(GOPKG)/pkg/version.version=$(VERSION) \
#       -X $(GOPKG)/pkg/version.gitCommit=$(COMMIT) \
#       -X $(GOPKG)/pkg/version.gitBranch=$(BRANCH) \

UNAME_OS   ?= $(shell uname -s)
UNAME_ARCH ?= $(shell uname -m)

ifndef CACHE_BASE
CACHE_BASE := $(HOME)/.cache/$(PROJECT)
endif
CACHE := $(CACHE_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
CACHE_BIN := $(CACHE)/bin
CACHE_INCLUDE := $(CACHE)/include
CACHE_VERSIONS := $(CACHE)/versions

export PATH := $(abspath $(CACHE_BIN)):$(PATH)

ifndef BUILD_DIR
BUILD_DIR := .build
endif

export GO111MODULE=on
export DOCKER_BUILDKIT=1

ifeq ($(NOSTRIP),)
LDFLAGS += -w -s
endif

ifeq ($(NOTRIM),)
GO_BUILD_ARGS += -trimpath
endif

ifeq ($(BUILDPIE),y)
GO_BUILD_ARGS += -buildmode=pie
LDFLAGS += -extldflags=-Wl,-z,now,-z,relro
endif

ifeq ($(V),1)
GO_BUILD_ARGS += -v
endif

COVER_DIR ?= /tmp

help:
        @echo "List of make targets:"
        @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT = help

-include scripts/make/buf.make

build: cmd

clean: clean-cmd

grpcnassmf: ## Build grpcNASSMF
	@echo "# installing grpcnassmf ${VERSION}"
	@go install -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/grpcNasSmf




install: ## Install commands
	@echo "# installing ${VERSION}"
	GOBIN=${GOBIN} go install -buildvcs=false -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/grpcNasSmf
	
# GOBIN=${GOBIN} go install -buildvcs=false -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/smf/pdusmsp
# GOBIN=${GOBIN} go install -buildvcs=false -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/amf/csp
# GOBIN=${GOBIN} go install -buildvcs=false -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/smf/upfgw
#	go install -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/smf/pdusmsp-init
#        go install -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/smf/upfgw
#        go install -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/smf/upfgw-init
#       go install -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS} ./cmd/oamctl

ginkgo-test:
	@ginkgo -r

code-check:
	@go fmt ./...
	@go vet ./...

cmd: ## Build commands
	@echo "# building ${VERSION}"
	cd cmd/grpcNasSmf && go build -buildvcs=false -ldflags "${LDFLAGS}" -tags="${GO_BUILD_TAGS}" ${GO_BUILD_ARGS}


clean-cmd: ## Clean commands
	@echo "# cleaning command binaries"
	rm -f ./cmd/grpcNasSmf/grpcnassmf
# rm -f ./cmd/smf/pdusmsp/pdusmsp
# rm -f ./cmd/smf/pdusmsp/pdusmsp-init
# rm -f ./cmd/amf/csp/csp
# rm -f ./cmd/amf/csp/csp-init
#        rm -f ./cmd/smf/upfgw/upfgw
#        rm -f ./cmd/smf/upfgw/upfgw-init
#       rm -f ./cmd/oamctl/oamctl

purge: ## Purge cached files
	go clean -testcache -cache ./...

debug-remote: ## Debug remotely
    cd ./cmd/grpcNasSmf && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient
#cd ./cmd/smf/pdusmsp && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient
#cd ./cmd/amf/csp && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient
#        cd ./cmd/smf/upfgw && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient

# -------------------------------
#  Testing
# -------------------------------

test: ## Run unit tests
	@echo "# running unit tests"
	go test -tags="${GO_BUILD_TAGS}" ./...

test-cover: ## Run unit tests with coverage
	@echo "# running unit tests with coverage"
	go test -tags="${GO_BUILD_TAGS}" -covermode=count -coverprofile=${COVER_DIR}/coverage.out ./...
test-cover-html: test-cover
	go tool cover -html=${COVER_DIR}/coverage.out -o ${COVER_DIR}/coverage.html
	@echo "# coverage report generated into ${COVER_DIR}/coverage.html"

perf: ## Run quick performance test
	@echo "# running perf test"
	./tests/perf/perf_test.sh grpc-perf 1000

perf-all: ## Run all performance tests
	@echo "# running all perf tests"
	./tests/perf/run_all.sh

integration-tests: test-tools ## Run integration tests
	@echo "# running integration tests"
	W5GC_IMG=$(W5GC_IMG) ./tests/integration/run_integration.sh

e2e-tests: images test-tools ## Run end-to-end tests
	@echo "# running end-to-end tests"
	W5GC=prod_w5gc ./tests/e2e/run_e2e.sh

# -------------------------------
#  Code generation
# -------------------------------

checknodiffgenerated:  ## Check no diff generated
	bash scripts/checknodiffgenerated.sh $(MAKE) generate

generate: generate-proto generate-binapi generate-desc-adapters ## Generate all

generate-proto: protocgengo ## Generate Protobuf files

get-binapi-generators:
	go install -mod=readonly w5gc.io/wipro5gcore/cmd/binapi-generator

generate-binapi: get-binapi-generators ## Generate Go code for binary API
	@echo "# generating binapi"
	BINAPI=$(BINAPI) ./scripts/genbinapi.sh

verify-binapi: ## Verify generated binary API
	@echo "# verifying generated binapi"
	docker build -f docker/dev/smf/pdusmsp/Dockerfile \
                --target verify-binapi .
	docker build -f docker/dev/amf/csp/Dockerfile \
                --target verify-binapi .

get-desc-adapter-generator:
	go install ./plugins/kvscheduler/descriptor-adapter

generate-desc-adapters: get-desc-adapter-generator ## Generate Go code for descriptors
	@echo "# generating descriptor adapters"
	go generate -x -run=descriptor-adapter ./...

get-bindata:
	go get -v github.com/jteeuwen/go-bindata/...
	go get -v github.com/elazarl/go-bindata-assetfs/...

bindata: get-bindata
	@echo "# generating bindata"
	go generate -x -run=go-bindata-assetfs ./...

proto-schema: ## Generate Protobuf schema image
	@echo "# generating proto schema"
	@$(MAKE) --no-print-directory buf-image

# -------------------------------
#  Dependencies
# -------------------------------

dep-install:
	@echo "# downloading project's dependencies"
	go mod download

dep-update:
	@echo "# updating all dependencies"
	@echo go mod tidy -v

dep-check:
	@echo "# checking dependencies"
	@if ! git --no-pager diff go.mod ; then \
		echo >&2 "go.mod has uncommitted changes!"; \
		exit 1; \
	fi
	go mod verify
	go mod tidy -v
	@if ! git --no-pager diff go.mod ; then \
		echo >&2 "go mod tidy check failed!"; \
		exit 1; \
	fi

# # -------------------------------
# #  Linters
# # -------------------------------

# gotestsumcmd := $(shell command -v gotestsum 2> /dev/null)

# test-tools: ## install test tools
# ifndef gotestsumcmd
# 	go install gotest.tools/gotestsum@v1.8.1

# endif
# 	@env CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/test2json cmd/test2json

# LINTER := $(shell command -v gometalinter 2> /dev/null)

# get-linters:
# ifndef LINTER
# 	@echo "# installing linters"
# 	go install github.com/alecthomas/gometalinter@latest
# 	gometalinter --install
# endif

# lint: get-linters ## Lint Go code
# 	@echo "# running code analysis"
# 	./scripts/static_analysis.sh golint vet

# format: ## Format Go code
# 	@echo "# formatting the code"
# 	./scripts/gofmt.sh

# MDLINKCHECK := $(shell command -v markdown-link-check 2> /dev/null)

# get-linkcheck: ## Check links in Markdown files
# ifndef MDLINKCHECK
# 	sudo apt-get update && sudo apt-get install -y npm
# 	npm install -g markdown-link-check@3.6.2
# endif

# check-links: get-linkcheck
# 	./scripts/check_links.sh

# get-yamllint:
# 	pip install --user yamllint

# yamllint: get-yamllint ## Lint YAML files
# 	@echo "# linting the yaml files"
# 	yamllint -c .yamllint.yml $(shell git ls-files '*.yaml' '*.yml' | grep -v 'vendor/')

# lint-proto: ## Lint Protobuf files
# 	@echo "# linting Protobuf files"
# 	@$(MAKE) --no-print-directory buf-lint

# check-proto: lint-proto ## Check proto files for breaking changes
# 	@echo "# checking proto files"
# 	@$(MAKE) --no-print-directory buf-breaking


.PHONY: help \
        grpcnassmf build clean install purge \
        cmd \
        test test-cover test-cover-html \
        generate checknodiffgenerated generate-binapi generate-proto get-binapi-generators \
        get-dep dep-install dep-update dep-check \
        get-linters lint format lint-proto check-proto \
        get-linkcheck check-links \
        get-yamllint yamllint \
		perf perf-all

