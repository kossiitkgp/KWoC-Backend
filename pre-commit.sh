#!/bin/sh

echo "Pre-commit hook running"

echo "Running go mod tidy..."
go mod tidy

echo "Running golangci-lint..."
golangci-lint run

echo "Running go vet..."
go vet ./...

echo "Hook finished running"
