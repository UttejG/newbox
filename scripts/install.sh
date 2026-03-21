#!/usr/bin/env bash
set -euo pipefail

REPO="UttejG/newbox"
BINARY="newbox"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info()    { echo -e "${BLUE}[newbox]${NC} $*"; }
success() { echo -e "${GREEN}[newbox]${NC} $*"; }
warn()    { echo -e "${YELLOW}[newbox]${NC} $*"; }
error()   { echo -e "${RED}[newbox]${NC} $*" >&2; exit 1; }

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) error "Unsupported architecture: $ARCH" ;;
esac
case "$OS" in
  darwin|linux) ;;
  *) error "Unsupported OS: $OS. Use install.ps1 for Windows." ;;
esac

# Get latest release version
info "Fetching latest release..."
if command -v curl &>/dev/null; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
elif command -v wget &>/dev/null; then
  VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
else
  error "curl or wget required"
fi

[ -z "$VERSION" ] && error "Could not determine latest version"
info "Installing newbox ${VERSION} (${OS}/${ARCH})..."

# Download URL
FILENAME="newbox_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

# Download to temp dir
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

if command -v curl &>/dev/null; then
  curl -fsSL "$URL" -o "$TMP/$FILENAME"
  curl -fsSL "https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt" -o "$TMP/checksums.txt"
else
  wget -qO "$TMP/$FILENAME" "$URL"
  wget -qO "$TMP/checksums.txt" "https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
fi

# Verify checksum (anchored match; errors on missing or ambiguous entry)
info "Verifying checksum..."
if command -v sha256sum &>/dev/null; then
  (
    cd "$TMP"
    pattern="^[0-9a-fA-F]{64}  ${FILENAME}$"
    match_count=$(grep -cE "$pattern" checksums.txt || true)
    if [ "$match_count" -eq 0 ]; then error "No checksum entry found for $FILENAME"; fi
    if [ "$match_count" -gt 1 ]; then error "Multiple checksum entries found for $FILENAME"; fi
    grep -E "$pattern" checksums.txt | sha256sum -c --status
  ) || error "Checksum verification failed!"
elif command -v shasum &>/dev/null; then
  (
    cd "$TMP"
    pattern="^[0-9a-fA-F]{64}  ${FILENAME}$"
    match_count=$(grep -cE "$pattern" checksums.txt || true)
    if [ "$match_count" -eq 0 ]; then error "No checksum entry found for $FILENAME"; fi
    if [ "$match_count" -gt 1 ]; then error "Multiple checksum entries found for $FILENAME"; fi
    grep -E "$pattern" checksums.txt | sed 's/  / */' | shasum -a 256 -c --status
  ) || error "Checksum verification failed!"
else
  warn "sha256sum/shasum not found; skipping checksum verification"
fi

# Extract
tar -xzf "$TMP/$FILENAME" -C "$TMP"

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
  chmod +x "$INSTALL_DIR/$BINARY"
else
  info "Requesting sudo to install to $INSTALL_DIR..."
  sudo mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
  sudo chmod +x "$INSTALL_DIR/$BINARY"
fi

success "newbox ${VERSION} installed to ${INSTALL_DIR}/${BINARY}"
success "Run 'newbox' to get started!"
