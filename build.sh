#!/bin/bash

set -e

echo "Running go mod tidy"
go mod tidy

echo "Running lint"
gofmt -w -s -l .

echo "Building"
go build ./cmd/backend.go

echo "Build complete"