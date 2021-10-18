#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi


# Create Alice user on cheqd
ALICE_KEY_NAME=$(echo $RANDOM | base64 | head -c 20)
ALICE_ACCOUNT=$(docker-compose exec cheqd cheqd-noded keys add ${ALICE_KEY_NAME} --output json)
ALICE_ADDRESS=$(echo ${ALICE_ACCOUNT} | jq --raw-output '.address')
ALICE_MNEMONIC=$(echo ${ALICE_ACCOUNT} | jq --raw-output '.mnemonic')

# Send some tokens to her
RES=$(docker-compose exec cheqd cheqd-noded tx bank send cheqd-user ${ALICE_ADDRESS} 1000000000000ncheq --gas-prices 25ncheq --chain-id cheqd -y)
[[ $(echo $RES | jq --raw-output '.code') == 0 ]] && echo "tx successfull" || (echo "non zero tx return code"; exit 1)


# Create Bob user on osmosis
BOB_KEY_NAME=$(echo $RANDOM | base64 | head -c 20)
BOB_ACCOUNT=$(docker-compose exec osmosis osmosisd keys add ${BOB_KEY_NAME} --output json --keyring-backend test)
BOB_ADDRESS=$(echo ${BOB_ACCOUNT} | jq --raw-output '.address')
BOB_MNEMONIC=$(echo ${BOB_ACCOUNT} | jq --raw-output '.mnemonic')

# Send some tokens to him
RES=$(docker-compose exec osmosis osmosisd tx bank send osmosis-user ${BOB_ADDRESS} 1000valtoken --chain-id osmosis -y --keyring-backend test --output json)
[[ $(echo $RES | jq --raw-output '.code') == 0 ]] && echo "tx successfull" || (echo "non zero tx return code"; exit 1)


# Configure and start hermes
docker-compose exec hermes hermes keys restore cheqd --mnemonic "$ALICE_MNEMONIC" --name cheqd-key
docker-compose exec hermes hermes keys restore osmosis --mnemonic "$BOB_MNEMONIC" --name osmosis-key

docker-compose exec hermes hermes create channel cheqd osmosis --port-a transfer --port-b transfer

docker-compose exec -d hermes hermes start


# balances
CHEQD_USER_ADDRESS=$(docker-compose exec cheqd cheqd-noded keys show --address cheqd-user | sed 's/\r//g')
CHEQD_BALANCE_1=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)

OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')
docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS

# forward transfer
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec cheqd cheqd-noded tx ibc-transfer transfer $PORT $CHANNEL $OSMOSIS_USER_ADDRESS 10000000000ncheq --from cheqd-user --chain-id cheqd --gas-prices 25ncheq -y
sleep 60

# balances
CHEQD_USER_ADDRESS=$(docker-compose exec cheqd cheqd-noded keys show --address cheqd-user | sed 's/\r//g')
CHEQD_BALANCE_2=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)

OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')
BALANCES=$(docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS --output json)

# denom trace
DENOM=$(echo "$BALANCES" | jq --raw-output '.balances[0].denom')
DENOM_CUT=$(echo "$DENOM" | cut -c 5-)
docker-compose exec osmosis osmosisd query ibc-transfer denom-trace $DENOM_CUT

# back transfer
PORT="transfer"
CHANNEL="channel-0"
docker-compose exec osmosis osmosisd tx ibc-transfer transfer $PORT $CHANNEL $CHEQD_USER_ADDRESS 10000000000${DENOM} --from osmosis-user --chain-id osmosis --keyring-backend test -y
sleep 60

# balances
CHEQD_USER_ADDRESS=$(docker-compose exec cheqd cheqd-noded keys show --address cheqd-user | sed 's/\r//g')
CHEQD_BALANCE_3=$(docker-compose exec cheqd cheqd-noded query bank balances $CHEQD_USER_ADDRESS --output json)

OSMOSIS_USER_ADDRESS=$(docker-compose exec osmosis osmosisd keys show --address osmosis-user --keyring-backend test | sed 's/\r//g')
docker-compose exec osmosis osmosisd query bank balances $OSMOSIS_USER_ADDRESS

echo $CHEQD_BALANCE_1 | jq '.balances[0].amount'
echo $CHEQD_BALANCE_2 | jq '.balances[0].amount'
echo $CHEQD_BALANCE_3 | jq '.balances[0].amount'

# TODO:
# - Test with gravity
# - Look at osmosis, atom, gravity genesis params
# - Test back transfers
# - Read white paper
