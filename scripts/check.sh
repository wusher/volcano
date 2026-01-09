#!/bin/bash
set -e

echo "Running lint..."
golangci-lint run

echo "Running tests with race detector..."
go test -race ./...

echo "Checking coverage..."
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Total coverage: $COVERAGE%"

if (( $(echo "$COVERAGE < 70" | bc -l) )); then
  echo "FAIL: Coverage $COVERAGE% is below 70% threshold"
  exit 1
fi

echo ""
echo "All checks passed! Ready for next story."
