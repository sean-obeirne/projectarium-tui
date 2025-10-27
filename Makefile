.PHONY: help build run test clean fmt vet lint install-hooks install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make check          - Run all checks (fmt, vet, lint, test)"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install-hooks  - Install git pre-commit hooks"
	@echo "  make install-tools  - Install development tools"

# Build the application
build:
	@echo "🏗️  Building..."
	go build -o bin/projectarium-tui

# Run the application
run:
	@echo "🚀 Running..."
	go run .

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Format code
fmt:
	@echo "📝 Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	go vet ./...

# Run golangci-lint
lint:
	@echo "🔒 Running golangci-lint..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run: make install-tools"; \
	fi

# Run all checks
check:
	@echo "🔧 Running go mod tidy..."
	@go mod tidy
	@echo "📝 Formatting code..."
	@go fmt ./...
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "🧪 Running tests..."
	@go test ./... -short
	@if command -v golangci-lint > /dev/null; then \
		echo "🔒 Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, skipping..."; \
	fi
	@echo "✅ All checks passed!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	rm -f build-errors.log

# Install git hooks
install-hooks:
	@echo "🪝 Installing git hooks..."
	@mkdir -p scripts/hooks
	@echo '#!/bin/sh' > scripts/hooks/pre-commit
	@echo 'make check' >> scripts/hooks/pre-commit
	@chmod +x scripts/hooks/pre-commit
	@ln -sf ../../scripts/hooks/pre-commit .git/hooks/pre-commit
	@echo "✅ Git hooks installed!"

# Install development tools
install-tools:
	@echo "📦 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed!"
