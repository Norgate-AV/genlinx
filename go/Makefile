.PHONY: all build clean test fmt vet mod-tidy

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Main package
MAIN_PACKAGE=.
BINARY_NAME=genlinx
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the project
all: clean mod-tidy fmt vet test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PACKAGE)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v $(MAIN_PACKAGE)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

mod-tidy:
	$(GOMOD) tidy

mod-download:
	$(GOMOD) download

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BINARY_NAME)

# Development helpers
dev: fmt vet test build

# Install dependencies
deps:
	$(GOGET) -v ./...

# Cross-platform build
build-all: build-linux build-windows

# Docker build (if needed)
docker-build:
	docker build -t genlinx .

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run clean, mod-tidy, fmt, vet, test, and build"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  mod-tidy     - Tidy go modules"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Development build (fmt, vet, test, build)"
	@echo "  deps         - Install dependencies"
	@echo "  build-all    - Build for Linux and Windows"
	@echo "  help         - Show this help"
