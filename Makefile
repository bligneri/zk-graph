# Variables
BINARY_NAME = zk-graph
CMD_DIR = ./cmd/zk-graph
PKG_DIR = ./pkg
OUTPUT_DIR = ./out

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) $(CMD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Clean up generated files and binary
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	find $(OUTPUT_DIR) -type f -name '*.html' -delete

# Build and then run tests
.PHONY: build-and-test
build-and-test: build test

# Install the project binary (optional)
.PHONY: install
install: build
	@echo "Installing binary..."
	mv $(BINARY_NAME) /usr/local/bin/

# Remove installed binary (if installed)
.PHONY: uninstall
uninstall:
	@echo "Uninstalling binary..."
	rm -f /usr/local/bin/$(BINARY_NAME)
