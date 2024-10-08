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

CONFIG_TOML="$HOME/.osmosisd/config/config.toml"
CLIENT_TOML="$HOME/.osmosisd/config/client.toml"

# Config
sed -i $SED_EXT 's/"stake"/"uosmo"/' "$HOME/.osmosisd/config/genesis.json"
sed -i $SED_EXT 's/"1000000000"/"1000000"/' "$HOME/.osmosisd/config/genesis.json"

sed -i 's/"allow_queries": \[\]/"allow_queries": \["\/osmosis.twap.v1beta1.Query\/ArithmeticTwap", "\/osmosis.twap.v1beta1.Query\/ArithmeticTwapToNow", "\/osmosis.twap.v1beta1.Query\/GeometricTwap", "\/osmosis.twap.v1beta1.Query\/GeometricTwapToNow", "\/osmosis.twap.v1beta1.Query\/Params"\]/' "$HOME/.osmosisd/config/genesis.json"

sed -i $SED_EXT 's/"stake"/"uosmo"/' "$HOME/.osmosisd/config/genesis.json"
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.osmosisd/config/config.toml"
sed -i $SED_EXT 's|address = "localhost:9090"|address = "0.0.0.0:9090"|g' "$HOME/.osmosisd/config/app.toml"

sed -i $SED_EXT 's/timeout_propose = "3s"/timeout_propose = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_commit = "3s"/timeout_commit = "2s"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/output = "text"/output = "json"/g' "${CLIENT_TOML}"

# sed -i $SED_EXT 's/timeout_commit = "3s"/timeout_commit = "2s"/g' "${CONFIG_TOML}"

osmosisd keys add osmosis-user --keyring-backend=test

# Genesis
osmosisd add-genesis-account "$(osmosisd keys show osmosis-user -a --keyring-backend=test)" 200000000000uosmo
osmosisd gentx osmosis-user 50000000000uosmo --keyring-backend=test --chain-id "$CHAIN_ID"
osmosisd collect-gentxs
