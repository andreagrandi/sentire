.PHONY: build test clean install help

# Default target
all: build

# Variables for version injection
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X sentire/internal/version.BuildTime=$(BUILD_TIME) -X sentire/internal/version.GitCommit=$(GIT_COMMIT)

# Build the application
build:
	@echo "Building sentire..."
	@go build -ldflags="$(LDFLAGS)" -o sentire ./cmd/sentire

# Build with version information
build-release:
	@echo "Building sentire with version info..."
	@go build -ldflags="$(LDFLAGS) -s -w" -o sentire ./cmd/sentire

# Run tests
test:
	@echo "Running tests..."
	@go test ./tests/ -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./tests/ -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f sentire
	@rm -f coverage.out
	@rm -f coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Install the binary to GOPATH/bin
install: build
	@echo "Installing sentire..."
	@cp sentire $(GOPATH)/bin/

# Cross-compile for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o sentire-linux-amd64 ./cmd/sentire
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o sentire-darwin-amd64 ./cmd/sentire
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o sentire-darwin-arm64 ./cmd/sentire
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o sentire-windows-amd64.exe ./cmd/sentire

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-release - Build with optimizations"
	@echo "  build-all     - Cross-compile for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Remove build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  help          - Show this help message"