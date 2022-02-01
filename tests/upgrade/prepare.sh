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

# Wait for start ordering, till height 1
bash ../networks/tools/wait-for-chain.sh 1

# Send tokens before upgrade
send_tokens

# Send DID transactions
send_did

# Check that token transaction exists
check_tx_hashes
