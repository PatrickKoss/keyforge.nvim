.PHONY: build run clean test watch install help lint lint-go lint-lua install-tools

# Default target
all: build

# Build the game binary
build:
	@echo "Building keyforge..."
	@cd game && go build -o bin/keyforge ./cmd/keyforge
	@echo "Built: game/bin/keyforge"

# Run the game
run: build
	@./game/bin/keyforge

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf game/bin
	@echo "Done"

# Run tests
test:
	@echo "Running tests..."
	@cd game && go test -v ./...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@cd game && go test -v ./internal/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@cd game && go test -coverprofile=coverage.out ./...
	@cd game && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: game/coverage.html"

# Run quick tests (unit tests only, skip integration)
test-unit:
	@echo "Running unit tests..."
	@cd game && go test -v -short ./...

# Watch for changes and rebuild (requires entr)
watch:
	@echo "Watching for changes..."
	@find game -name "*.go" | entr -r make run

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	@cd game && go mod tidy

# Install development tools (golangci-lint, luacheck)
install-tools:
	@./scripts/install-tools.sh

# Format code
fmt:
	@echo "Formatting code..."
	@cd game && go fmt ./...

# Lint code (uses local binaries from bin/)
lint: lint-go lint-lua

lint-go:
	@echo "Linting Go code..."
	@cd game && ../bin/golangci-lint run --fix

lint-lua:
	@echo "Linting Lua code..."
	@./bin/luacheck lua/ plugin/ tests/ --config .luacheckrc

# Build for multiple platforms
release:
	@echo "Building releases..."
	@mkdir -p game/bin/release
	@cd game && GOOS=linux GOARCH=amd64 go build -o bin/release/keyforge-linux-amd64 ./cmd/keyforge
	@cd game && GOOS=darwin GOARCH=amd64 go build -o bin/release/keyforge-darwin-amd64 ./cmd/keyforge
	@cd game && GOOS=darwin GOARCH=arm64 go build -o bin/release/keyforge-darwin-arm64 ./cmd/keyforge
	@cd game && GOOS=windows GOARCH=amd64 go build -o bin/release/keyforge-windows-amd64.exe ./cmd/keyforge
	@echo "Releases built in game/bin/release/"

# Install the plugin (symlink for development)
install-dev:
	@echo "Creating development symlink..."
	@mkdir -p ~/.local/share/nvim/lazy
	@ln -sf $(PWD) ~/.local/share/nvim/lazy/keyforge.nvim
	@echo "Symlinked to ~/.local/share/nvim/lazy/keyforge.nvim"

# Help
help:
	@echo "Keyforge Makefile targets:"
	@echo "  build            - Build the game binary"
	@echo "  run              - Build and run the game"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run all tests"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  watch            - Watch for changes and rebuild"
	@echo "  deps             - Install Go dependencies"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run all linters (Go + Lua)"
	@echo "  lint-go          - Run golangci-lint on Go code"
	@echo "  lint-lua         - Run luacheck on Lua code"
	@echo "  install-tools    - Install dev tools (golangci-lint, luacheck)"
	@echo "  release          - Build for all platforms"
	@echo "  install-dev      - Symlink plugin for development"
