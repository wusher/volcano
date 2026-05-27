.PHONY: help setup build test lint clean coverage all e2e

# Go bin path
GOBIN := $(shell go env GOPATH)/bin

# Default target — show help when `make` is invoked with no args
.DEFAULT_GOAL := help

help:
	@echo "Volcano — make targets"
	@echo ""
	@echo "  make setup     Install Go deps and tools (golangci-lint v2)"
	@echo "  make build     Build the volcano binary"
	@echo "  make test      Run tests and enforce 90% coverage threshold"
	@echo "  make coverage  Run tests and print total coverage"
	@echo "  make lint      Run golangci-lint --fix and prettier on e2e/"
	@echo "  make e2e       Build then run Playwright end-to-end tests"
	@echo "  make all       lint + test + build (full quality check)"
	@echo "  make clean     Remove build artifacts and generated output"
	@echo "  make help      Show this help (default when run with no args)"

# Full quality check
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

# Run linter with auto-fix
lint:
	$(GOBIN)/golangci-lint run --fix
	@cd e2e && npx prettier --write . --log-level error

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total
	@rm -f coverage.out

# Run end-to-end tests (Playwright)
e2e: build
	@echo "Running Playwright e2e tests..."
	cd e2e && npm test

# Clean build artifacts
clean:
	rm -f volcano
	rm -f coverage.out
	rm -rf output/
	rm -rf test-output/
