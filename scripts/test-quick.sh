#!/bin/bash
# GibRAM Quick Test
# Runs fast tests without race detection for quick feedback

set -e

echo "Quick Test (no race detection)"
echo "================================"

go test -short ./pkg/...

echo ""
echo "Quick tests passed!"
