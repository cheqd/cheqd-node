#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

# Send proposal to pool
# shellcheck disable=SC2086
local_client_tx tx gov submit-proposal software-upgrade \
    "$UPGRADE_NAME" \
    --title "Upgrade-to--the-new-version" \
    --description "Description-of-the-upgrade-to-the-new-version" \
    --upgrade-height "$UPGRADE_HEIGHT" \
    --upgrade-info "$UPGRADE_INFO" \
    --from operator1 \
    ${TX_PARAMS}

# Set the deposit from operator0
# shellcheck disable=SC2086
local_client_tx tx gov deposit 1 \
    "${DEPOSIT_AMOUNT}ncheq" \
    --from operator0 \
    ${TX_PARAMS}

# Make a vote for operator0
# shellcheck disable=SC2086
local_client_tx tx gov vote 1 \
    yes \
    --from operator0 \
    ${TX_PARAMS}

# Make a vote for operator1
# shellcheck disable=SC2086
local_client_tx tx gov vote 1 \
    yes \
    --from operator1 \
    ${TX_PARAMS}



# End of voting


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
bash ../../tools/wait-for-chain.sh "$UPGRADE_HEIGHT" $((3 * VOTING_PERIOD))

# Stop docker-compose service
docker_compose_down

# Make all the data accessible
make_775

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
