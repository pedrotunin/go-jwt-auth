# Define variables for your project
GO=go
AIR=air
MAIN_FILE=cmd/api/api.go
BUILD_DIR=bin
APP_NAME=api

# Default target, which builds and runs the app
.PHONY: all
all: run

# Install dependencies if needed
.PHONY: deps
deps:
	$(GO) mod tidy

# Build the project
.PHONY: build
build:
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run the project using air (for live-reloading during development)
.PHONY: run
run: build
	$(AIR) -c .air.toml

# Clean the build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)/*

# Run the project directly without air (useful for production)
.PHONY: run-prod
run-prod: build
	./$(BUILD_DIR)/$(APP_NAME)

# Watch for changes and restart the app automatically (air with live reload)
.PHONY: watch
watch:
	$(AIR) -c .air.toml

