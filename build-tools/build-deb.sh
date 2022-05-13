#!/bin/bash

set -euox pipefail

if [ -z ${1+x} ]; then
  echo "Binary path must be passed as the first parameter"
fi

 if [ ! -z "${2}" ]; then
  COSMOVISOR_BIN=${2}
 fi

BINARY_PATH="$1"
# It's needed for creating an RC deb package
if [ -z ${2+x} ];
then
  VERSION=$("${BINARY_PATH}" version 2>&1)
else
  VERSION="$2"
fi

PKG_NAME="cheqd-node"

BUILD_DIR="build"
OUTPUT_DIR="output"

mkdir -p "${BUILD_DIR}"
mkdir -p "${OUTPUT_DIR}"

# Prepare content
PACKAGE_CONTENT="${BUILD_DIR}/deb-package-content"
mkdir -p "${PACKAGE_CONTENT}/usr/bin/"
cp "${BINARY_PATH}" "${PACKAGE_CONTENT}/usr/bin/"

 if [ ! -z "${COSMOVISOR_BIN}" ]; then
  cp "${COSMOVISOR_BIN}" "${PACKAGE_CONTENT}/usr/bin/"
 fi

# Make intermediate archive
PACKAGE_CONTENT_TAR="${BUILD_DIR}/deb-package-content.tar.gz"
tar -cvzf "${PACKAGE_CONTENT_TAR}" -C "${PACKAGE_CONTENT}" "."

# Build deb based on the archive
ARCH="amd64"
DEB_PACKAGE="${OUTPUT_DIR}/${PKG_NAME}_${VERSION}_${ARCH}.deb"

fpm \
  --input-type "tar" \
  --output-type "deb" \
  --version "${VERSION}" \
  --name "cheqd-node" \
  --description "cheqd node" \
  --url "https://github.com/cheqd/cheqd-node" \
  --architecture "${ARCH}" \
  --deb-generate-changes \
  --deb-compression gz \
  --after-install "postinst" \
  --deb-after-purge "postpurge" \
  --deb-systemd-enable \
  --deb-systemd cheqd-noded.service \
  --depends "logrotate" \
  --verbose \
  --package "${DEB_PACKAGE}" \
  "${PACKAGE_CONTENT_TAR}"
