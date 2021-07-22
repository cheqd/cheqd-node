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


source ./common.sh

mkdir -p output
mkdir -p $TMP_DIR
cp $PATH_TO_BIN $TMP_DIR

tar -czf $PATH_TAR $TMP_DIR


