#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh





# Generate config files
bash gen_node_configs.sh

# Add all needed permissions
make_775

# Network setup
docker network create ${NETWORK_NAME}

# Start the network on version which will be upgraded from
docker_compose_up "${CHEQD_IMAGE_FROM}" "$(pwd)"


# TODO: Remove the workaround
docker-compose cp node_configs/client/.cheqdnode/keyring-test node0:/home/cheqd/.cheqdnode


# Wait for start ordering, till height 1
# bash ../../tools/wait-for-chain.sh 1
compose_wait_for_chain_height node0 cheqd-noded
