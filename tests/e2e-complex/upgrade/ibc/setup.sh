#!/bin/bash

set -euox pipefail

source common.sh


info "Cleanup"
docker-compose down --volumes --remove-orphans

info "Running osmosis network"
docker-compose up -d osmosis
docker-compose cp osmosis/osmosis_init.sh osmosis:/home/osmosis/osmosis_init.sh
docker-compose exec osmosis bash /home/osmosis/osmosis_init.sh
docker-compose exec -d osmosis osmosisd start

info "Waiting for osmosis chain"
compose_wait_for_chain_height osmosis osmosisd




# TODO: Remove the workaround
export NODE_0_ID="qwe"

info "Create relayer user on cheqd"  # ---
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"
CHEQD_RELAYER_ACCOUNT=$(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json 2>&1)
CHEQD_RELAYER_ADDRESS=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
# TODO: Refactor users
RES=$(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded tx bank send ${CHEQD_USER} "${CHEQD_RELAYER_ADDRESS}" 1000000000000ncheq --gas-prices 25ncheq --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

info "Create relayer user on osmosis" # ---
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(docker-compose exec osmosis osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test 2>&1)
OSMOSIS_RELAYER_ADDRESS=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
RES=$(docker-compose exec osmosis osmosisd tx bank send osmosis-user "${OSMOSIS_RELAYER_ADDRESS}" 1000stake --chain-id osmosis -y --keyring-backend test --output json)
assert_tx_successful "${RES}"

# TODO: Get rid of sleep
sleep 10 # Wait for state


info "Import accounts in hermes" # ---
docker-compose up -d hermes
docker-compose exec hermes hermes keys restore cheqd --mnemonic "$CHEQD_RELAYER_MNEMONIC" --name cheqd-key
docker-compose exec hermes hermes keys restore osmosis --mnemonic "$OSMOSIS_RELAYER_MNEMONIC" --name osmosis-key

info "Open channel" # ---
docker-compose exec hermes bash -c "yes | hermes create channel cheqd --chain-b osmosis --port-a transfer --port-b transfer --new-client-connection"

info "Start hermes" # ---
docker-compose exec -d hermes hermes start


info "Check balances" # ---
CHEQD_USER_ADDRESS=$(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded keys show --address ${CHEQD_USER} --keyring-backend test | sed 's/\r//g')
OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')

CHEQD_BALANCE_1=$(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
docker-compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS"

echo ${CHEQD_USER_ADDRESS} > cheqd-user-address.txt
echo ${OSMOSIS_USER_ADDRESS} > osmosis-user-address.txt

echo ${CHEQD_BALANCE_1} > cheqd-balance-1.txt
