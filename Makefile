.PHONY: build clean test install help

# Default target
all: build

# Build the application
build:
	go build -o arbor

# Install the application
install:
	go install

# Run tests
test:
	go test ./...

# Run tests with race detection
test-race:
	go test -race ./...

# Clean build artifacts
clean:
	rm -f arbor
	rm -f *.html

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  install   - Install the application"
	@echo "  test      - Run tests"
	@echo "  test-race - Run tests with race detection"
	@echo "  clean     - Clean build artifacts"
	@echo "  help      - Show this help message"

# Development helpers
dev: build
	./arbor --help

# Cross-compilation examples
build-linux:
	GOOS=linux GOARCH=amd64 go build -o arbor-linux

build-windows:
	GOOS=windows GOARCH=amd64 go build -o arbor.exe

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o arbor-darwin
