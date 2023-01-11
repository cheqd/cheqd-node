#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

# sed in MacOS requires extra argument
if [[ "$OSTYPE" == "darwin"* ]]; then
  SED_EXT='.orig'
else 
  SED_EXT=''
fi

CHAIN_ID="cheqd"

# Node
cheqd-noded init node0 --chain-id "$CHAIN_ID"
NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)

# User
cheqd-noded keys add cheqd-user --keyring-backend test

# Config
sed -i $SED_EXT 's|minimum-gas-prices = ""|minimum-gas-prices = "50ncheq"|g' "$HOME/.cheqdnode/config/app.toml"

# shellcheck disable=SC2086
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.cheqdnode/config/config.toml"

# Genesis
GENESIS="$HOME/.cheqdnode/config/genesis.json"
sed -i $SED_EXT 's/"stake"/"ncheq"/' "$GENESIS"

cheqd-noded add-genesis-account cheqd-user 1000000000000000000ncheq --keyring-backend test
cheqd-noded gentx cheqd-user 10000000000000000ncheq --chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY" --keyring-backend test

cheqd-noded collect-gentxs
cheqd-noded validate-genesis
