.PHONY: build build-all build-darwin build-linux install test lint clean run

# Variables
BINARY_NAME=cfctl
VERSION=1.0.0
BUILD_DIR=bin
MAIN_PATH=./cmd/cfctl

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all: build-darwin build-linux
	@echo "✓ All builds complete"

# Build for macOS (Intel + Apple Silicon)
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "✓ macOS builds complete"

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "✓ Linux builds complete"

# Install locally
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

# Run the application
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✓ Dependencies tidied"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Display help
help:
	@echo "Available targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-darwin   - Build for macOS (Intel + ARM)"
	@echo "  build-linux    - Build for Linux (AMD64 + ARM64)"
	@echo "  install        - Build and install locally"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  clean          - Remove build artifacts"
	@echo "  help           - Display this help message"
