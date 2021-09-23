#!/bin/bash

# Generates configurations for 2 nodes.

set -euox pipefail

CHAIN_ID="cheqd"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR


echo "##### [Node 0] Generates key" 

NODE_0_HOME="$NODE_CONFIGS_DIR/node0"
cheqd-noded init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME
NODE_0_ID=$(cheqd-noded tendermint show-node-id --home $NODE_0_HOME)
NODE_0_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_0_HOME)


echo "##### [Node 1] Generates key" 

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"
cheqd-noded init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME
NODE_1_ID=$(cheqd-noded tendermint show-node-id --home $NODE_1_HOME)
NODE_1_VAL_PUBKEY=$(cheqd-noded tendermint show-validator --home $NODE_1_HOME)


echo "##### [Validator operators] Generate keys" 

CLIENT_HOME=$NODE_CONFIGS_DIR/client
cheqd-noded keys add alice --home $CLIENT_HOME
cheqd-noded keys add bob --home $CLIENT_HOME


echo "##### [Validator operators] Init genesis" 

cheqd-noded init dummy_node --chain-id $CHAIN_ID --home $CLIENT_HOME

sed -i 's/"stake"/"cheq"/' ${CLIENT_HOME}/config/genesis.json

echo "##### [Validator operators] Add them to the genesis" 

cheqd-noded add-genesis-account alice 20000000cheq --home $CLIENT_HOME
cheqd-noded add-genesis-account bob 20000000cheq --home $CLIENT_HOME



echo "##### [Test pool] Add test accounts to the genesis" 

ALICE_OLD_ACCOUNT_ID="cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96"
# Mnemonic: sketch mountain erode window enact net enrich smoke claim kangaroo another visual write meat latin bacon pulp similar forum guilt father state erase bright
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${ALICE_OLD_ACCOUNT_ID}'", "coins": [{"denom": "cheq", "amount": "100001000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${ALICE_OLD_ACCOUNT_ID}'", "pub_key": null,"account_number": "0","sequence": "0"}]') > ${CLIENT_HOME}/config/genesis.json

ALICE_NEW_ACCOUNT_ID="cheqd1l9sq0se0jd3vklyrrtjchx4ua47awug5vsyeeh"
# Mnemonic: ugly dirt sorry girl prepare argue door man that manual glow scout bomb pigeon matter library transfer flower clown cat miss pluck drama dizzy
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${ALICE_NEW_ACCOUNT_ID}'", "coins": [{"denom": "cheq", "amount": "100001000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.auth.v1beta1.BaseAccount","address": "'${ALICE_NEW_ACCOUNT_ID}'", "pub_key": null,"account_number": "0","sequence": "0"}]') > ${CLIENT_HOME}/config/genesis.json


BASE_VESTING_ACCOUNT="cheqd1lkqddnapqvz2hujx2trpj7xj6c9hmuq7uhl0md"
# Mnemonic: coach index fence broken very cricket someone casino dial truth fitness stay habit such three jump exotic spawn planet fragile walk enact angry great
BASE_VESTING_COIN='{"denom":"cheq","amount":"10001000"}'
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${BASE_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.BaseVestingAccount", "base_account": {"address": "'${BASE_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}]') > ${CLIENT_HOME}/config/genesis.json

CONTINOUS_VESTING_ACCOUNT="cheqd1353p46macvn444rupg2jstmx3tmz657yt9gl4l"
# Mnemonic: phone worry flame safe panther dirt picture pepper purchase tiny search theme issue genre orange merit stove spoil surface color garment mind chuckle image
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${CONTINOUS_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount", "base_vesting_account": { "base_account": {"address": "'${CONTINOUS_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630352459"}]') > ${CLIENT_HOME}/config/genesis.json

DELAYED_VESTING_ACCOUNT="cheqd1njwu33lek5jt4kzlmljkp366ny4qpqusahpyrj"
# Mnemonic: pilot text keen deal economy donkey use artist divide foster walk pink breeze proud dish brown icon shaft infant level labor lift will tomorrow
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${DELAYED_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.DelayedVestingAccount", "base_vesting_account": { "base_account": {"address": "'${DELAYED_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}}]') > ${CLIENT_HOME}/config/genesis.json

PERIODIC_VESTING_ACCOUNT="cheqd1uyngr0l3xtyj07js9sdew9mk50tqeq8lghhcfr"
# Mnemonic: want merge flame plate trouble moral submit wing whale sick meat lonely yellow lens enable oyster slight health vast weird radar mesh grab olive
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.bank.balances += [{"address": "'${PERIODIC_VESTING_ACCOUNT}'", "coins": [{"denom": "cheq", "amount": "5000000"}] }]') > ${CLIENT_HOME}/config/genesis.json
echo $(cat ${CLIENT_HOME}/config/genesis.json | jq '.app_state.auth.accounts += [{"@type": "/cosmos.vesting.v1beta1.PeriodicVestingAccount", "base_vesting_account": { "base_account": {"address": "'${PERIODIC_VESTING_ACCOUNT}'","pub_key": null,"account_number": "0","sequence": "0"}, "original_vesting": ['${BASE_VESTING_COIN}'], "delegated_free": [], "delegated_vesting": [], "end_time": "1630362459"}, "start_time": "1630362439", "vesting_periods": [{"length": "20", "amount": ['${BASE_VESTING_COIN}']}]}]') > ${CLIENT_HOME}/config/genesis.json

echo "##### [Validator operators] Generate stake transactions"

cheqd-noded gentx alice 1000000cheq --chain-id $CHAIN_ID --node-id $NODE_0_ID --pubkey $NODE_0_VAL_PUBKEY --home $CLIENT_HOME
cheqd-noded gentx bob 1000000cheq --chain-id $CHAIN_ID --node-id $NODE_1_ID --pubkey $NODE_1_VAL_PUBKEY --home $CLIENT_HOME


echo "##### [Validator operators] Collect them"

cheqd-noded collect-gentxs --home $CLIENT_HOME
cheqd-noded validate-genesis --home $CLIENT_HOME


echo "##### [Validator operators] Propagate genesis to nodes"

cp $CLIENT_HOME/config/genesis.json $NODE_0_HOME/config/
cp $CLIENT_HOME/config/genesis.json $NODE_1_HOME/config/
