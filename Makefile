.PHONY: build test clean install run dev lint help

# Variables
BINARY_NAME=gh-flow
BUILD_DIR=.
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Default target
.DEFAULT_GOAL := help

## build: Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	go clean
	@echo "✅ Clean complete"

## install: Install binary to GOPATH/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	go install .
	@echo "✅ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

## run: Build and run the binary
run: build
	./$(BINARY_NAME)

## dev: Run with auto-rebuild (requires air)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

## lint: Run linter
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed"; \
		go vet ./...; \
	fi

## fmt: Format Go code
fmt:
	@echo "📝 Formatting code..."
	go fmt ./...

## tidy: Tidy and download dependencies
tidy:
	@echo "📦 Tidying dependencies..."
	go mod tidy
	go mod download

## setup-test: Create test repositories for manual testing
setup-test:
	@echo "🔧 Setting up test repositories..."
	@mkdir -p test-repos
	@cd test-repos && \
		if [ ! -d "frontend" ]; then \
			mkdir frontend && cd frontend && git init && git checkout -b main && \
			echo "# Frontend" > README.md && \
			git add . && git commit -m "Initial commit" && \
			git checkout -b dev && \
			cd ..; \
		fi && \
		if [ ! -d "backend" ]; then \
			mkdir backend && cd backend && git init && git checkout -b main && \
			echo "# Backend" > README.md && \
			git add . && git commit -m "Initial commit" && \
			git checkout -b dev && \
			cd ..; \
		fi
	@echo "✅ Test repositories created in test-repos/"
	@echo ""
	@echo "Directory structure:"
	@ls -la test-repos/

## clean-test: Remove test repositories
clean-test:
	@echo "🧹 Removing test repositories..."
	rm -rf test-repos/
	@echo "✅ Test repositories removed"

## test-start: Test 'start' command in test environment
test-start: build setup-test
	@echo "🚀 Testing 'start' command..."
	cd test-repos && ../$(BINARY_NAME) start

## test-finish: Test 'finish' command in test environment
test-finish: build setup-test
	@echo "🚀 Testing 'finish' command..."
	@echo "Creating test changes in frontend..."
	@echo "// Test change" >> test-repos/frontend/README.md
	@echo "Creating test changes in backend..."
	@echo "// Test change" >> test-repos/backend/README.md
	cd test-repos && ../$(BINARY_NAME) finish

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed 's/## //g' | column -t -s ':'
