#!/usr/bin/env bash
set -euo pipefail

# Tool versions (pinned)
GOLANGCI_LINT_VERSION="v2.7.2"
LUACHECK_VERSION="1.2.0"

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_ROOT/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Create bin directory
mkdir -p "$BIN_DIR"

install_golangci_lint() {
    local target="$BIN_DIR/golangci-lint"

    if [[ -x "$target" ]]; then
        local current_version
        current_version=$("$target" --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
        if [[ "$current_version" == "$GOLANGCI_LINT_VERSION" ]]; then
            info "golangci-lint $GOLANGCI_LINT_VERSION already installed"
            return 0
        fi
        warn "golangci-lint version mismatch (have: $current_version, want: $GOLANGCI_LINT_VERSION), reinstalling..."
    fi

    info "Installing golangci-lint $GOLANGCI_LINT_VERSION..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$BIN_DIR" "$GOLANGCI_LINT_VERSION"

    if [[ -x "$target" ]]; then
        info "golangci-lint installed successfully"
    else
        error "Failed to install golangci-lint"
    fi
}

install_luacheck() {
    local target="$BIN_DIR/luacheck"

    if [[ -x "$target" ]]; then
        local current_version
        current_version=$("$target" --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown")
        if [[ "$current_version" == "$LUACHECK_VERSION" ]]; then
            info "luacheck $LUACHECK_VERSION already installed"
            return 0
        fi
        warn "luacheck version mismatch (have: $current_version, want: $LUACHECK_VERSION), reinstalling..."
    fi

    # Try different installation methods in order of preference
    if command -v luarocks &> /dev/null; then
        install_luacheck_luarocks "$target"
    elif command -v docker &> /dev/null; then
        install_luacheck_docker "$target"
    else
        error "Neither luarocks nor docker is available. Please install one of:
  - luarocks: sudo apt install luarocks (Ubuntu/Debian)
  - docker: https://docs.docker.com/get-docker/"
    fi
}

install_luacheck_luarocks() {
    local target="$1"
    local lua_modules="$PROJECT_ROOT/lua_modules"

    info "Installing luacheck $LUACHECK_VERSION via luarocks..."

    # Install luacheck locally using luarocks
    luarocks install --tree="$lua_modules" luacheck "$LUACHECK_VERSION"

    # Create wrapper script in bin directory
    cat > "$target" << 'WRAPPER'
#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
export LUA_PATH="$PROJECT_ROOT/lua_modules/share/lua/5.1/?.lua;$PROJECT_ROOT/lua_modules/share/lua/5.1/?/init.lua;;"
export LUA_CPATH="$PROJECT_ROOT/lua_modules/lib/lua/5.1/?.so;;"
exec "$PROJECT_ROOT/lua_modules/bin/luacheck" "$@"
WRAPPER
    chmod +x "$target"

    if [[ -x "$target" ]]; then
        info "luacheck installed successfully via luarocks"
    else
        error "Failed to install luacheck"
    fi
}

install_luacheck_docker() {
    local target="$1"
    local image="ghcr.io/lunarmodules/luacheck:v${LUACHECK_VERSION}"

    info "Installing luacheck $LUACHECK_VERSION via docker..."

    # Pull the docker image
    if ! docker pull "$image" > /dev/null 2>&1; then
        error "Failed to pull docker image: $image"
    fi

    # Create wrapper script that runs luacheck via docker
    cat > "$target" << WRAPPER
#!/usr/bin/env bash
# Luacheck wrapper using Docker
# Image: $image

# Get the project root (parent of bin/)
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="\$(dirname "\$SCRIPT_DIR")"

# Run luacheck in docker, mounting project root
exec docker run --rm -v "\$PROJECT_ROOT:/data" -w /data "$image" "\$@"
WRAPPER
    chmod +x "$target"

    # Verify wrapper works
    if [[ -x "$target" ]] && "$target" --version > /dev/null 2>&1; then
        info "luacheck installed successfully via docker"
    else
        error "Failed to install luacheck docker wrapper"
    fi
}

main() {
    info "Installing development tools to $BIN_DIR..."
    echo ""

    install_golangci_lint
    echo ""
    install_luacheck
    echo ""

    info "All tools installed successfully!"
    info "Add $BIN_DIR to your PATH or use 'make lint' to run linters"
}

# Allow running individual installers
case "${1:-all}" in
    golangci-lint) install_golangci_lint ;;
    luacheck) install_luacheck ;;
    all) main ;;
    *) error "Unknown tool: $1. Use 'golangci-lint', 'luacheck', or 'all'" ;;
esac
