.PHONY: build test install clean run help lint coverage

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=tatsu

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOVET=$(GOCMD) vet
GOFMT=gofmt

# Build the project
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v

# Install the binary to $GOPATH/bin
install: ## Install the binary globally
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL)

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -cover ./...

# Run tests with detailed coverage
coverage-html: ## Generate HTML coverage report
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linters
lint: ## Run linters and formatters
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Checking formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Format code
fmt: ## Format all Go code
	@echo "Formatting code..."
	$(GOFMT) -w .

# Clean build artifacts
clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Run tatsu (requires tatsu.yaml in current directory)
run: build ## Build and run tatsu with example task
	@echo "Running tatsu..."
	./$(BINARY_NAME) version

# Development workflow: format, lint, test, build
all: fmt lint test build ## Run fmt, lint, test, and build

# Help target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
