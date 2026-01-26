#!/bin/bash
# Modern installation script for cfctl

set -e

# ============================================================================
# Configuration
# ============================================================================
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
REPO="siyamsarker/cfctl"
VERSION="${1:-latest}"
CONFIG_DIR="$HOME/.config/cfctl"

# ============================================================================
# Visual Styling
# ============================================================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'
NC='\033[0m'

# ============================================================================
# Helper Functions
# ============================================================================

# Animated spinner for background processes
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

# Status indicators
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

# Cleanup handler
cleanup() {
    if [ $? -ne 0 ]; then
        echo ""
        show_error "Installation failed. Please check the errors above."
    fi
}
trap cleanup EXIT

# ============================================================================
# Visual Banner
# ============================================================================
echo ""
echo -e "${BOLD}${CYAN}"
cat << "EOF"
   ╔═══════════════════════════════════════════════╗
   ║                                               ║
   ║          CFCTL INSTALLER v1.0.0               ║
   ║          Cloudflare CLI Management            ║
   ║                                               ║
   ╚═══════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# ============================================================================
# Pre-Installation Checks
# ============================================================================

# Check required commands
show_progress "Checking dependencies..."
MISSING_DEPS=()
for cmd in uname file du; do
    if ! command -v $cmd &> /dev/null; then
        MISSING_DEPS+=("$cmd")
    fi
done

