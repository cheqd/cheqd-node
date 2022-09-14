#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/../../tools/helpers.sh"

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"


echo "=> Genegrating network configuration"
in_localnet_path bash "gen-network-config.sh"

echo "=> Creating docker network"
docker network create "${LOCALNET_NETWORK}" || true

echo "=> Starting network"
set_old_compose_env
localnet_compose up -d

echo "=> Waiting for network to start"
in_localnet_path compose_wait_for_chain_height "validator-0" "cheqd-noded"

echo "=> Copying keys"
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
