APP_NAME := genlinx
BIN_DIR := bin

MODULE := github.com/Norgate-AV/genlinx

TARGET := $(APP_NAME)
ifeq ($(OS),Windows_NT)
	TARGET := $(APP_NAME).exe
endif

# Get GOBIN, fallback to GOPATH/bin, then HOME/go/bin
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif
ifeq ($(GOBIN),)
GOBIN := $(HOME)/go/bin
endif

# Version info
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

VERSION_PKG := $(MODULE)/internal/version

# Compiler Flags
GOFLAGS := -ldflags="-s -w -extldflags=-static \
	-X '$(VERSION_PKG).Version=$(VERSION)' \
	-X '$(VERSION_PKG).GitCommit=$(GIT_COMMIT)' \
	-X '$(VERSION_PKG).BuildDate=$(BUILD_DATE)'" \
	-trimpath

.PHONY: build
build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(BIN_DIR)/$(TARGET) .

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: run
run: build
	./$(BIN_DIR)/$(TARGET)

install: build
	cp $(BIN_DIR)/$(TARGET) $(GOBIN)/$(APP_NAME)-go.exe

.PHONY: fmt
fmt:
	goimports -local $(MODULE) -w $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

.PHONY: lint
lint:
	golangci-lint run

.PHONY:
man:
	@pandoc docs/genlinx.1.md --to man docs/genlinx.1 > docs/genlinx.1
