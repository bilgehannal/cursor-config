#!/bin/bash
set -euo pipefail

REPO="bilgehannal/cursor-config"
BINARY="curset"
INSTALL_DIR="/usr/local/bin"

echo "Installing curset..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is required but not installed."
    echo "Install Go from https://go.dev/dl/ and try again."
    exit 1
fi

# Clone, build, install
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

git clone --depth 1 "https://github.com/${REPO}.git" "$TMPDIR/repo" 2>/dev/null
cd "$TMPDIR/repo/curset"
go mod tidy
go build -ldflags "-X main.version=latest" -o "$TMPDIR/$BINARY" .
sudo mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"

echo "curset installed to $INSTALL_DIR/$BINARY"
echo "Run 'curset list' to see available collections."
