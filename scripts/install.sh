#!/bin/bash
# Modern installation script for cfctl

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
REPO="siyamsarker/cfctl"
VERSION="${1:-latest}"
CONFIG_DIR="$HOME/.config/cfctl"

# Colors and styling
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# Spinner animation
spinner() {
    local pid=$1
    local delay=0.1
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    while ps -p $pid > /dev/null 2>&1; do
        local temp=${spinstr#?}
        printf " [${CYAN}%c${NC}]  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

# Progress indicator
show_progress() {
    echo -ne "${BLUE}▸${NC} $1"
}

show_success() {
    echo -e "\r${GREEN}✓${NC} $1"
}

show_error() {
    echo -e "\r${RED}✗${NC} $1"
}

show_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

show_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Banner
echo -e "${BOLD}${CYAN}"
cat << "EOF"
   ╔═══════════════════════════════════╗
   ║                                   ║
   ║     CFCTL INSTALLER v1.0.0        ║
   ║     Cloudflare CLI Management     ║
   ║                                   ║
   ╚═══════════════════════════════════╝
EOF
echo -e "${NC}"

# System detection
show_progress "Detecting system configuration..."
sleep 0.5

# System detection
show_progress "Detecting system configuration..."
sleep 0.5
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        show_error "Unsupported architecture: $ARCH"
        echo -e "${DIM}Supported: x86_64 (amd64), arm64${NC}"
        exit 1
        ;;
esac

show_success "System detected: ${BOLD}${OS}-${ARCH}${NC}"

# Check for existing installation
if command -v cfctl &> /dev/null; then
    CURRENT_VERSION=$(cfctl --version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
    show_warning "cfctl is already installed (version: ${CURRENT_VERSION})"
    echo -ne "${YELLOW}?${NC} Do you want to ${BOLD}upgrade/reinstall${NC}? [Y/n]: "
    read -r response
    if [[ "$response" =~ ^[Nn]$ ]]; then
        echo -e "${CYAN}Installation cancelled.${NC}"
        exit 0
    fi
    echo ""
fi

# Check sudo requirements
SUDO=""
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    show_info "Sudo privileges required for installation to ${BOLD}${INSTALL_DIR}${NC}"
    SUDO="sudo"
    # Test sudo access
    if ! sudo -v; then
        show_error "Unable to obtain sudo privileges"
        exit 1
    fi
fi

# Find binary
show_progress "Locating cfctl binary..."
sleep 0.3

BINARY_PATH=""
if [ -f "./bin/cfctl-${OS}-${ARCH}" ]; then
    BINARY_PATH="./bin/cfctl-${OS}-${ARCH}"
elif [ -f "./bin/cfctl" ]; then
    BINARY_PATH="./bin/cfctl"
elif [ -f "./cfctl" ]; then
    BINARY_PATH="./cfctl"
else
    show_error "Binary not found"
    echo -e "${DIM}Expected locations:${NC}"
    echo -e "  ${DIM}• ./bin/cfctl-${OS}-${ARCH}${NC}"
    echo -e "  ${DIM}• ./bin/cfctl${NC}"
    echo -e "  ${DIM}• ./cfctl${NC}"
    echo ""
    echo -e "${YELLOW}Tip:${NC} Build the binary first with: ${BOLD}make build${NC}"
    exit 1
fi

show_success "Binary found: ${DIM}${BINARY_PATH}${NC}"

# Verify binary
show_progress "Verifying binary..."
if file "$BINARY_PATH" | grep -q "executable"; then
    BINARY_SIZE=$(du -h "$BINARY_PATH" | cut -f1)
    show_success "Binary verified (${BINARY_SIZE})"
else
    show_error "Invalid binary file"
    exit 1
fi

# Installation
echo ""
echo -e "${BOLD}Installation Summary:${NC}"
echo -e "  ${DIM}Source:${NC}      ${BINARY_PATH}"
echo -e "  ${DIM}Destination:${NC} ${INSTALL_DIR}/${BINARY_NAME}"
echo -e "  ${DIM}Permissions:${NC} executable"
echo ""

if [ -n "$SUDO" ]; then
    echo -ne "${YELLOW}?${NC} Proceed with installation? [Y/n]: "
    read -r response
    if [[ "$response" =~ ^[Nn]$ ]]; then
        echo -e "${CYAN}Installation cancelled.${NC}"
        exit 0
    fi
    echo ""
fi

show_progress "Installing binary..."
(
    $SUDO cp "$BINARY_PATH" "${INSTALL_DIR}/${BINARY_NAME}" && 
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
) &
spinner $!
wait $!

if [ $? -eq 0 ]; then
    show_success "Binary installed successfully"
else
    show_error "Installation failed"
    exit 1
fi

# Verify installation
show_progress "Verifying installation..."
sleep 0.3
if command -v cfctl &> /dev/null; then
    INSTALLED_VERSION=$(cfctl --version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
    show_success "Installation verified (version: ${INSTALLED_VERSION})"
else
    show_error "Verification failed - cfctl not found in PATH"
    exit 1
fi

# Create config directory
if [ ! -d "$CONFIG_DIR" ]; then
    show_progress "Creating configuration directory..."
    mkdir -p "$CONFIG_DIR"
    show_success "Configuration directory created: ${DIM}${CONFIG_DIR}${NC}"
fi

# Success banner
echo ""
echo -e "${GREEN}${BOLD}╔══════════════════════════════════╗${NC}"
echo -e "${GREEN}${BOLD}║                                  ║${NC}"
echo -e "${GREEN}${BOLD}║    ✓ Installation Successful!    ║${NC}"
echo -e "${GREEN}${BOLD}║                                  ║${NC}"
echo -e "${GREEN}${BOLD}╚══════════════════════════════════╝${NC}"
echo ""
echo -e "${BOLD}Quick Start:${NC}"
echo -e "  ${GREEN}▸${NC} Run the application:    ${BOLD}cfctl${NC}"
echo -e "  ${GREEN}▸${NC} View help:              ${BOLD}cfctl --help${NC}"
echo -e "  ${GREEN}▸${NC} Configure account:      ${BOLD}cfctl${NC} ${DIM}(then select 'Configure Account')${NC}"
echo -e "  ${GREEN}▸${NC} Check version:          ${BOLD}cfctl --version${NC}"
echo ""
echo -e "${DIM}Documentation: ${CYAN}https://github.com/${REPO}${NC}"
echo ""
