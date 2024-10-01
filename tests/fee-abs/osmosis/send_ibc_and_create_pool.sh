#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function info() {
    printf "${GREEN}[info] %s${NC}\n" "${1}"
}

function err() {
    printf "${RED}[err] %s${NC}\n" "${1}"
}

function assert_tx_successful() {
  RES="$1"

  if [[ $(echo "${RES}" | jq --raw-output '.code') == 0 ]]
  then
    info "tx successful"
  else
    err "non zero tx return code"
    exit 1
  fi
}

CHEQD_USER_ADDRESS=$(cheqd-noded keys show --address cheqd-user --keyring-backend test | tr -d '\r')
OSMOSIS_USER_ADDRESS=$(osmosisd keys show --address osmosis-user --keyring-backend test | tr -d '\r')

info "Transfer cheqd -> osmosis" # ---
PORT="transfer"
CHANNEL="channel-0"
cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL "$OSMOSIS_USER_ADDRESS" 10000000000ncheq --from cheqd-user --chain-id cheqd --gas-prices 50ncheq --keyring-backend test -y
sleep 30 # Wait for relayer

info "Get balances" # ---
CHEQD_BALANCE_2=$(cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
BALANCES=$(osmosisd query bank balances "$OSMOSIS_USER_ADDRESS" --output json)

info "Denom trace" # ---
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_CUT=$(echo "$DENOM" | cut -c 5-)
osmosisd query ibc-transfer denom-trace "$DENOM_CUT"

info "Send 100OSMO to cheqd"
osmosisd tx ibc-transfer transfer $PORT $CHANNEL "$CHEQD_USER_ADDRESS" 100000000uosmo --from osmosis-user --chain-id osmosis --fees 500uosmo --keyring-backend test -y
sleep 30

CHEQD_BALANCE_2=$(cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
DENOM=$(echo "$CHEQD_BALANCE_2" | jq --raw-output '.balances[0].denom')


info "create pool"
# create pool
TX_HASH=$(osmosisd tx gamm create-pool --pool-file pool.json --from $OSMOSIS_USER_ADDRESS --keyring-backend test  --gas-prices 1uosmo --gas-adjustment 1 -y --chain-id osmosis --output json --gas 350000 | jq -r '.txhash')
echo "tx hash: $TX_HASH"
sleep 5

POOL_ID=$(osmosisd q tx $TX_HASH --output json | jq -r '.logs[0].events[-10].attributes[-1].value')
echo "pool id: $POOL_ID"

info "enable fee abs"
cheqd-noded tx gov submit-legacy-proposal param-change proposal.json --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 50ncheq --gas 350000
sleep 5 
cheqd-noded tx gov vote 1 yes --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 50ncheq
sleep 5 

info "add host zone config"
cheqd-noded tx gov submit-legacy-proposal add-hostzone-config host_zone.json --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 50ncheq --gas 350000
sleep 5
cheqd-noded tx gov vote 2 yes --from $CHEQD_USER_ADDRESS --keyring-backend test --chain-id cheqd --yes --gas-prices 50ncheq
