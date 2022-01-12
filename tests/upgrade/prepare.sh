#!/bin/bash

set -euox pipefail

. common.sh

# Stop docker compose
docker_compose_down

# Clean environment (for reproducable purposes in future)
clean_env

# Generate config files
bash gen_node_configs.sh

# Start the network on version which will be upgraded from
docker_compose_up $CHEQD_VERSION_FROM $(pwd)