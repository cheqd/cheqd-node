#!/bin/bash

set -euox pipefail

source common.sh
source ../common.sh


CHEQD_USER_ADDRESS="$(cat cheqd-user-address.txt)"
OSMOSIS_USER_ADDRESS="$(cat osmosis-user-address.txt)"

set_old_compose_env

info "Forward transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL "$OSMOSIS_USER_ADDRESS" 10000000000ncheq --from ${CHEQD_USER} --chain-id cheqd --gas-prices 25ncheq --keyring-backend test -y)
sleep 30 # Wait for relayer

info "Check balances the second time" # ---
CHEQD_BALANCE_2=$(set +x && localnet_compose exec ${CHEQD_SERVICE} cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
BALANCES=$(docker-compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS" --output json)

info "Denom trace" # ---
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_CUT=$(echo "$DENOM" | cut -c 5-)
docker-compose exec osmosis osmosisd query ibc-transfer denom-trace "$DENOM_CUT"


echo "$CHEQD_BALANCE_2" > cheqd-balance-2.txt

echo "$DENOM" > denom.txt
