#!/bin/bash

set -euo pipefail

CHEQD_IMAGE_FROM="ghcr.io/cheqd/cheqd-node"
CHEQD_TAG_FROM="0.6.7"

CHEQD_IMAGE_TO="cheqd-node"
CHEQD_TAG_TO="latest"

VOTING_PERIOD="10"
EXPECTED_BLOCK_SECOND="1"
EXTRA_BLOCKS="5"

UPGRADE_NAME="v0.6"
DEPOSIT_AMOUNT="10000000"

CHAIN_ID="cheqd"

GAS="auto"
GAS_ADJUSTMENT="1.3"
GAS_PRICES="25ncheq"

TX_PARAMS="--gas ${GAS} \
    --gas-adjustment ${GAS_ADJUSTMENT} \
    --gas-prices ${GAS_PRICES} \
    --chain-id ${CHAIN_ID} \
    --keyring-backend test \
    -y"
QUERY_PARAMS="--output json"

function set_old_compose_env() {
    export CHEQD_NODE_IMAGE=${CHEQD_IMAGE_FROM}
    export DOCKER_IMAGE_VERSION=${CHEQD_TAG_FROM}
    export NETWORK_EXTERNAL="true"
}

function set_new_compose_env() {
    export CHEQD_NODE_IMAGE=${CHEQD_IMAGE_TO}
    export DOCKER_IMAGE_VERSION=${CHEQD_TAG_TO}
    export NETWORK_EXTERNAL="true"
}
