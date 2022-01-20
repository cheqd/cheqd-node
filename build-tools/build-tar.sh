#!/bin/bash

set -euox pipefail

if [ -z ${1+x} ]; then
  echo "Binary path must be passed as the first parameter"
fi

BINARY_PATH="$1"
VERSION=$("${BINARY_PATH}" version)
PKG_NAME="cheqd-node"

BUILD_DIR="build"
OUTPUT_DIR="output"

mkdir -p "${BUILD_DIR}"
mkdir -p "${OUTPUT_DIR}"

# Prepare content
PACKAGE_CONTENT="${BUILD_DIR}/tar-package-content"
mkdir -p "${PACKAGE_CONTENT}"
cp "${BINARY_PATH}" "${PACKAGE_CONTENT}"

# Make an archive
TAR_PACKAGE="${OUTPUT_DIR}/${PKG_NAME}_${VERSION}.tar.gz"
tar -cvzf "${TAR_PACKAGE}" -C "${PACKAGE_CONTENT}" "."
