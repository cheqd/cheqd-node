#!/bin/bash

set -euo pipefail

. "../../tools/helpers.sh"
. "common.sh"

# Shut down the network
set_new_compose_env
localnet_compose down --volumes --remove-orphans

# Remove configuration
in_localnet_path rm -rf "network-config"

# Remove docker network
docker network remove "${LOCALNET_NETWORK}" || true
