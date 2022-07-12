#!/bin/bash

set -euox pipefail

# shellcheck disable=SC1091
. "../../tools/helpers.sh"
. "common.sh"

set_old_compose_env

CURRENT_HEIGHT=$(cd ${LOCALNET_PATH} && docker-compose exec validator-0 cheqd-noded status 2>&1 | jq -r '.SyncInfo.latest_block_height')
VOTING_END_HEIGHT=$((CURRENT_HEIGHT + VOTING_PERIOD / EXPECTED_BLOCK_SECOND + EXTRA_BLOCKS))
UPGRADE_HEIGHT=$((CURRENT_HEIGHT + VOTING_PERIOD / EXPECTED_BLOCK_SECOND + EXTRA_BLOCKS * 2))

set_old_compose_env

# Send proposal
RES=$(localnet_compose exec validator-0 \
    cheqd-noded tx gov submit-proposal software-upgrade \
    "$UPGRADE_NAME" \
    --title "Upgrade title" \
    --description "Upgrade description" \
    --upgrade-height "$UPGRADE_HEIGHT" \
    --upgrade-info "Upgrade info" \
    --from operator-0 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

# Set the deposit from operator0
RES=$(localnet_compose exec validator-0 \
    cheqd-noded tx gov deposit 1 \
    "${DEPOSIT_AMOUNT}ncheq" \
    --from operator-0 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

# Make a vote for operator0
RES=$(localnet_compose exec validator-0 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-0 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

# Make a vote for operator1
RES=$(localnet_compose exec validator-1 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-1 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

# Make a vote for operator2
RES=$(localnet_compose exec validator-2 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-2 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

# Make a vote for operator3
RES=$(localnet_compose exec validator-3 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-3 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"


# Wait for the end of voting
(cd ${LOCALNET_PATH} && compose_wait_for_chain_height "validator-0" "cheqd-noded" "$VOTING_END_HEIGHT")

# TODO: Check that the proposal is accepted
STATUS=$(localnet_compose exec validator-0 cheqd-noded query gov proposal 1 --output json | jq -r '.status')
assert_eq "${STATUS}" "PROPOSAL_STATUS_PASSED"

# Wait for upgrade height
(cd ${LOCALNET_PATH} && compose_wait_for_chain_height "validator-0" "cheqd-noded" "$UPGRADE_HEIGHT")

# Shut down network
(cd ${LOCALNET_PATH} && docker-compose down)

# Bump node version
set_new_compose_env

# Restart network
(cd ${LOCALNET_PATH} && docker-compose up -d)

# Wait for upgrade height + 2
(cd ${LOCALNET_PATH} && compose_wait_for_chain_height "validator-0" "cheqd-noded" "$((UPGRADE_HEIGHT + 2))")

echo "Upgrade successfull"
