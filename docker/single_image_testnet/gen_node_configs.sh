#!/bin/bash

# Generates configurations for 1 node.

# We can't pack several nodes without having `--home` flag working
# so reducing nodes count to 1 for now and creating follow up ticket.

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

CHAIN_ID="cheqd"

echo "##### [Node 0] Generating key"

cheqd-noded init node0 --chain-id $CHAIN_ID
NODE_0_ID=$(cheqd-noded tendermint show-node-id)
NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)

echo "##### [Validator operator] Generating key"

cheqd-noded keys add alice

echo "##### [Validator operator] Initializing genesis"

GENESIS="$HOME/.cheqdnode/config/genesis.json"
sed -i $sed_extension 's/"stake"/"cheq"/' $GENESIS

echo "##### [Validator operator] Creating genesis account"

cheqd-noded add-genesis-account alice 20000000cheq

echo "##### Adding test accounts to the genesis"

BASE_ACCOUNT_1="cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_ACCOUNT_1}'", "coins": [{"denom": "cheq", "amount": "100001000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${BASE_ACCOUNT_1}'", "pub_key": null,"account_number": "0","sequence": "0"}]' "$GENESIS")" > "$GENESIS"

BASE_ACCOUNT_2="cosmos1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_ACCOUNT_2}'", "coins": [{"denom": "cheq", "amount": "100001000"}] }]' "$GENESIS")" > $GENESIS
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${BASE_ACCOUNT_2}'", "pub_key": null,"account_number": "0","sequence": "0"}]' "$GENESIS")" > "$GENESIS"

BASE_VESTING_ACCOUNT="cosmos1h0ul2knsd6pa4spfxtvznxfy6qma34uhxtu8zd"
BASE_VESTING_COIN="{\"denom\":\"cheq\",\"amount\":\"10001000\"}"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.BaseVestingAccount", "base_account": {"address": "'${BASE_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}]' "$GENESIS")" > "$GENESIS"

CONTINOUS_VESTING_ACCOUNT="cosmos16gn9jhq4cztt9rkg8pt6r8zeruy6swzlwurcay"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${CONTINOUS_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount", "base_vesting_account": { "base_account": {"address": "'${CONTINOUS_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630352459"}]' "$GENESIS")" > "$GENESIS"

DELAYED_VESTING_ACCOUNT="cosmos1830c9zt72lwaqgk7yxjcxpgawwqhq9mlla8yhx"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${DELAYED_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.DelayedVestingAccount", "base_vesting_account": { "base_account": {"address": "'${DELAYED_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}}]' "$GENESIS")" > "$GENESIS"

PERIODIC_VESTING_ACCOUNT="cosmos1ecnhll5kery5td3ensfefeszfd208rv5tsa4u6"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${PERIODIC_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.PeriodicVestingAccount", "base_vesting_account": { "base_account": {"address": "'${PERIODIC_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630362439", "vesting_periods": [{"length": "20", "amount": ['${BASE_VESTING_COIN}']}]}]' "$GENESIS")" > "$GENESIS"

echo "##### [Validator operator] Creating genesis validator"

cheqd-noded gentx alice 1000000cheq --chain-id $CHAIN_ID --node-id "$NODE_0_ID" --pubkey "$NODE_0_VAL_PUBKEY"

echo "##### [Validator operator] Collect gentxs"

cheqd-noded collect-gentxs
cheqd-noded validate-genesis
