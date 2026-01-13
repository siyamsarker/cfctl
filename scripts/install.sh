#!/bin/bash
# Installation script for cfctl

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="cfctl"
REPO="siyamsarker/cfctl"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing cfctl...${NC}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Detected: ${OS}-${ARCH}"

# Check if running as root for installation
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Note: You may need sudo privileges to install to ${INSTALL_DIR}${NC}"
    SUDO="sudo"
else
    SUDO=""
fi

# Download from GitHub releases (when available)
# For now, check if binary exists locally
if [ -f "./bin/cfctl-${OS}-${ARCH}" ]; then
    echo "Installing from local build..."
    $SUDO cp "./bin/cfctl-${OS}-${ARCH}" "${INSTALL_DIR}/${BINARY_NAME}"
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
elif [ -f "./bin/cfctl" ]; then
    echo "Installing from local build..."
    $SUDO cp "./bin/cfctl" "${INSTALL_DIR}/${BINARY_NAME}"
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${RED}Error: Binary not found. Please build first with: make build${NC}"
    exit 1
fi

echo -e "${GREEN}✓ cfctl installed successfully to ${INSTALL_DIR}/${BINARY_NAME}${NC}"
echo ""
echo "To get started, run:"
echo "  cfctl"
echo ""
echo "For help, run:"
echo "  cfctl --help"
