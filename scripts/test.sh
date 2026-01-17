#!/bin/bash
# GibRAM Test Runner
# Runs all tests with coverage reporting

set -e

echo "============================================="
echo "GibRAM Test Suite"
echo "============================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_DIR="coverage"
COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"
MIN_COVERAGE=70 # Minimum coverage percentage

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Function to print header
print_header() {
    echo ""
    echo -e "${YELLOW}>>> $1${NC}"
}

# Check Go installation
print_header "Checking Go installation"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    print_status 0 "Go installed: $GO_VERSION"
else
    print_status 1 "Go is not installed"
fi

# Download dependencies
print_header "Downloading dependencies"
go mod download
print_status $? "Dependencies downloaded"

# Vet code
print_header "Running go vet"
go vet ./...
print_status $? "Go vet passed"

# Run tests with coverage
print_header "Running tests with coverage"
go test -v -race -coverprofile="$COVERAGE_FILE" -covermode=atomic ./pkg/...
TEST_EXIT=$?
print_status $TEST_EXIT "Tests completed"

# Generate HTML coverage report
print_header "Generating coverage report"
go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
print_status $? "Coverage report: $COVERAGE_HTML"

# Check coverage percentage
print_header "Checking coverage threshold"
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')
echo "Total coverage: ${COVERAGE}%"

# Compare with minimum (using bc for float comparison)
if [ $(echo "$COVERAGE >= $MIN_COVERAGE" | bc -l) -eq 1 ]; then
    print_status 0 "Coverage ${COVERAGE}% meets minimum ${MIN_COVERAGE}%"
else
    print_status 1 "Coverage ${COVERAGE}% below minimum ${MIN_COVERAGE}%"
fi

# Summary
echo ""
echo "============================================="
echo -e "${GREEN}All tests passed!${NC}"
echo "Coverage: ${COVERAGE}%"
echo "Report: $COVERAGE_HTML"
echo "============================================="
