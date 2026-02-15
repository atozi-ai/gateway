.PHONY: build run clean test

# Binary name
BINARY_NAME=gateway

# Get version from git tag or use dev
VERSION := $(shell git describe --tags --exact-match HEAD 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "dev")
VERSION := $(shell echo $(VERSION) | sed 's/^v//')
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# Build the application
build:
	@echo "Building $(BINARY_NAME) (version: $(VERSION))..."
	@go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Build for Linux
build-linux:
	@echo "Building for Linux (version: $(VERSION))..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-$(VERSION)-linux-amd64 ./cmd/server

# Build for Windows
build-windows:
	@echo "Building for Windows (version: $(VERSION))..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-$(VERSION)-windows-amd64.exe ./cmd/server

# Build for macOS (Intel)
build-darwin:
	@echo "Building for macOS Intel (version: $(VERSION))..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-$(VERSION)-darwin-amd64 ./cmd/server

# Build for macOS (Apple Silicon)
build-darwin-arm64:
	@echo "Building for macOS Apple Silicon (version: $(VERSION))..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-$(VERSION)-darwin-arm64 ./cmd/server

# Build for all platforms
build-all: build-linux build-windows build-darwin build-darwin-arm64
	@echo "Build complete for all platforms (version: $(VERSION))"
