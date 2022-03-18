#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

cheqd_noded_docker() {
    docker run --rm \
        -v "$(pwd):/home/cheqd" \
        --network host \
        -u root \
        -e HOME=/home/cheqd \
        --entrypoint "cheqd-noded" \
        "${CHEQD_IMAGE_TO}" "$@"
}

# Wait for upgrade height
bash ../../tools/wait-for-chain.sh "$UPGRADE_HEIGHT" $((2 * VOTING_PERIOD))

# Stop docker-compose service
docker_compose_down

# Make all the data accessible
make_777

# Start docker-compose with new base image on new version
docker_compose_up "$CHEQD_IMAGE_TO" "$(pwd)"

# Check that upgrade was successful

# Wait for upgrade height
bash ../../tools/wait-for-chain.sh $((UPGRADE_HEIGHT + 2))

CURRENT_VERSION=$(docker run --entrypoint cheqd-noded "$CHEQD_IMAGE_TO" version 2>&1)

if [ "$CURRENT_VERSION" != "$CHEQD_VERSION_TO" ] ; then
     echo "Upgrade to version $CHEQD_VERSION_TO was not successful"
     exit 1
fi

get_addresses
# "To" address was used for sending tokens before upgrade
# shellcheck disable=SC2154
OP_ADDRESS_BEFORE=${addresses[2]}
# "To" address was used for sending tokens after upgrade
OP_ADDRESS_AFTER=${addresses[3]}

# Check balances after tokens sending
check_balance "$OP_ADDRESS_BEFORE"

# Check that did written before upgrade stil exist
check_did "$DID_1"

# Send tokens for checking functionality after upgrade
send_tokens "$OP_ADDRESS_AFTER"

# Send DID after upgrade
send_did_new "$DID_2"

# Check balance after token sending
check_balance "$OP_ADDRESS_AFTER"

# Check that did written before upgrade stil exist
check_did "$DID_2"

# Check that token transaction exists after upgrade too
check_tx_hashes

# Stop docker compose
docker_compose_down

# Clean environment after test
clean_env