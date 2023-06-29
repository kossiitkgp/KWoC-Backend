#!/bin/sh

set -eu -o pipefail

echo "Pre-commit hook running"

echo "Running go mod tidy..."
go mod tidy

if ! command -v golangci-lint &> /dev/null ; then
    echo "golangci-lint not installed or available in the PATH" >&2
    echo "please check https://github.com/golangci/golangci-lint" >&2
    exit 1
fi
echo "Running golangci-lint..."
golangci-lint run

echo "Running go vet..."
go vet ./...

echo "Hook finished running"
