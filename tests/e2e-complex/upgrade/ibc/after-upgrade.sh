#!/bin/bash

set -euox pipefail

source common.sh

CHEQD_USER_ADDRESS="$(cat cheqd-user-address.txt)"
OSMOSIS_USER_ADDRESS="$(cat osmosis-user-address.txt)"

CHEQD_BALANCE_1=$(cat cheqd-balance-1.txt)
CHEQD_BALANCE_2=$(cat cheqd-balance-2.txt)

DENOM=$(cat denom.txt)


info "Back transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL "$CHEQD_USER_ADDRESS" 10000000000"${DENOM}" --from osmosis-user --chain-id osmosis --keyring-backend test -y
sleep 30 # Wait for relayer

info "Check balances for the last time" # ---
CHEQD_BALANCE_3=$(cd .. && docker-compose exec ${CHEQD_SERVICE} cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
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
