#!/usr/bin/env bash

# UVE Installer
set -euo pipefail

BASE_URL="https://github.com/iamshreeram/uve/releases/latest/download"
TEMP_DIR=$(mktemp -d)

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64)  ARCH="arm64" ;;
    *)      echo "Unsupported architecture"; exit 1 ;;
esac

# Download appropriate binary
echo "Downloading UVE for ${OS}-${ARCH}..."
case $OS in
    linux|darwin)
        curl -L "${BASE_URL}/uve-${OS}-${ARCH}.tar.gz" -o "${TEMP_DIR}/uve.tar.gz"
        tar -xzf "${TEMP_DIR}/uve.tar.gz" -C "${TEMP_DIR}"
        BIN_PATH="${HOME}/.local/bin"
        mkdir -p "${BIN_PATH}"
        mv "${TEMP_DIR}/uve-bin" "${BIN_PATH}/uve"
        chmod +x "${BIN_PATH}/uve"
        ;;
    *)
        echo "Unsupported OS"; exit 1 ;;
esac

# Initialize shell
echo "Setting up shell integration..."
"${BIN_PATH}/uve" init

echo "UVE installed successfully to ${BIN_PATH}"
echo "Please restart your shell or run: source ~/.bashrc (or equivalent)"