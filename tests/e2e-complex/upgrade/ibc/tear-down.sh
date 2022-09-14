#!/bin/bash

set -euox pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"

info "Tear down" # ---
docker compose down --timeout 20 --volumes --remove-orphans
