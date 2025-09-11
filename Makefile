# Makefile for the lazynet project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTIDY=$(GOCMD) mod tidy
GORUN=$(GOCMD) run

# Binary name
BINARY_NAME=lazynet
BINARY_UNIX=$(BINARY_NAME)

# Main package location
MAIN_PACKAGE=./cmd/lazynet/main.go

all: help

# Build the binary for the current operating system
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(BINARY_NAME) built successfully."

# Run the application
run:
	@echo "Running lazynet..."
	@$(GORUN) $(MAIN_PACKAGE)

# Clean up the binary
clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@echo "Cleanup complete."

# Tidy the go modules
tidy:
	@echo "Tidying go modules..."
	@$(GOTIDY)
	@echo "Modules tidied."

# Display help information
help:
	@echo "-----------------------------------------"
	@echo " lazynet Makefile Help                 "
	@echo "-----------------------------------------"
	@echo " make build   - Build the application binary"
	@echo " make run     - Run the application"
	@echo " make clean   - Remove the built binary"
	@echo " make tidy    - Tidy go.mod dependencies"
	@echo "-----------------------------------------"

.PHONY: all build run clean tidy help
