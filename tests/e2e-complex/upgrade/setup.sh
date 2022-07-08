#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh





# Generate config files
bash gen_node_configs.sh

# Add all needed permissions
make_775

# Start the network on version which will be upgraded from
docker_compose_up "${CHEQD_IMAGE_FROM}" "$(pwd)"

# Wait for start ordering, till height 1
bash ../../tools/wait-for-chain.sh 1
