#!/bin/bash

set -euo pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/../../tools/helpers.sh"

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"

set_old_compose_env

CURRENT_HEIGHT=$(localnet_compose exec validator-0 cheqd-noded status 2>&1 | jq -r '.SyncInfo.latest_block_height')
VOTING_END_HEIGHT=$((CURRENT_HEIGHT + VOTING_PERIOD / EXPECTED_BLOCK_SECOND + EXTRA_BLOCKS))
UPGRADE_HEIGHT=$((CURRENT_HEIGHT + VOTING_PERIOD / EXPECTED_BLOCK_SECOND + EXTRA_BLOCKS * 2))

set_old_compose_env

echo "=> Sending upgrade proposal"

# shellcheck disable=SC2086
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

echo "=> Setting deposit"

# shellcheck disable=SC2086
RES=$(localnet_compose exec validator-0 \
    cheqd-noded tx gov deposit 1 \
    "${DEPOSIT_AMOUNT}ncheq" \
    --from operator-0 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

echo "=> Making a vote for operator0"

# shellcheck disable=SC2086
RES=$(localnet_compose exec validator-0 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-0 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

echo "=> Making a vote for operator1"

# shellcheck disable=SC2086
RES=$(localnet_compose exec validator-1 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-1 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

echo "=> Making a vote for operator2"

# shellcheck disable=SC2086
RES=$(localnet_compose exec validator-2 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-2 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"

echo "=> Making a vote for operator3"

# shellcheck disable=SC2086
RES=$(localnet_compose exec validator-3 \
    cheqd-noded tx gov vote 1 yes \
    --from operator-3 \
    ${TX_PARAMS})
assert_tx_successful "${RES}"


echo "=> Waiting for the end of voting"
in_localnet_path compose_wait_for_chain_height "validator-0" "cheqd-noded" "$VOTING_END_HEIGHT"

echo "=> Checking that the proposal is accepted"
# shellcheck disable=SC2086
STATUS=$(localnet_compose exec validator-0 cheqd-noded query gov proposal 1 ${QUERY_PARAMS} | jq -r '.status')
assert_eq "${STATUS}" "PROPOSAL_STATUS_PASSED"

echo "=> Waiting for upgrade height"
in_localnet_path compose_wait_for_chain_height "validator-0" "cheqd-noded" "$UPGRADE_HEIGHT"

echo "=> Shutting down network"
localnet_compose down

echo "=> Bumping node version"
set_new_compose_env

echo "=> Restarting network"
localnet_compose up -d

echo "=> Waiting for upgrade height + 2"
in_localnet_path compose_wait_for_chain_height "validator-0" "cheqd-noded" "$((UPGRADE_HEIGHT + 2))"

echo "=> Upgrade successfull"
