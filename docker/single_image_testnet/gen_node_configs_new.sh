#!/bin/bash

# Generates configurations for 4 nodes.

# We can't pack several nodes without having `--home` flag working
# so reducing nodes count to 1 for now and creating follow up ticket.

set -euox pipefail

# sed in macos requires extra argument

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sed_extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    sed_extension='.orig'
fi

function info() {
    printf "[info] %s\n" "${1}"
}


CHAIN_ID="cheqd"
NODES_COUNT="4"

for ((i=0 ; i<$NODES_COUNT ; i++))
do
    cheqd-noded init "node$i" --chain-id $CHAIN_ID --home "node$i"
    cheqd-noded keys add "operator$i" --keyring-backend test
done


# NODE_0_ID=$(cheqd-noded tendermint show-node-id)
# NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator)

# Operator 0
GENESIS="node0/config/genesis.json"
sed -i $sed_extension 's/"stake"/"ncheq"/' $GENESIS
cheqd-noded add-genesis-account "operator0" 20000000000000000ncheq --home "node0" --keyring-backend test
cheqd-noded gentx alice 1000000000000000ncheq --chain-id $CHAIN_ID --node-id "$NODE_0_ID" --pubkey "$NODE_0_VAL_PUBKEY" --home "node0" --keyring-backend test

BASE_ACCOUNT_1="cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96"
# Mnemonic: sketch mountain erode window enact net enrich smoke claim kangaroo another visual write meat latin bacon pulp similar forum guilt father state erase bright
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_ACCOUNT_1}'", "coins": [{"denom": "ncheq", "amount": "100001000000000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${BASE_ACCOUNT_1}'", "pub_key": null,"account_number": "0","sequence": "0"}]' "$GENESIS")" > "$GENESIS"

BASE_ACCOUNT_2="cheqd1l9sq0se0jd3vklyrrtjchx4ua47awug5vsyeeh"
# Mnemonic: ugly dirt sorry girl prepare argue door man that manual glow scout bomb pigeon matter library transfer flower clown cat miss pluck drama dizzy
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_ACCOUNT_2}'", "coins": [{"denom": "ncheq", "amount": "100001000000000000"}] }]' "$GENESIS")" > $GENESIS
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${BASE_ACCOUNT_2}'", "pub_key": null,"account_number": "0","sequence": "0"}]' "$GENESIS")" > "$GENESIS"

BASE_VESTING_ACCOUNT="cheqd1lkqddnapqvz2hujx2trpj7xj6c9hmuq7uhl0md"
# Mnemonic: coach index fence broken very cricket someone casino dial truth fitness stay habit such three jump exotic spawn planet fragile walk enact angry great
BASE_VESTING_COIN="{\"denom\":\"ncheq\",\"amount\":\"10001000000000000\"}"
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${BASE_VESTING_ACCOUNT}'", "coins": [{"denom": "ncheq", "amount": "5000000000000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.BaseVestingAccount", "base_account": {"address": "'${BASE_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}]' "$GENESIS")" > "$GENESIS"

CONTINOUS_VESTING_ACCOUNT="cheqd1353p46macvn444rupg2jstmx3tmz657yt9gl4l"
# Mnemonic: phone worry flame safe panther dirt picture pepper purchase tiny search theme issue genre orange merit stove spoil surface color garment mind chuckle image
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${CONTINOUS_VESTING_ACCOUNT}'", "coins": [{"denom": "ncheq", "amount": "5000000000000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount", "base_vesting_account": { "base_account": {"address": "'${CONTINOUS_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630352459"}]' "$GENESIS")" > "$GENESIS"

DELAYED_VESTING_ACCOUNT="cheqd1njwu33lek5jt4kzlmljkp366ny4qpqusahpyrj"
# Mnemonic: pilot text keen deal economy donkey use artist divide foster walk pink breeze proud dish brown icon shaft infant level labor lift will tomorrow
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${DELAYED_VESTING_ACCOUNT}'", "coins": [{"denom": "ncheq", "amount": "5000000000000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.DelayedVestingAccount", "base_vesting_account": { "base_account": {"address": "'${DELAYED_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}}]' "$GENESIS")" > "$GENESIS"

PERIODIC_VESTING_ACCOUNT="cheqd1uyngr0l3xtyj07js9sdew9mk50tqeq8lghhcfr"
# Mnemonic: want merge flame plate trouble moral submit wing whale sick meat lonely yellow lens enable oyster slight health vast weird radar mesh grab olive
cat <<< "$(jq '.app_state.bank.balances += [{"address": "'${PERIODIC_VESTING_ACCOUNT}'", "coins": [{"denom": "ncheq", "amount": "5000000000000000"}] }]' "$GENESIS")" > "$GENESIS"
cat <<< "$(jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.PeriodicVestingAccount", "base_vesting_account": { "base_account": {"address": "'${PERIODIC_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630362439", "vesting_periods": [{"length": "20", "amount": ['${BASE_VESTING_COIN}']}]}]' "$GENESIS")" > "$GENESIS"


# Operators 1..N

for ((i=1 ; i<$VALIDATORS_COUNT ; i++))
do
    cp "node$((i-1))/config/genesis.json" "node$i/config/genesis.json"
    cp "node$((i-1))/config/genesis.json" "node$i/config/genesis.json"

    cheqd-noded add-genesis-account "operator$i" 20000000000000000ncheq --home "node$i" --keyring-backend test
    cheqd-noded gentx alice 1000000000000000ncheq --chain-id $CHAIN_ID --node-id "$NODE_0_ID" --pubkey "$NODE_0_VAL_PUBKEY" --home "node$i" --keyring-backend test
done

echo "##### [Validator operator] Collect gentxs"

cheqd-noded collect-gentxs
cheqd-noded validate-genesis
