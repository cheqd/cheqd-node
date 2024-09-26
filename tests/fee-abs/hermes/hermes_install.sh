#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

# Check if the version is provided as an argument
# if [ -z "$1" ]; then
#   echo "Usage: $0 <version>"
#   echo "Example: $0 v1.0.0"
#   exit 1
# fi

VERSION="v1.8.0"
BINARY_NAME="hermes"

# Set the URL for the binary based on the provided version
URL="https://github.com/informalsystems/hermes/releases/download/$VERSION/hermes-$VERSION-x86_64-unknown-linux-gnu.tar.gz"
# https://github.com/informalsystems/hermes/releases/download/v1.8.0/hermes-v1.8.0-x86_64-unknown-linux-gnu.tar.gz

# Download the binary tarball
curl -LO $URL

# Extract the binary
tar -xzf hermes-$VERSION-x86_64-unknown-linux-gnu.tar.gz

# Move the binary to /usr/local/bin
sudo mv $BINARY_NAME /usr/local/bin/

# Set executable permissions
sudo chmod +x /usr/local/bin/$BINARY_NAME

# Verify the installation
$BINARY_NAME --version

# Cleanup the downloaded files
rm hermes-$VERSION-x86_64-unknown-linux-gnu.tar.gz

echo "Hermes $VERSION installed successfully!"
