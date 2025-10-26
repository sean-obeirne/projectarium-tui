.PHONY: build run clean install help

# Build the application
build:
	go build -v -o projectarium-tui

# Run the application
run: build
	./projectarium-tui

# Clean build artifacts
clean:
	rm -f projectarium-tui
	go clean

# Install dependencies
install:
	go mod download
	go mod tidy

# Display help
help:
	@echo "Available targets:"
	@echo "  build   - Build the application"
	@echo "  run     - Build and run the application"
	@echo "  clean   - Clean build artifacts"
	@echo "  install - Install/update dependencies"
	@echo "  help    - Display this help message"
