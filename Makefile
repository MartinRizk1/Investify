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

# Setup Python environment for TensorFlow models (with graceful failure)
python-setup:
	@echo "Setting up Python virtual environment for TensorFlow models..."
	python3 -m venv models/venv || true
	@echo "Installing required Python packages (TensorFlow is optional)..."
	. models/venv/bin/activate && pip3 install numpy pandas matplotlib scikit-learn joblib || true
	@echo "Trying to install TensorFlow (may fail on some systems)..."
	. models/venv/bin/activate && pip3 install tensorflow || true
	@echo "Installing yfinance for stock data..."
	. models/venv/bin/activate && pip3 install yfinance || true
	@echo "Python setup completed (with possibly limited functionality)"

# Test Python environment and TensorFlow
test-python:
	@echo "Testing Python environment and TensorFlow setup..."
	. models/venv/bin/activate && cd models && python3 test_environment.py

# Train a basic model for testing
train-test-model:
	@echo "Training a basic model for testing..."
	. models/venv/bin/activate && cd models && python3 train_stock_model.py AAPL --quick-test

# Train models for popular stocks
train-models:
	@echo "Training TensorFlow models for popular stocks..."
	. models/venv/bin/activate && cd models && python3 train_stock_model.py AAPL
	. models/venv/bin/activate && cd models && python3 train_stock_model.py MSFT
	. models/venv/bin/activate && cd models && python3 train_stock_model.py GOOGL
	. models/venv/bin/activate && cd models && python3 train_stock_model.py AMZN
	. models/venv/bin/activate && cd models && python3 train_stock_model.py TSLA

# Full setup - Go and Python
setup: deps python-setup test-python

.PHONY: all build run test test-coverage clean build-linux deps install benchmark fmt vet python-setup test-python train-test-model train-models setup
