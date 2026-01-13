#!/bin/bash
# Uninstallation script for cfctl

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
CONFIG_DIR="$HOME/.config/cfctl"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Uninstalling cfctl...${NC}"

# Check if running as root for uninstallation
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    SUDO="sudo"
else
    SUDO=""
fi

# Remove binary
if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    $SUDO rm "${INSTALL_DIR}/${BINARY_NAME}"
    echo -e "${GREEN}✓ Removed binary from ${INSTALL_DIR}${NC}"
else
    echo -e "${YELLOW}Binary not found in ${INSTALL_DIR}${NC}"
fi

# Ask about removing config
echo ""
read -p "Do you want to remove configuration and stored credentials? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "$CONFIG_DIR" ]; then
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}✓ Removed configuration directory${NC}"
    fi
    
    # Note about keyring
    echo -e "${YELLOW}Note: Credentials stored in system keyring were NOT automatically removed.${NC}"
    echo -e "${YELLOW}You may want to remove them manually using your system's keychain management tool.${NC}"
fi

echo ""
echo -e "${GREEN}✓ cfctl uninstalled successfully${NC}"
