#!/bin/bash

set -euox pipefail

. "../../tools/helpers.sh"
. "common.sh"

# Shut down the network
set_new_compose_env
localnet_compose down --volumes --remove-orphans

# Remove configuration
(cd ${LOCALNET_PATH} && rm -rf "network-config")

# Remove docker network
docker network remove "${LOCALNET_NETWORK}" || true
