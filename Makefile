# Variables
TEST_DIR ?= ./...

# Default target
all: help

# Help target to list available options
help:
	@echo "Usage: make [target] [TEST_DIR=./<directory>]"
	@echo ""
	@echo "Targets:"
	@echo "  test              Run tests in a specific directory"
	@echo "  cover             Shows code coverage"
	@echo "  clean             Clean test cache"

test:
	@echo "Running tests in directory '$(TEST_DIR)'"
	go test $(TEST_DIR)

cover:
	@echo "Running tests in directory '$(TEST_DIR)'"
	go test $(TEST_DIR) -cover

# Clean test cache
clean:
	@echo "Cleaning test cache..."
	go clean -testcache
