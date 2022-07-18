#!/bin/bash

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

KEYS_ARGS="--keyring-backend=test"

CHAIN_ID="gravity"
USER="gravity-user"

# Node
gravity init --chain-id $CHAIN_ID testing

# User
gravity keys add ${USER} ${KEYS_ARGS}

# TODO: Is it correct?
FAKE_ETHER_ADDR="0x63972388a8C03Db64a9f734edE93Eb00abc30d18"
ORCHESTRATOR_ADDRESS=$(gravity keys show ${USER} ${KEYS_ARGS} --output json | jq -r '.address')

# Genesis
gravity add-genesis-account "$(gravity keys show ${USER} -a ${KEYS_ARGS})" 1000000000stake,1000000000valtoken
gravity gentx ${USER} 500000000stake ${FAKE_ETHER_ADDR} ${ORCHESTRATOR_ADDRESS} --keyring-backend=test --chain-id $CHAIN_ID
gravity collect-gentxs

GENESIS="$HOME/.gravity/config/genesis.json"
TMP_GENESIS="$HOME/.gravity/config/tmp_genesis.json"

jq '.app_state["gov"]["voting_params"]["voting_period"]="10s"' "${GENESIS}" > "${TMP_GENESIS}" && mv "${TMP_GENESIS}" "${GENESIS}"

# Config
CONFIG_TOML="$HOME/.gravity/config/config.toml"

sed -i $sed_extension 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "${CONFIG_TOML}"

sed -i $sed_extension 's/timeout_propose = "3s"/timeout_propose = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/g' "${CONFIG_TOML}"
sed -i $sed_extension 's/timeout_commit = "5s"/timeout_commit = "500ms"/g' "${CONFIG_TOML}"
