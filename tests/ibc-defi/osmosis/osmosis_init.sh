#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

CHAIN_ID="osmosis"

# Node
osmosisd init --chain-id $CHAIN_ID testing

# User
osmosisd keys add osmosis-user --keyring-backend=test

# Genesis
osmosisd add-genesis-account $(osmosisd keys show osmosis-user -a --keyring-backend=test) 1000000000stake,1000000000valtoken
osmosisd gentx osmosis-user 500000000stake --keyring-backend=test --chain-id $CHAIN_ID
osmosisd collect-gentxs

cat $HOME/.osmosisd/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="10s"' > $HOME/.osmosisd/config/tmp_genesis.json && mv $HOME/.osmosisd/config/tmp_genesis.json $HOME/.osmosisd/config/genesis.json

# Config
sed -i $sed_extension 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.osmosisd/config/config.toml"
