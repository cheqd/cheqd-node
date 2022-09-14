#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/../../tools/helpers.sh"

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"


echo "=> Shutting down network"
set_new_compose_env
localnet_compose down --volumes --remove-orphans

echo "=> Removing network configuration"
in_localnet_path rm -rf "network-config"

echo "=> Removing docker network"
docker network remove "${LOCALNET_NETWORK}" || true
