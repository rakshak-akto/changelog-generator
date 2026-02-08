.PHONY: build test run clean install help

# Binary name
BINARY_NAME=changelog-generator

# Build variables
BUILD_DIR=./bin
GO=go
GOFLAGS=-v

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cli

install: ## Install the binary to GOPATH/bin
	@echo "Installing..."
	$(GO) install ./cmd/cli

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

run: build ## Build and run the binary
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@$(GO) clean

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

lint: fmt vet ## Run formatters and linters
	@echo "Linting complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

.DEFAULT_GOAL := help
