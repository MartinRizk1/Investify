# Makefile for Investify

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
BINARY_NAME=investify
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/main.go

# Default target
all: test build

# Build target
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Run target
run:
	$(GORUN) $(MAIN_PATH)

# Test target
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

# Cross-compile for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Get dependencies
deps:
	$(GOGET) -v ./...

# Install target
install:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)

# Run benchmark tests
benchmark:
	$(GOTEST) -bench=. ./...

# Format go code
fmt:
	$(GOCMD) fmt ./...

# Vet code for potential issues
vet:
	$(GOCMD) vet ./...

.PHONY: all build run test test-coverage clean build-linux deps install benchmark fmt vet
