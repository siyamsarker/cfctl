#!/bin/bash
# Modern uninstallation script for cfctl

set -e

# ============================================================================
# Configuration
# ============================================================================
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
CONFIG_DIR="$HOME/.config/cfctl"
CACHE_DIR="$HOME/.cache/cfctl"

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

# Status indicators
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

# Get directory size safely
get_dir_size() {
    if [ -d "$1" ]; then
        du -sh "$1" 2>/dev/null | cut -f1 || echo "unknown"
    else
        echo "0B"
    fi
}

# Count files in directory safely
count_files() {
    if [ -d "$1" ]; then
        find "$1" -type f 2>/dev/null | wc -l | tr -d ' ' || echo "0"
    else
        echo "0"
    fi
}

# Cleanup handler
cleanup() {
    if [ $? -ne 0 ]; then
        echo ""
        show_error "Uninstallation encountered errors. Please check the output above."
    fi
}
trap cleanup EXIT

# ============================================================================
# Visual Banner
# ============================================================================
echo ""
echo -e "${BOLD}${YELLOW}"
cat << "EOF"
   ╔════════════════════════════════════════════╗
   ║                                            ║
   ║         CFCTL UNINSTALLER v1.0.0           ║
   ║         Cloudflare CLI Management          ║
   ║                                            ║
   ╚════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# ============================================================================
# Pre-Uninstallation Checks
# ============================================================================

# Check if cfctl is installed
if ! command -v cfctl &> /dev/null && [ ! -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    show_warning "cfctl is not installed"
    echo -e "${DIM}   No binary found in ${INSTALL_DIR}${NC}"
    echo ""
    
    # Check for orphaned config
    if [ -d "$CONFIG_DIR" ] || [ -d "$CACHE_DIR" ]; then
        show_info "Found leftover configuration files"
        echo -ne "${YELLOW}?${NC} Do you want to remove them? [Y/n]: "
        read -r response
        if [[ ! "$response" =~ ^[Nn]$ ]]; then
            echo ""
            if [ -d "$CONFIG_DIR" ]; then
                rm -rf "$CONFIG_DIR" && show_success "Removed configuration directory"
            fi
            if [ -d "$CACHE_DIR" ]; then
                rm -rf "$CACHE_DIR" && show_success "Removed cache directory"
            fi
        fi
    fi
    echo ""
    exit 0
fi

# Get current version
if command -v cfctl &> /dev/null; then
    CURRENT_VERSION=$(cfctl --version 2>/dev/null | grep "Version:" | awk '{print $2}' || echo "unknown")
    show_info "Currently installed version: ${BOLD}${CURRENT_VERSION}${NC}"
fi

# ============================================================================
# Scan Installation
# ============================================================================
echo ""
echo -e "${BOLD}Scanning installation...${NC}"
echo ""

ITEMS_TO_REMOVE=()

# Check binary
if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    BINARY_SIZE=$(du -h "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null | cut -f1 || echo "unknown")
    show_item "Binary: ${DIM}${INSTALL_DIR}/${BINARY_NAME}${NC} ${DIM}(${BINARY_SIZE})${NC}"
    ITEMS_TO_REMOVE+=("binary")
fi

# Check config directory
if [ -d "$CONFIG_DIR" ]; then
    CONFIG_SIZE=$(get_dir_size "$CONFIG_DIR")
    CONFIG_FILES=$(count_files "$CONFIG_DIR")
    show_item "Config: ${DIM}${CONFIG_DIR}${NC} ${DIM}(${CONFIG_FILES} files, ${CONFIG_SIZE})${NC}"
    ITEMS_TO_REMOVE+=("config")
fi

# Check cache directory
if [ -d "$CACHE_DIR" ]; then
    CACHE_SIZE=$(get_dir_size "$CACHE_DIR")
    CACHE_FILES=$(count_files "$CACHE_DIR")
    show_item "Cache:  ${DIM}${CACHE_DIR}${NC} ${DIM}(${CACHE_FILES} files, ${CACHE_SIZE})${NC}"
    ITEMS_TO_REMOVE+=("cache")
fi

if [ ${#ITEMS_TO_REMOVE[@]} -eq 0 ]; then
    echo ""
    show_warning "Nothing to uninstall"
    echo ""
    exit 0
fi

# ============================================================================
# Confirmation Prompt
# ============================================================================
echo ""
echo -e "${BOLD}${YELLOW}╭──────────────────────────────────────╮${NC}"
echo -e "${BOLD}${YELLOW}│${NC}  ${BOLD}${YELLOW}⚠  Warning - Permanent Removal${NC}      ${BOLD}${YELLOW}│${NC}"
echo -e "${BOLD}${YELLOW}╰──────────────────────────────────────╯${NC}"
echo ""
echo -e "${DIM}The following items will be permanently deleted:${NC}"
echo ""
for item in "${ITEMS_TO_REMOVE[@]}"; do
    echo -e "   ${RED}✗${NC}  ${item}"
done

echo ""
echo -ne "${YELLOW}?${NC} ${BOLD}Proceed with uninstallation?${NC} [y/N]: "
read -r response

if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo ""
    echo -e "${CYAN}Uninstallation cancelled by user.${NC}"
    echo ""
    exit 0
fi

# Check sudo requirements
SUDO=""
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    echo ""
    show_info "Sudo privileges required for binary removal"
    SUDO="sudo"
    if ! sudo -v; then
        show_error "Unable to obtain sudo privileges"
        exit 1
    fi
fi

# ============================================================================
# Uninstallation Process
# ============================================================================
echo ""
echo -e "${BOLD}Uninstalling cfctl...${NC}"
echo ""

ERRORS=0

# Remove binary
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " binary " ]]; then
    if $SUDO rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null; then
        show_success "Removed binary from ${INSTALL_DIR}"
    else
        show_error "Failed to remove binary"
        ERRORS=$((ERRORS + 1))
    fi
fi

# Remove config (with separate confirmation for data protection)
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " config " ]]; then
    echo ""
    echo -e "${YELLOW}⚠${NC}  Configuration contains account credentials and settings"
    echo -ne "${YELLOW}?${NC} ${BOLD}Remove configuration and account data?${NC} [y/N]: "
    read -r config_response
    echo ""
    if [[ "$config_response" =~ ^[Yy]$ ]]; then
        if rm -rf "$CONFIG_DIR" 2>/dev/null; then
            show_success "Removed configuration directory"
        else
            show_error "Failed to remove configuration"
            ERRORS=$((ERRORS + 1))
        fi
    else
        show_info "Configuration preserved at ${DIM}${CONFIG_DIR}${NC}"
    fi
