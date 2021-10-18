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


# Configure hermes
docker-compose exec hermes hermes keys restore cheqd --mnemonic "$ALICE_MNEMONIC" --name cheqd-key
docker-compose exec hermes hermes keys restore osmosis --mnemonic "$BOB_MNEMONIC" --name osmosis-key


# docker-compose exec hermes rm -f config.toml && touch config.toml
# docker-compose exec hermes rm -rf cheqd && mkdir cheqd
# docker-compose exec hermes rm -rf osmosis && mkdir osmosis

# docker-compose exec hermes hermes -c ./config.toml light add tcp://cheqd:26657 -c chain-a -s ./cheqd -p -y -f
