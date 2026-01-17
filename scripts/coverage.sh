#!/bin/bash
# GibRAM Coverage Check
# Displays coverage per package

set -e

echo "Coverage by Package"
echo "==================="

# Run tests with coverage per package
for pkg in $(go list ./pkg/...); do
    COVERAGE=$(go test -cover "$pkg" 2>&1 | grep coverage | awk '{print $2}')
    PKG_SHORT=$(echo "$pkg" | sed 's|.*/||')
    if [ -n "$COVERAGE" ]; then
        printf "%-20s %s\n" "$PKG_SHORT" "$COVERAGE"
    fi
done

echo ""
echo "==================="
go test -cover ./pkg/... 2>&1 | grep total || echo "Run 'scripts/test.sh' for detailed coverage"
