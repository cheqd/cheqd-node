#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

CHAIN_ID="osmosis"
USER="osmosis-user"

# Node
osmosisd init --chain-id $CHAIN_ID testing

# User
osmosisd keys add ${USER} --keyring-backend=test

# Genesis
osmosisd add-genesis-account "$(osmosisd keys show ${USER} -a --keyring-backend=test)" 1000000000stake,1000000000valtoken
osmosisd gentx ${USER} 500000000stake --keyring-backend=test --chain-id $CHAIN_ID
osmosisd collect-gentxs

GENESIS="$HOME/.osmosisd/config/genesis.json"
TMP_GENESIS="$HOME/.osmosisd/config/tmp_genesis.json"

jq '.app_state["gov"]["voting_params"]["voting_period"]="10s"' "${GENESIS}" > "${TMP_GENESIS}" && mv "${TMP_GENESIS}" "${GENESIS}"

# Config
CONFIG_TOML="$HOME/.osmosisd/config/config.toml"

sed -i $sed_extension 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "${CONFIG_TOML}"

sed -i $sed_extension 's/timeout_propose = "3s"/timeout_propose = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_commit = "5s"/timeout_commit = "500ms"/g' "${CONFIG_TOML}"
