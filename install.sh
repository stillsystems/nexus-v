#!/bin/bash
# NEXUS-V Linux/macOS Install Script

echo "Installing NEXUS-V..."

# Determine OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH="amd64"
if [[ "$(uname -m)" == "arm64" || "$(uname -m)" == "aarch64" ]]; then
    ARCH="arm64"
fi

# Download and extract binary
URL="https://github.com/stillsystems/nexus-v/releases/latest/download/nexus-v_${OS}_${ARCH}.tar.gz"
curl -L "$URL" -o nexus-v.tar.gz
tar -xzf nexus-v.tar.gz nexus-v
rm nexus-v.tar.gz

# Make executable and move to bin
chmod +x nexus-v
sudo mv nexus-v /usr/local/bin/

echo "NEXUS-V installed successfully! Run 'nexus-v version' to verify."
