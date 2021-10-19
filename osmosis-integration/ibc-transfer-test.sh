#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi


# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function info() {
    printf "${GREEN}[info] ${1}${NC}\n"
}

function err() {
    printf "${RED}[err] ${1}${NC}\n"
}


info "Run docker" # ---
docker-compose up --detach --build --force-recreate --remove-orphans
sleep 15 # Wait for chains

info "Create relayer user on cheqd"  # ---
CHEQD_RELAYER_KEY_NAME="cheqd-relayer"
CHEQD_RELAYER_ACCOUNT=$(docker-compose exec cheqd cheqd-noded keys add ${CHEQD_RELAYER_KEY_NAME} --output json)
CHEQD_RELAYER_ADDRESS=$(echo ${CHEQD_RELAYER_ACCOUNT} | jq --raw-output '.address')
CHEQD_RELAYER_MNEMONIC=$(echo ${CHEQD_RELAYER_ACCOUNT} | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
RES=$(docker-compose exec cheqd cheqd-noded tx bank send cheqd-user ${CHEQD_RELAYER_ADDRESS} 1000000000000ncheq --gas-prices 25ncheq --chain-id cheqd -y)
[[ $(echo $RES | jq --raw-output '.code') == 0 ]] && info "tx successfull" || (err "non zero tx return code"; exit 1)


info "Create relayer user on osmosis" # ---
OSMOSIS_RELAYER_KEY_NAME="osmosis-relayer"
OSMOSIS_RELAYER_ACCOUNT=$(docker-compose exec osmosis osmosisd keys add ${OSMOSIS_RELAYER_KEY_NAME} --output json --keyring-backend test)
OSMOSIS_RELAYER_ADDRESS=$(echo ${OSMOSIS_RELAYER_ACCOUNT} | jq --raw-output '.address')
OSMOSIS_RELAYER_MNEMONIC=$(echo ${OSMOSIS_RELAYER_ACCOUNT} | jq --raw-output '.mnemonic')

info "Send some tokens to it" # ---
RES=$(docker-compose exec osmosis osmosisd tx bank send osmosis-user ${OSMOSIS_RELAYER_ADDRESS} 1000stake --chain-id osmosis -y --keyring-backend test --output json)
[[ $(echo $RES | jq --raw-output '.code') == 0 ]] && info "tx successfull" || (err "non zero tx return code"; exit 1)
sleep 10 # Wait for state

info "Import accounts in hermes" # ---
docker-compose exec hermes hermes keys restore cheqd --mnemonic "$CHEQD_RELAYER_MNEMONIC" --name cheqd-key
docker-compose exec hermes hermes keys restore osmosis --mnemonic "$OSMOSIS_RELAYER_MNEMONIC" --name osmosis-key

info "Open channel" # ---
docker-compose exec hermes hermes create channel cheqd osmosis --port-a transfer --port-b transfer

info "Start hermes" # ---
docker-compose exec -d hermes hermes start


info "Check balances" # ---
CHEQD_USER_ADDRESS=$(docker-compose exec cheqd cheqd-noded keys show --address cheqd-user | sed 's/\r//g')
OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')

CHEQD_BALANCE_1=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)
docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS

info "Forward transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec cheqd cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL $OSMOSIS_USER_ADDRESS 10000000000ncheq --from cheqd-user --chain-id cheqd --gas-prices 25ncheq -y
sleep 30 # Wait for relayer

info "Check balances the second time" # ---
CHEQD_BALANCE_2=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)
BALANCES=$(docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS --output json)

log "Denom trace" # ---
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_CUT=$(echo "$DENOM" | cut -c 5-)
docker-compose exec osmosis osmosisd query ibc-transfer denom-trace $DENOM_CUT

info "Back transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL $CHEQD_USER_ADDRESS 10000000000${DENOM} --from osmosis-user --chain-id osmosis --keyring-backend test -y
sleep 30 # Wait for relayer

info "Check balances the last time" # ---
CHEQD_BALANCE_3=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)
docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS

CHEQD_BALANCE_1=$(echo $CHEQD_BALANCE_1 | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_2=$(echo $CHEQD_BALANCE_2 | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_3=$(echo $CHEQD_BALANCE_3 | jq --raw-output '.balances[0].amount')

info "Assert balances" # ---
[[ $CHEQD_BALANCE_2 < $CHEQD_BALANCE_1 ]] && info "cheqd -> osmosis transfer is successfull" || (err "cheqd -> osmosis transfer error"; exit 1)
[[ $CHEQD_BALANCE_3 > $CHEQD_BALANCE_2 ]] && info "osmosis -> cheqd transfer is successfull" || (err "osmosis -> cheqd transfer error"; exit 1)
[[ $CHEQD_BALANCE_3 < $CHEQD_BALANCE_1 ]] && info "fee processed successfully" || (err "fee error"; exit 1)


log "Tear down" # ---
docker-compose down --timeout 20


# Ready:
# - Positive case test
# ToDo:
# - Test with gravity
# - Look at osmosis, atom, gravity genesis params
# - Test back transfers
# - Read white paper
# Questions:
# - Back transfer via other channel?
# - relayers scalability?
# - What to backup on relyer?
