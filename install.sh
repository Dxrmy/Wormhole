#!/bin/bash
set -e

REPO="Dxrmy/Wormhole"
BRANCH="main"
BASE_URL="https://raw.githubusercontent.com/$REPO/$BRANCH"

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
    linux) OS_NAME="linux" ;;
    darwin) OS_NAME="mac" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH_NAME="amd64" ;;
    arm64|aarch64) ARCH_NAME="arm64" ;;
    i386|i686) ARCH_NAME="386" ;;
    *) echo "Unsupported Architecture: $ARCH"; exit 1 ;;
esac

# Construct binary name
if [ "$OS_NAME" = "mac" ] && [ "$ARCH_NAME" = "amd64" ]; then
    BINARY="proxy-mac-amd64"
elif [ "$OS_NAME" = "mac" ] && [ "$ARCH_NAME" = "arm64" ]; then
    BINARY="proxy-mac-arm64"
else
    BINARY="proxy-${OS_NAME}-${ARCH_NAME}"
fi

URL="$BASE_URL/$BINARY"
DEST_DIR="$HOME/.wormhole"
DEST_FILE="$DEST_DIR/wormhole"

echo "Detected System: $OS ($ARCH)"
echo "Downloading Terraria Proxy..."

mkdir -p "$DEST_DIR"
curl -# -L "$URL" -o "$DEST_FILE"
chmod +x "$DEST_FILE"

echo "Starting Wormhole Proxy..."
"$DEST_FILE" "$@"