if [ ${#MISSING_DEPS[@]} -gt 0 ]; then
    show_error "Missing required commands: ${MISSING_DEPS[*]}"
    exit 1
fi
show_success "Dependencies satisfied"

# System detection
show_progress "Detecting system configuration..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        show_error "Unsupported architecture: $ARCH"
        echo -e "${DIM}   Supported architectures: x86_64 (amd64), arm64${NC}"
        exit 1
        ;;
esac

show_success "System detected: ${BOLD}${OS}-${ARCH}${NC}"

# Check for existing installation
if command -v cfctl &> /dev/null; then
    CURRENT_VERSION=$(cfctl --version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
    echo ""
    show_warning "cfctl is already installed ${DIM}(version: ${CURRENT_VERSION})${NC}"
    echo -ne "${YELLOW}?${NC} Do you want to ${BOLD}upgrade/reinstall${NC}? [Y/n]: "
    read -r response
    if [[ "$response" =~ ^[Nn]$ ]]; then
        echo ""
        echo -e "${CYAN}Installation cancelled by user.${NC}"
        echo ""
        exit 0
    fi
fi

# Check sudo requirements
SUDO=""
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    echo ""
    show_info "Sudo privileges required for installation to ${BOLD}${INSTALL_DIR}${NC}"
    SUDO="sudo"
    # Test sudo access
    if ! sudo -v; then
        show_error "Unable to obtain sudo privileges"
        exit 1
    fi
fi

# ============================================================================
# Binary Location & Verification
# ============================================================================
echo ""
show_progress "Locating cfctl binary..."

BINARY_PATH=""
if [ -f "./bin/cfctl-${OS}-${ARCH}" ]; then
    BINARY_PATH="./bin/cfctl-${OS}-${ARCH}"
elif [ -f "./bin/cfctl" ]; then
    BINARY_PATH="./bin/cfctl"
elif [ -f "./cfctl" ]; then
    BINARY_PATH="./cfctl"
else
    show_error "Binary not found"
    echo ""
    echo -e "${DIM}   Expected locations:${NC}"
    echo -e "   ${DIM}• ./bin/cfctl-${OS}-${ARCH}${NC}"
    echo -e "   ${DIM}• ./bin/cfctl${NC}"
    echo -e "   ${DIM}• ./cfctl${NC}"
    echo ""
    echo -e "${YELLOW}   Tip:${NC} Build the binary first with: ${BOLD}make build${NC}"
    echo ""
    exit 1
fi

show_success "Binary located: ${DIM}${BINARY_PATH}${NC}"

# Verify binary integrity
show_progress "Verifying binary integrity..."
if ! file "$BINARY_PATH" | grep -q "executable"; then
    show_error "Invalid binary file (not an executable)"
    exit 1
fi

BINARY_SIZE=$(du -h "$BINARY_PATH" | cut -f1)
show_success "Binary verified ${DIM}(${BINARY_SIZE})${NC}"

# ============================================================================
# Installation Summary & Confirmation
# ============================================================================
echo ""
echo -e "${BOLD}╭─────────────────────────────────────────────╮${NC}"
echo -e "${BOLD}│${NC}  ${BOLD}Installation Summary${NC}                       ${BOLD}│${NC}"
echo -e "${BOLD}├─────────────────────────────────────────────┤${NC}"
echo -e "${BOLD}│${NC}  ${DIM}Source:${NC}       ${BINARY_PATH}${BOLD}                  │${NC}"
echo -e "${BOLD}│${NC}  ${DIM}Destination:${NC}  ${INSTALL_DIR}/${BINARY_NAME}${BOLD}         │${NC}"
echo -e "${BOLD}│${NC}  ${DIM}Permissions:${NC}  executable                   ${BOLD}│${NC}"
echo -e "${BOLD}│${NC}  ${DIM}Size:${NC}         ${BINARY_SIZE}                         ${BOLD}│${NC}"
echo -e "${BOLD}╰─────────────────────────────────────────────╯${NC}"

if [ -n "$SUDO" ]; then
    echo ""
    echo -ne "${YELLOW}?${NC} ${BOLD}Proceed with installation?${NC} [Y/n]: "
    read -r response
    if [[ "$response" =~ ^[Nn]$ ]]; then
        echo ""
        echo -e "${CYAN}Installation cancelled by user.${NC}"
        echo ""
        exit 0
    fi
fi

# ============================================================================
# Installation Process
# ============================================================================
echo ""
echo -e "${BOLD}Installing cfctl...${NC}"
echo ""

show_progress "Copying binary to ${INSTALL_DIR}..."
(
    $SUDO cp "$BINARY_PATH" "${INSTALL_DIR}/${BINARY_NAME}" && 
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
) &
spinner $!
wait $!

if [ $? -eq 0 ]; then
    show_success "Binary installed successfully"
else
    show_error "Failed to install binary"
    exit 1
fi

# Verify installation
show_progress "Verifying installation..."
sleep 0.2
if command -v cfctl &> /dev/null; then
    INSTALLED_VERSION=$(cfctl --version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
    show_success "Installation verified ${DIM}(version: ${INSTALLED_VERSION})${NC}"
else
    show_error "Verification failed - cfctl not found in PATH"
    echo ""
    echo -e "${DIM}   Please ensure ${INSTALL_DIR} is in your PATH${NC}"
    exit 1
fi

# Create config directory
if [ ! -d "$CONFIG_DIR" ]; then
    show_progress "Creating configuration directory..."
    if mkdir -p "$CONFIG_DIR" 2>/dev/null; then
        show_success "Configuration directory created"
    else
        show_warning "Failed to create config directory (non-fatal)"
    fi
fi

# ============================================================================
# Success Message
# ============================================================================
echo ""
echo -e "${GREEN}${BOLD}╔═══════════════════════════════════════════╗${NC}"
echo -e "${GREEN}${BOLD}║                                           ║${NC}"
echo -e "${GREEN}${BOLD}║       ✓  Installation Successful!         ║${NC}"
echo -e "${GREEN}${BOLD}║                                           ║${NC}"
echo -e "${GREEN}${BOLD}╚═══════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BOLD}Quick Start Guide:${NC}"
echo ""
echo -e "  ${GREEN}▸${NC} Launch application     ${BOLD}cfctl${NC}"
echo -e "  ${GREEN}▸${NC} View help              ${BOLD}cfctl --help${NC}"
echo -e "  ${GREEN}▸${NC} Check version          ${BOLD}cfctl --version${NC}"
echo -e "  ${GREEN}▸${NC} Configure account      ${BOLD}cfctl${NC} ${DIM}→ select 'Configure Account'${NC}"
echo ""
echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${DIM}Documentation:${NC} ${CYAN}https://github.com/${REPO}${NC}"
echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
