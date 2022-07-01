#!/bin/bash

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
  RES=$1

  if [[ $(echo "${RES}" | jq --raw-output '.code') == 0 ]]
  then
    info "tx successful"
  else
    err "non zero tx return code"
    exit 1
  fi
}


info "Run docker" # ---
docker-compose up --detach --build --force-recreate --remove-orphans
sleep 15 # Wait for chains

info "Create relayer user on cheqd"  # ---
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"
CHEQD_RELAYER_ACCOUNT=$(docker-compose exec cheqd cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --keyring-backend test --output json 2>&1)
CHEQD_RELAYER_ADDRESS=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo "${CHEQD_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
RES=$(docker-compose exec cheqd cheqd-noded tx bank send cheqd-user "${CHEQD_RELAYER_ADDRESS}" 1000000000000ncheq --gas-prices 25ncheq --chain-id cheqd -y --keyring-backend test)
assert_tx_successful "${RES}"

info "Create relayer user on osmosis" # ---
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(docker-compose exec osmosis osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test 2>&1)
OSMOSIS_RELAYER_ADDRESS=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo "${OSMOSIS_RELAYER_ACCOUNT}" | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
RES=$(docker-compose exec osmosis osmosisd tx bank send osmosis-user "${OSMOSIS_RELAYER_ADDRESS}" 1000stake --chain-id osmosis -y --keyring-backend test --output json)
assert_tx_successful "${RES}"
sleep 10 # Wait for state

info "Import accounts in hermes" # ---
docker-compose exec hermes hermes keys restore cheqd --mnemonic "$CHEQD_RELAYER_MNEMONIC" --name cheqd-key
docker-compose exec hermes hermes keys restore osmosis --mnemonic "$OSMOSIS_RELAYER_MNEMONIC" --name osmosis-key

info "Open channel" # ---
docker-compose exec hermes hermes create channel cheqd --chain-b osmosis --port-a transfer --port-b transfer --new-client-connection

info "Start hermes" # ---
docker-compose exec -d hermes hermes start


info "Check balances" # ---
CHEQD_USER_ADDRESS=$(docker-compose exec cheqd cheqd-noded keys show --address cheqd-user --keyring-backend test | sed 's/\r//g')
OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')

CHEQD_BALANCE_1=$(docker-compose exec cheqd cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
docker-compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS"

info "Forward transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec cheqd cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL "$OSMOSIS_USER_ADDRESS" 10000000000ncheq --from cheqd-user --chain-id cheqd --gas-prices 25ncheq --keyring-backend test -y
sleep 30 # Wait for relayer

info "Check balances the second time" # ---
CHEQD_BALANCE_2=$(docker-compose exec cheqd cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
BALANCES=$(docker-compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS" --output json)

info "Denom trace" # ---
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_CUT=$(echo "$DENOM" | cut -c 5-)
docker-compose exec osmosis osmosisd query ibc-transfer denom-trace "$DENOM_CUT"

info "Back transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL "$CHEQD_USER_ADDRESS" 10000000000"${DENOM}" --from osmosis-user --chain-id osmosis --keyring-backend test -y
sleep 30 # Wait for relayer

info "Check balances the last time" # ---
CHEQD_BALANCE_3=$(docker-compose exec cheqd cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
docker-compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS"

CHEQD_BALANCE_1=$(echo "$CHEQD_BALANCE_1" | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_2=$(echo "$CHEQD_BALANCE_2" | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_3=$(echo "$CHEQD_BALANCE_3" | jq --raw-output '.balances[0].amount')

info "Assert balances" # ---
if [[ $CHEQD_BALANCE_2 < $CHEQD_BALANCE_1 ]]
then
  info "cheqd -> osmosis transfer is successfull"
else
  err "cheqd -> osmosis transfer error"
  exit 1
fi

if [[ $CHEQD_BALANCE_3 > $CHEQD_BALANCE_2 ]]
then
  info "osmosis -> cheqd transfer is successfull"
else
  err "osmosis -> cheqd transfer error"
  exit 1
fi

if [[ $CHEQD_BALANCE_3 < $CHEQD_BALANCE_1 ]]
then
  info "fee processed successfully"
else
  err "fee error"
  exit 1
fi


info "Tear down" # ---
docker-compose down --timeout 20
