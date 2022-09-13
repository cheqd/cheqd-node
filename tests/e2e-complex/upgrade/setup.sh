#!/bin/bash

set -euo pipefail

. "../../tools/helpers.sh"
. "common.sh"

# Network configuration
in_localnet_path bash "gen-network-config.sh"

# Docker network
docker network create "${LOCALNET_NETWORK}" || true

# Run network
set_old_compose_env
localnet_compose up -d

# Wait for the network
(cd ${LOCALNET_PATH} && compose_wait_for_chain_height "validator-0" "cheqd-noded")

# Copy keys
VALIDATORS_COUNT=4

for ((i=0 ; i<VALIDATORS_COUNT ; i++))
do
    MONIKER="validator-$i"

    USER="cheqd"
    GROUP="cheqd"
    DOCKER_HOME="/home/cheqd"

    localnet_compose cp network-config/${MONIKER}/keyring-test ${MONIKER}:/home/cheqd/.cheqdnode
    localnet_compose exec -it --user root ${MONIKER} chown -R ${USER}:${GROUP} ${DOCKER_HOME}
done
