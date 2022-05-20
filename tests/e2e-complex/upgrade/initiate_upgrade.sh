#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. common.sh

# Send proposal to pool
# shellcheck disable=SC2086
local_client_tx tx gov submit-proposal software-upgrade \
    "$UPGRADE_NAME" \
    --title "Upgrade-to-new-version" \
    --description "Description-of-upgrade-to-new-version" \
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
