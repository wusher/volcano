.PHONY: setup build test lint clean coverage all

# Default target
all: lint test build

# Install dependencies and tools
setup:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build the binary
build:
	go build -o volcano .

# Run tests with race detector
test:
	go test -race ./...

# Run linter
lint:
	golangci-lint run

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total
	@rm -f coverage.out

# Clean build artifacts
clean:
	rm -f volcano
	rm -f coverage.out
	rm -rf output/
