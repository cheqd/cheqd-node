#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

# Stop docker compose
docker_compose_down

# Clean environment (for reproducable purposes in future)
clean_env

# Generate config files
bash gen_node_configs.sh

# Start the network on version which will be upgraded from
docker_compose_up "${CHEQD_IMAGE_FROM}" "$(pwd)"

# Wait for start ordering, till height 1
bash ../../tools/wait-for-chain.sh 1

# Get address of operator which will be used for sending tokens before upgrade
get_addresses
# shellcheck disable=SC2154
OP2_ADDRESS=${addresses[2]}

# Send tokens before upgrade
send_tokens "$OP2_ADDRESS"

# Send DID transaction
send_did "$DID_1"

sleep 5

# Check that token transaction exists
check_tx_hashes

# Check that $DID was written
check_did "$DID_1"

# Check balance after token sending
check_balance "$OP2_ADDRESS"