fi

# Remove cache
if [[ " ${ITEMS_TO_REMOVE[*]} " =~ " cache " ]]; then
    if rm -rf "$CACHE_DIR" 2>/dev/null; then
        show_success "Removed cache directory"
    else
        show_warning "Could not remove cache (may not exist)"
    fi
fi

# ============================================================================
# Verification & Success Message
# ============================================================================
echo ""

if ! command -v cfctl &> /dev/null && [ ! -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    # Success banner
    echo -e "${GREEN}${BOLD}╔═══════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}${BOLD}║                                               ║${NC}"
    echo -e "${GREEN}${BOLD}║        ✓  Uninstallation Complete!            ║${NC}"
    echo -e "${GREEN}${BOLD}║                                               ║${NC}"
    echo -e "${GREEN}${BOLD}╚═══════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Keyring notice
    echo -e "${BOLD}Additional Cleanup (Optional):${NC}"
    echo ""
    show_info "Credentials may still exist in system keyring"
    echo ""
    echo -e "${DIM}   To remove keyring credentials manually:${NC}"
    echo ""
    
    OS=$(uname -s)
    case $OS in
        Darwin)
            echo -e "   ${DIM}1. Open ${BOLD}Keychain Access${NC}${DIM} application${NC}"
            echo -e "   ${DIM}2. Search for:${NC} ${BOLD}cfctl${NC}"
            echo -e "   ${DIM}3. Delete any matching entries${NC}"
            ;;
        Linux)
            echo -e "   ${DIM}Run:${NC} ${BOLD}secret-tool search service cfctl${NC}"
            echo -e "   ${DIM}Then delete entries as needed${NC}"
            ;;
        *)
            echo -e "   ${DIM}Consult your system's keyring documentation${NC}"
            ;;
    esac
    
    echo ""
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${DIM}Thank you for using cfctl!${NC}"
    echo -e "${DIM}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
else
    show_error "Uninstallation may be incomplete"
    echo ""
    show_info "Please manually check: ${BOLD}${INSTALL_DIR}/${BINARY_NAME}${NC}"
    
    if [ $ERRORS -gt 0 ]; then
        echo ""
        show_warning "Encountered ${ERRORS} error(s) during uninstallation"
    fi
fi

echo ""
