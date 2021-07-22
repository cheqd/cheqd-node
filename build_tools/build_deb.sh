#!/bin/bash

if [ -n "$1" ]; then
    PKG_NAME=$1
else
    echo "It seems that parameter 'PKG_NAME' was missed. Try: "
    echo "$0 <package name> <version of deb>"
    exit 1
fi

if [ -n "$2" ]; then
    VERSION=$2
else
    echo "It seems that parameter 'VERSION' was missed. Try: "
    echo "$0 <package name> <version of deb>"
    exit 1
fi

source common.sh

ARCH="amd64"
FULL_PKG_NAME=${PKG_NAME}_${VERSION}_${ARCH}.deb
PKG_PATH=$OUTPUT_DIR/$FULL_PKG_NAME

fpm \
    --input-type "tar" \
    --output-type "deb" \
    --version "${VERSION}" \
    --name "cheqd-node" \
    --description "cheqd node" \
    --architecture "${ARCH}" \
    --pre-install "postinst" \
    --after-remove "postremove" \
    --verbose \
    --package "${PKG_PATH}" \
    $PATH_TAR
