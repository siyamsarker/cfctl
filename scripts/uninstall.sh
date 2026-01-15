#!/bin/bash
# Modern uninstallation script for cfctl

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
CONFIG_DIR="$HOME/.config/cfctl"
CACHE_DIR="$HOME/.cache/cfctl"

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

# Helper functions
show_success() {
    echo -e "${GREEN}✓${NC} $1"
}

show_error() {
    echo -e "${RED}✗${NC} $1"
}

show_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

show_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

show_item() {
    echo -e "  ${DIM}•${NC} $1"
}

# Get directory size
get_dir_size() {
    if [ -d "$1" ]; then
        du -sh "$1" 2>/dev/null | cut -f1
    else
        echo "0B"
    fi
}

# Count files in directory
count_files() {
    if [ -d "$1" ]; then
        find "$1" -type f 2>/dev/null | wc -l | tr -d ' '
    else
        echo "0"
    fi
}

# Banner
echo -e "${BOLD}${YELLOW}"
cat << "EOF"
   ╔═══════════════════════════════════╗
   ║                                   ║
   ║     CFCTL UNINSTALLER v1.0.0      ║
   ║     Cloudflare CLI Management     ║
   ║                                   ║
   ╚═══════════════════════════════════╝
EOF
echo -e "${NC}"

# Check if cfctl is installed
if ! command -v cfctl &> /dev/null && [ ! -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    show_warning "cfctl is not installed"
    echo -e "${DIM}No binary found in ${INSTALL_DIR}${NC}"
    echo ""
    
    # Check for orphaned config
    if [ -d "$CONFIG_DIR" ] || [ -d "$CACHE_DIR" ]; then
        show_info "Found leftover configuration files"
        echo -ne "${YELLOW}?${NC} Do you want to remove them? [Y/n]: "
        read -r response
        if [[ ! "$response" =~ ^[Nn]$ ]]; then
            [ -d "$CONFIG_DIR" ] && rm -rf "$CONFIG_DIR" && show_success "Removed ${CONFIG_DIR}"
            [ -d "$CACHE_DIR" ] && rm -rf "$CACHE_DIR" && show_success "Removed ${CACHE_DIR}"
        fi
    fi
    exit 0
fi

# Get current version
if command -v cfctl &> /dev/null; then
    CURRENT_VERSION=$(cfctl --version 2>/dev/null | sed -n 's/.*v\([0-9.]\+\).*/\1/p' || echo "unknown")
    show_info "Current version: ${BOLD}${CURRENT_VERSION}${NC}"
fi

# Scan installation
echo ""
echo -e "${BOLD}Scanning installation...${NC}"
echo ""

ITEMS_TO_REMOVE=()
TOTAL_SIZE=0

# Check binary
if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    BINARY_SIZE=$(du -h "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null | cut -f1 || echo "unknown")
    show_item "Binary: ${DIM}${INSTALL_DIR}/${BINARY_NAME}${NC} (${BINARY_SIZE})"
    ITEMS_TO_REMOVE+=("binary")
fi

# Check config directory
if [ -d "$CONFIG_DIR" ]; then
    CONFIG_SIZE=$(get_dir_size "$CONFIG_DIR")
    CONFIG_FILES=$(count_files "$CONFIG_DIR")
    show_item "Config: ${DIM}${CONFIG_DIR}${NC} (${CONFIG_FILES} files, ${CONFIG_SIZE})"
    ITEMS_TO_REMOVE+=("config")
fi

# Check cache directory
if [ -d "$CACHE_DIR" ]; then
    CACHE_SIZE=$(get_dir_size "$CACHE_DIR")
    CACHE_FILES=$(count_files "$CACHE_DIR")
    show_item "Cache:  ${DIM}${CACHE_DIR}${NC} (${CACHE_FILES} files, ${CACHE_SIZE})"
    ITEMS_TO_REMOVE+=("cache")
fi

if [ ${#ITEMS_TO_REMOVE[@]} -eq 0 ]; then
    show_warning "Nothing to uninstall"
    exit 0
fi

# Confirmation prompt
echo ""
echo -e "${BOLD}${YELLOW}⚠  Warning${NC}"
echo -e "${DIM}The following will be permanently removed:${NC}"
for item in "${ITEMS_TO_REMOVE[@]}"; do
    echo -e "  ${RED}✗${NC} ${item}"
done

echo ""
echo -ne "${YELLOW}?${NC} ${BOLD}Proceed with uninstallation?${NC} [y/N]: "
read -r response

if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo -e "${CYAN}Uninstallation cancelled.${NC}"
    exit 0
fi

# Check sudo requirements
SUDO=""
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    show_info "Sudo privileges required"
    SUDO="sudo"
    if ! sudo -v; then
        show_error "Unable to obtain sudo privileges"
        exit 1
    fi
fi

echo ""
echo -e "${BOLD}Uninstalling...${NC}"
echo ""

# Remove binary
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " binary " ]]; then
    if $SUDO rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null; then
        show_success "Removed binary"
    else
        show_error "Failed to remove binary"
    fi
fi

# Remove config (with separate confirmation for data)
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " config " ]]; then
    echo -ne "${YELLOW}?${NC} Remove configuration and account data? [y/N]: "
    read -r config_response
    if [[ "$config_response" =~ ^[Yy]$ ]]; then
        if rm -rf "$CONFIG_DIR" 2>/dev/null; then
            show_success "Removed configuration"
        else
            show_error "Failed to remove configuration"
        fi
    else
        show_info "Configuration preserved"
    fi
fi

# Remove cache
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " cache " ]]; then
    if rm -rf "$CACHE_DIR" 2>/dev/null; then
        show_success "Removed cache"
    else
        show_warning "No cache to remove"
    fi
fi

# Verify removal
echo ""
if ! command -v cfctl &> /dev/null && [ ! -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    # Success banner
    echo -e "${GREEN}${BOLD}╔════════════════════════════════════╗${NC}"
    echo -e "${GREEN}${BOLD}║                                    ║${NC}"
    echo -e "${GREEN}${BOLD}║    ✓ Uninstallation Complete!      ║${NC}"
    echo -e "${GREEN}${BOLD}║                                    ║${NC}"
    echo -e "${GREEN}${BOLD}╚════════════════════════════════════╝${NC}"
    echo ""
    
    # Keyring notice
    show_info "${BOLD}Important:${NC} Credentials in system keyring were preserved"
    echo -e "${DIM}  To remove them manually:${NC}"
    
    OS=$(uname -s)
    case $OS in
        Darwin)
            echo -e "${DIM}  • Open Keychain Access app${NC}"
            echo -e "${DIM}  • Search for 'cfctl'${NC}"
            echo -e "${DIM}  • Delete any matching entries${NC}"
            ;;
        Linux)
            echo -e "${DIM}  • Run: secret-tool search service cfctl${NC}"
            echo -e "${DIM}  • Delete entries as needed${NC}"
            ;;
    esac
    
    echo ""
    echo -e "${DIM}Thank you for using cfctl!${NC}"
else
    show_error "Uninstallation may be incomplete"
    show_info "Please check ${INSTALL_DIR}/${BINARY_NAME} manually"
fi

echo ""
