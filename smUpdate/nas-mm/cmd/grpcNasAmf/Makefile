GRPCNASAMF_BINARY_NAME=grpcnasamf

# Go related variables.
GOBASE=$(shell pwd)/../../../
GOPATH="$(GOBASE)/vendor:$(GOBASE)/../../../"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

go-build:
        @echo "  >  Building grpc nas AMF binary..."
        @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(GRPCNASAMF_BINARY_NAME) $(GOFILES)


go-get:
        @echo "  >  Checking if there is any missing grpc dependencies..."
        @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-install:
        @echo "  >  Installing grpc nas AMF binary..."
        @go install $(GOFILES)

go-test:
        @echo "  >  Executing grpc nas AMF test cases ..."
        @GOPATH=$(GOPATH) go test ./...

go-clean:
        @echo "  >  Cleaning build cache"
        @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean


.PHONY: help build test
all: help
help: Makefile
        @echo
        @echo " Choose a command run in "$(PROJECTNAME)":"
        @echo
        @sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
        @echo

