#!/bin/bash

set -euox pipefail

source common.sh
. ../common.sh

CHEQD_USER_ADDRESS="$(cat cheqd-user-address.txt)"
OSMOSIS_USER_ADDRESS="$(cat osmosis-user-address.txt)"

CHEQD_BALANCE_1=$(cat cheqd-balance-1.txt)
CHEQD_BALANCE_2=$(cat cheqd-balance-2.txt)

DENOM=$(cat denom.txt)

set_new_compose_env

info "Check balances for the pre last time" # ---
CHEQD_BALANCE_3=$(set +x && localnet_compose exec ${CHEQD_SERVICE} cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)

info "Back transfer" # ---
PORT="transfer"
CHANNEL="channel-0"
# read -n 1 -p "Press any key to continue..."
docker compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL "$CHEQD_USER_ADDRESS" 10000000000"${DENOM}" --from osmosis-user --chain-id osmosis --keyring-backend test -y

# TODO: Fix the bug with time shift
sleep $((5*60)) # Wait for relayer

info "Check balances for the last time" # ---
CHEQD_BALANCE_4=$(set +x && localnet_compose exec ${CHEQD_SERVICE} cheqd-noded query bank balances "$CHEQD_USER_ADDRESS" --output json)
docker compose exec osmosis osmosisd query bank balances "$OSMOSIS_USER_ADDRESS"

CHEQD_BALANCE_1=$(echo "$CHEQD_BALANCE_1" | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_2=$(echo "$CHEQD_BALANCE_2" | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_3=$(echo "$CHEQD_BALANCE_3" | jq --raw-output '.balances[0].amount')
CHEQD_BALANCE_4=$(echo "$CHEQD_BALANCE_4" | jq --raw-output '.balances[0].amount')

echo "$CHEQD_BALANCE_1" > cheqd-balance-1.txt
echo "$CHEQD_BALANCE_2" > cheqd-balance-2.txt
echo "$CHEQD_BALANCE_3" > cheqd-balance-3.txt
echo "$CHEQD_BALANCE_4" > cheqd-balance-4.txt

info "Assert balances" # ---
if [[ $CHEQD_BALANCE_2 < $CHEQD_BALANCE_1 ]]
then
  info "cheqd -> osmosis transfer is successfull"
else
  err "cheqd -> osmosis transfer error"
  exit 1
fi

if [[ $CHEQD_BALANCE_4 > $CHEQD_BALANCE_3 ]]
then
  info "osmosis -> cheqd transfer is successfull"
else
  err "osmosis -> cheqd transfer error"
  exit 1
fi
