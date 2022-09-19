#!/bin/bash

set -euox pipefail

BASE_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# shellcheck disable=SC1091
. "${BASE_DIR}/common.sh"


# echo "=> Running osmosis network"
# ibc_compose up -d osmosis
# ibc_compose cp "${BASE_DIR}/osmosis/osmosis_init.sh" osmosis:/home/osmosis/osmosis_init.sh
# ibc_compose exec osmosis bash /home/osmosis/osmosis_init.sh
# ibc_compose exec -d osmosis osmosisd start

# echo "=> Waiting for osmosis chain"
# compose_wait_for_chain_height osmosis osmosisd


echo "=> Create relayer user on cheqd"  # ---
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"

set_old_compose_env

localnet_compose exec "${CHEQD_SERVICE}" cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json

exit 1

CHEQD_RELAYER_ACCOUNT=$(localnet_compose exec "${CHEQD_SERVICE}" cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json 2>&1)
CHEQD_RELAYER_ADDRESS=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "=> Send some tokens to it" # ---
# TODO: Refactor users
RES=$(set +x && localnet_compose exec "${CHEQD_SERVICE}" cheqd-noded tx bank send "${CHEQD_USER}" "${CHEQD_RELAYER_ADDRESS}" 1000000000000ncheq --gas-prices 25ncheq --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

echo "=> Create relayer user on osmosis" # ---
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(ibc_compose exec osmosis osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test 2>&1)
OSMOSIS_RELAYER_ADDRESS=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

echo "=> Send some tokens to it" # ---
RES=$(ibc_compose exec osmosis osmosisd tx bank send osmosis-user "${OSMOSIS_RELAYER_ADDRESS}" 1000stake --chain-id osmosis -y --keyring-backend test --output json)
assert_tx_successful "${RES}"


# TODO: Get rid of sleep
sleep 10 # Wait for state


echo "=> Import accounts in hermes" # ---
ibc_compose up -d hermes
ibc_compose exec hermes hermes keys restore cheqd --mnemonic "$CHEQD_RELAYER_MNEMONIC" --name cheqd-key
ibc_compose exec hermes hermes keys restore osmosis --mnemonic "$OSMOSIS_RELAYER_MNEMONIC" --name osmosis-key

echo "=> Open channel" # ---
ibc_compose exec hermes bash -c "hermes create channel cheqd --chain-b osmosis --port-a transfer --port-b transfer --new-client-connection << 'y'"

echo "=> Start hermes" # ---
# read -n 1 -p "Press any key to continue..."
ibc_compose exec -d hermes bash -c "hermes start > log.txt 2>&1"


echo "=> Check balances" # ---
CHEQD_USER_ADDRESS=$(set +x && localnet_compose exec "${CHEQD_SERVICE}" cheqd-noded keys show --address "${CHEQD_USER}" --keyring-backend test | sed 's/\r//g')
OSMOSIS_USER_ADDRESS=$(ibc_compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')

CHEQD_BALANCE_1=$(set +x && localnet_compose exec "${CHEQD_SERVICE}" cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
ibc_compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS"

echo "${CHEQD_USER_ADDRESS}" > cheqd-user-address.txt
echo "${OSMOSIS_USER_ADDRESS}" > osmosis-user-address.txt

echo "${CHEQD_BALANCE_1}" > cheqd-balance-1.txt
