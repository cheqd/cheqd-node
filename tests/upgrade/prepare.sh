#!/bin/bash

set -euox pipefail

. common.sh

# Stop docker compose
docker_compose_down

# Clean environment (for reproducable purposes in future)
clean_env

# Generate config files
bash gen_node_configs.sh

# Make all the data accessable
make_777

# Start the network on version which will be upgraded from
docker_compose_up "${CHEQD_IMAGE_FROM}" $(pwd)