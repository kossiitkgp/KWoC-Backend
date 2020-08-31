build:
	go build

lint:
	gofmt -s -w .

help:
	@echo "Makefile for automating tasks"
	@echo "build : run make build for building the codebase"
	@echo "lint : run make lint for running lint checks"
