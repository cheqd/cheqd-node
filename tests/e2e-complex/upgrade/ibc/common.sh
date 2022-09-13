#!/bin/bash

set -euox pipefail


BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

. "${BASE_DIR}/../../../tools/helpers.sh"

# TMP
CHEQD_SERVICE="validator-0"
CHEQD_USER="operator-0"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function info() {
    printf "${GREEN}[info] %s${NC}\n" "${1}"
}

function err() {
    printf "${RED}[err] %s${NC}\n" "${1}"
}
