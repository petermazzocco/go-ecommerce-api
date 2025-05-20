.PHONY: dev build run clean install-air

# Default binary output name
BINARY_NAME=app

# Go and Air settings
GO=go
AIR=air
AIR_CONFIG=.air.toml

# Directory paths
CMD_DIR=cmd/api
DB_DIR=cmd/db
MAIN_FILE=$(CMD_DIR)/main.go
DB_FILE=$(DB_DIR)/main.go

# Run the application with air for hot reloading
dev:
	$(AIR) -c $(AIR_CONFIG) 

# Build the application
build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application without air
run:
	$(GO) run $(MAIN_FILE)

run-db:
	$(GO) run $(DB_FILE)

# Clean build artifacts
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# Install air if not already installed
install-air:
	$(GO) install github.com/cosmtrek/air@latest

# Default command when running make without args
default: dev
