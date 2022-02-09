set -euox pipefail

BINARY_NAME=${PKG_NAME}d
PATH_TO_BIN=/home/runner/go/bin/${BINARY_NAME} # TODO: Must be parameters
TMP_DIR=usr/bin
OUTPUT_DIR=output
TAR_ARCHIVE=${PKG_NAME}_${VERSION}.tar.gz
PATH_TAR=$OUTPUT_DIR/$TAR_ARCHIVE
