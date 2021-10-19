#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

CHAIN_ID="gravity"

# Node
liquidityd init node0 --chain-id $CHAIN_ID
NODE_0_VAL_PUBKEY=$(liquidityd tendermint show-validator)

# User
liquidityd keys add gravity-user --keyring-backend test

# Config
sed -i $sed_extension 's|minimum-gas-prices = ""|minimum-gas-prices = "0stake"|g' "$HOME/.liquidityapp/config/app.toml"
sed -i $sed_extension 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.liquidityapp/config/config.toml"

# Genesis
GENESIS="$HOME/.cheqdnode/config/genesis.json"

liquidityd add-genesis-account gravity-user 1000000000stake --keyring-backend test
liquidityd gentx gravity-user 1000000000stake --chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY"  --keyring-backend test

liquidityd collect-gentxs
liquidityd validate-genesis
