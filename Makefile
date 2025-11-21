default: build

ROOT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
OUTPUT_PATH := $(ROOT_DIR)/bin
BIN_NAME := ddns-go
UNAME_S := $(shell uname -s)
GOPATH := $(shell go env GOPATH)
GO_VERSION=$(shell go version | awk '{print $$3}')
GO_MODULE=$(shell go list -m)
BRANCH_NAME=$(shell git rev-parse --abbrev-ref HEAD)
LAST_COMMIT_ID=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
VERSION=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
VERSION_INFO :=-X '$(GO_MODULE)/pkg/version.GoVersion=$(GO_VERSION)' -X '$(GO_MODULE)/pkg/version.BranchName=$(BRANCH_NAME)' \
-X '$(GO_MODULE)/pkg/version.CommitID=$(LAST_COMMIT_ID)' -X '$(GO_MODULE)/pkg/version.BuildTime=$(BUILD_TIME)' \
-X '$(GO_MODULE)/pkg/version.Version=$(VERSION)'

GOLDFLAGS=-ldflags="-s -w -buildid=$(VERSION)-$(LAST_COMMIT_ID) $(VERSION_INFO)"

.PHONY: config
config:
	@cd $(ROOT_DIR) && go mod tidy

.PHONY: build
build: config
	@CGO_ENABLED=0 go build -trimpath $(GOLDFLAGS) -o $(OUTPUT_PATH)/$(BIN_NAME) ./cmd/ddns-go

.PHONY: clean
clean:
	@rm -rf $(OUTPUT_PATH)

.PHONY: image
image:
	@$(ROOT_DIR)/optools/image/build.sh
