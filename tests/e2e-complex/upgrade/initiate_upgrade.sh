#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

# Send proposal to pool
local_client_tx tx gov submit-proposal software-upgrade \
    "$UPGRADE_NAME" \
    --title "Upgrade-to-new-version" \
    --description "Description-of-upgrade-to-new-version" \
    --upgrade-height "$UPGRADE_HEIGHT" \
    --from operator1 \
    --gas auto \
    --gas-adjustment 1.2 \
    --gas-prices "25ncheq" \
    --chain-id "$CHAIN_ID" \
    -y

# Set the deposit from operator0
local_client_tx tx gov deposit 1 \
    "${DEPOSIT_AMOUNT}ncheq" \
    --from operator0 \
    --gas auto \
    --gas-adjustment 1.2 \
    --gas-prices 25ncheq \
    --chain-id "$CHAIN_ID" \
    -y

# Make a vote for operator0
local_client_tx tx gov vote 1 \
    yes \
    --from operator0 \
    --gas auto \
    --gas-adjustment 1.3 \
    --gas-prices 25ncheq \
    --chain-id "$CHAIN_ID" \
    -y

# Make a vote for operator1
local_client_tx tx gov vote 1 \
    yes \
    --from operator1 \
    --gas auto \
    --gas-adjustment 1.3 \
    --gas-prices 25ncheq \
    --chain-id "$CHAIN_ID" \
    -y
