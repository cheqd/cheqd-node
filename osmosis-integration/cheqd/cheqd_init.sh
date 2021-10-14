#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

CHAIN_ID="cheqd"

# Node
cheqd-noded init node0 --chain-id $CHAIN_ID
NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)

# User
cheqd-noded keys add cheqd-user

# Config
sed -i $sed_extension 's|minimum-gas-prices = ""|minimum-gas-prices = "25ncheq"|g' "$HOME/.cheqdnode/config/app.toml"
sed -i $sed_extension 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.cheqdnode/config/app.toml"

# Genesis
GENESIS="$HOME/.cheqdnode/config/genesis.json"
sed -i $sed_extension 's/"stake"/"ncheq"/' $GENESIS

cheqd-noded add-genesis-account cheqd-user 1000000000000000000ncheq
cheqd-noded gentx cheqd-user 10000000000000000ncheq --chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY"

cheqd-noded collect-gentxs
cheqd-noded validate-genesis
