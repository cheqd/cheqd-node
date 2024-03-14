#!/bin/bash
# shellcheck disable=SC2086

set -euox pipefail

# sed in macos requires extra argument

# sed in MacOS requires extra argument
if [[ "$OSTYPE" == "darwin"* ]]; then
  SED_EXT='.orig'
else
  SED_EXT=''
fi

CHAIN_ID="osmosis"

# Node
osmosisd init --chain-id "$CHAIN_ID" testing

# Config
sed -i $SED_EXT 's/"stake"/"uosmo"/' "$HOME/.osmosisd/config/genesis.json"
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.osmosisd/config/config.toml"
sed -i $SED_EXT 's|address = "localhost:9090"|address = "0.0.0.0:9090"|g' "$HOME/.osmosisd/config/app.toml"

osmosisd keys add osmosis-user --keyring-backend=test

# Genesis
osmosisd add-genesis-account "$(osmosisd keys show osmosis-user -a --keyring-backend=test)" 2000000000uosmo
osmosisd gentx osmosis-user 500000000uosmo --keyring-backend=test --chain-id "$CHAIN_ID"
osmosisd collect-gentxs
