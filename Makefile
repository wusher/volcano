.PHONY: setup build test lint clean coverage all e2e

# Go bin path
GOBIN := $(shell go env GOPATH)/bin

# Default target
all: lint test build

# Install dependencies and tools
setup:
	go mod download
	go mod tidy
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "Setup complete. All dependencies installed."

# Build the binary
build:
	go build -o volcano .

# Run tests with coverage threshold
test:
	go test -coverprofile=coverage.out ./...
	@COVERAGE=$$(go tool cover -func=coverage.out | awk '/total/ {gsub("%","",$$3); print $$3}'); \
		echo "Total coverage: $$COVERAGE%"; \
		awk -v c="$$COVERAGE" 'BEGIN { exit (c >= 90) ? 0 : 1 }' || (echo "Coverage $$COVERAGE% is below 90% threshold" && exit 1); \
		rm -f coverage.out

# Run linter
lint:
	$(GOBIN)/golangci-lint run

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total
	@rm -f coverage.out

# Run end-to-end tests
e2e: build
	@echo "Running e2e tests..."
	@rm -rf output/
	./volcano ./example
	@echo "Verifying output structure..."
	@test -f ./output/index.html || (echo "FAIL: index.html not found" && exit 1)
	@test -f ./output/404.html || (echo "FAIL: 404.html not found" && exit 1)
	@echo "Verifying HTML content..."
	@grep -q "<title>" ./output/index.html || (echo "FAIL: title tag not found" && exit 1)
	@grep -q "nav" ./output/index.html || (echo "FAIL: nav element not found" && exit 1)
	@echo "Testing serve mode..."
	@./volcano -s -p 8888 ./output & PID=$$!; \
		sleep 2; \
		curl -sf http://localhost:8888/ > /dev/null || (kill $$PID 2>/dev/null; echo "FAIL: serve mode not responding" && exit 1); \
		kill $$PID 2>/dev/null
	@rm -rf output/
	@echo "All e2e tests passed!"

# Clean build artifacts
clean:
	rm -f volcano
	rm -f coverage.out
	rm -rf output/
	rm -rf test-output/
