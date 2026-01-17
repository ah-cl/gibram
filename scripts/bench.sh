#!/bin/bash
# GibRAM Benchmark Runner
# Runs all benchmarks and generates reports

set -e

echo "============================================="
echo "GibRAM Benchmark Suite"
echo "============================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
BENCHMARK_DIR="benchmarks"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BENCHMARK_FILE="$BENCHMARK_DIR/bench_$TIMESTAMP.txt"
COMPARE_FILE="$BENCHMARK_DIR/bench_latest.txt"

# Create benchmark directory
mkdir -p "$BENCHMARK_DIR"

print_header() {
    echo ""
    echo -e "${YELLOW}>>> $1${NC}"
}

# Store benchmarks
print_header "Running benchmarks (this may take a while...)"

# Run benchmarks and save output
go test -bench=. -benchmem -run=^$ ./pkg/... 2>&1 | tee "$BENCHMARK_FILE"

# Copy to latest for future comparison
cp "$BENCHMARK_FILE" "$COMPARE_FILE"

echo ""
echo -e "${GREEN}Benchmarks completed!${NC}"
echo "Results saved to: $BENCHMARK_FILE"

# Summary of key benchmarks
print_header "Key Benchmark Results"
echo ""
grep -E "Benchmark.*-[0-9]" "$BENCHMARK_FILE" | head -20

echo ""
echo "============================================="
echo "Full results: $BENCHMARK_FILE"
echo "============================================="
