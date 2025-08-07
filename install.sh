#!/bin/bash

set -e

INSTALL_DIR="/usr/local/bin"
TARGET="$INSTALL_DIR/moesic"

echo "Installing Moesic..."

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

FILENAME=""

if [[ "$OS" == "linux" ]]; then
  FILENAME="moesic-linux"
elif [[ "$OS" == "darwin" ]]; then
  if [[ "$ARCH" == "arm64" ]]; then
    FILENAME="moesic-macos-arm64"
  else
    FILENAME="moesic-macos"
  fi
else
  echo "Unsupported OS: $OS"
  exit 1
fi

echo "Fetching latest release info..."
LATEST_TAG=$(curl -s https://api.github.com/repos/angga7togk/moesic/releases/latest | grep tag_name | cut -d '"' -f 4)

if [[ -z "$LATEST_TAG" ]]; then
  echo "Failed to get latest release."
  exit 1
fi

URL="https://github.com/angga7togk/moesic/releases/download/$LATEST_TAG/$FILENAME"

echo "Downloading $URL"
curl -L "$URL" -o moesic
chmod +x moesic

echo "Installing to $INSTALL_DIR (requires sudo)"
sudo mv moesic "$TARGET"

echo "Installed successfully!"
echo "You can now run: moesic"
