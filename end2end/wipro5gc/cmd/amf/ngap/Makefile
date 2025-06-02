NGAP_BINARY_NAME=ngap

# Go related variables.
GOBASE=$(shell pwd)/../../../
GOPATH="$(GOBASE)/vendor:$(GOBASE)/../../../"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

go-build:
        @echo "  >  Building sms ngap binary..."
        @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(NGAP_BINARY_NAME) $(GOFILES)


go-get:
        @echo "  >  Checking if there is any missing ngap dependencies..."
        @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-install:
        @echo "  >  Installing amf ngap binary..."
        @go install $(GOFILES)

go-test:
        @echo "  >  Executing amf ngap test cases ..."
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

