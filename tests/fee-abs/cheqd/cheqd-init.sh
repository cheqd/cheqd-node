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
echo 'success claw flock lecture soul photo jump pause sadness enter uncle stage cherry teach inside now fun fuel oval angry unveil horn tobacco penalty' | cheqd-noded keys add cheqd-user --keyring-backend test --recover

CONFIG_TOML="$HOME/.cheqdnode/config/config.toml"

# Config
sed -i $SED_EXT 's|minimum-gas-prices = ""|minimum-gas-prices = "50ncheq"|g' "$HOME/.cheqdnode/config/app.toml"
sed -i $SED_EXT 's|addr_book_strict = true|addr_book_strict = false|g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_propose = "3s"/timeout_propose = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/timeout_commit = "5s"/timeout_commit = "2s"/g' "${CONFIG_TOML}"
sed -i $SED_EXT 's/log_level = "info"/log_level = "debug"/g' "${CONFIG_TOML}"


sed -i $SED_EXT 's/"voting_period":"172800s"/"voting_period":"12s"/' "$HOME/.cheqdnode/config/genesis.json"
sed -i $SED_EXT 's/"vote_extensions_enable_height"[[:space:]]*:[[:space:]]*"0"/"vote_extensions_enable_height": "2"/' "$HOME/.cheqdnode/config/genesis.json"

sed -i $SED_EXT 's/"expedited_voting_period": "86400s"/"expedited_voting_period": "10s"/' "$HOME/.cheqdnode/config/genesis.json"

# shellcheck disable=SC2086
sed -i $SED_EXT 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME/.cheqdnode/config/config.toml"
sed -i $SED_EXT 's|address = "localhost:9090"|address = "0.0.0.0:9090"|g' "$HOME/.cheqdnode/config/app.toml"
sed -i $SED_EXT 's|log_level = "error"|log_level = "info"|g' "$HOME/.cheqdnode/config/config.toml"

# Genesis
GENESIS="$HOME/.cheqdnode/config/genesis.json"
sed -i $SED_EXT 's/"stake"/"ncheq"/' "$GENESIS"

cheqd-noded genesis add-genesis-account cheqd-user 1000000000000000000000ncheq --keyring-backend test
cheqd-noded genesis gentx cheqd-user 1000000000000000000000ncheq --chain-id $CHAIN_ID --pubkey "$NODE_0_VAL_PUBKEY" --keyring-backend test

cheqd-noded genesis collect-gentxs
cheqd-noded genesis validate-genesis

# set default output format to json
cheqd-noded config set client output json
