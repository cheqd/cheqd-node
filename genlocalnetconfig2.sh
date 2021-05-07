#!/bin/bash

# Creates network of 4 nodes. Mostly the same script as `genlocalnetconfig.sh` but discribes the production flow (keys don't leave the place where they where generated).

set -euox pipefail

CHAIN_ID="verim-cosmos-chain"

rm -rf localnet
mkdir localnet

# sed in macos requires extra argument
extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    extension='.orig'
fi


echo "################################ Jack's node0"

NODE_0_HOME="localnet/node0"

echo "# Generate key"
verim-cosmosd keys add jack --home $NODE_0_HOME

echo "# Initialze node"
verim-cosmosd init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME

echo "# Add genesis account"
verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home $NODE_0_HOME

echo "# Generate genesis node tx"
verim-cosmosd gentx jack 1000000stake --chain-id $CHAIN_ID --home $NODE_0_HOME

echo "# Publish validator pubkey"
NODE_0_ID=$(verim-cosmosd tendermint show-node-id --home $NODE_0_HOME)

echo "# Make RPC enpoint available externally (optional, allows cliens to connect to the node)"
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_0_HOME/config/config.toml


echo "################################ Alice's node1"

NODE_1_HOME="localnet/node1"

echo "# Generate key"
verim-cosmosd keys add alice --home $NODE_1_HOME

echo "# Initialze node"
verim-cosmosd init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME

echo "### Get genesis from Jack"
cp $NODE_0_HOME/config/genesis.json $NODE_1_HOME/config

echo "### Get genesis node txs form Jack"
mkdir $NODE_1_HOME/config/gentx
cp $NODE_0_HOME/config/gentx/* $NODE_1_HOME/config/gentx

echo "# Add genesis account"
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home $NODE_1_HOME

echo "# Generate genesis node tx"
verim-cosmosd gentx alice 1000000stake --chain-id $CHAIN_ID --home $NODE_1_HOME

echo "# Publish validator pubkey"
NODE_1_ID=$(verim-cosmosd tendermint show-node-id --home $NODE_1_HOME)

echo "# Make RPC enpoint available externally (optional, allows cliens to connect to the node)"
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_1_HOME/config/config.toml


echo "################################ Bob's node2"

NODE_2_HOME="localnet/node2"

echo "# Generate key"
verim-cosmosd keys add bob --home $NODE_2_HOME

echo "# Initialze node"
verim-cosmosd init node1 --chain-id $CHAIN_ID --home $NODE_2_HOME

echo "### Get genesis from Alice"
cp $NODE_1_HOME/config/genesis.json $NODE_2_HOME/config

echo "### Get genesis node txs form Alice"
mkdir $NODE_2_HOME/config/gentx
cp $NODE_1_HOME/config/gentx/* $NODE_2_HOME/config/gentx

echo "# Add genesis account"
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home $NODE_2_HOME

echo "# Generate genesis node tx"
verim-cosmosd gentx bob 1000000stake --chain-id $CHAIN_ID --home $NODE_2_HOME

echo "# Add the tx into genesis"
verim-cosmosd collect-gentxs --home $NODE_2_HOME

echo "# Publish validator pubkey"
NODE_2_ID=$(verim-cosmosd tendermint show-node-id --home $NODE_2_HOME)

echo "# Make RPC enpoint available externally (optional, allows cliens to connect to the node)"
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_2_HOME/config/config.toml


echo "################################ Anna's node3"

NODE_3_HOME="localnet/node3"

echo "# Generate key"
verim-cosmosd keys add anna --home $NODE_3_HOME

echo "# Initialze node"
verim-cosmosd init node3 --chain-id $CHAIN_ID --home $NODE_3_HOME

echo "### Get genesis from Bob"
cp $NODE_2_HOME/config/genesis.json $NODE_3_HOME/config

echo "### Get genesis node txs form Bob"
mkdir $NODE_3_HOME/config/gentx
cp $NODE_2_HOME/config/gentx/* $NODE_3_HOME/config/gentx

echo "# Add genesis account"
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home $NODE_3_HOME

echo "# Generate genesis node tx"
verim-cosmosd gentx anna 1000000stake --chain-id $CHAIN_ID --home $NODE_3_HOME

echo "# Publish validator pubkey"
NODE_3_ID=$(verim-cosmosd tendermint show-node-id --home $NODE_3_HOME)

echo "# Make RPC enpoint available externally (optional, allows cliens to connect to the node)"
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_3_HOME/config/config.toml


echo "################################ Anna (last participatn) shares genesis with everyone else"

echo "# Add genesis node txs into genesis"
verim-cosmosd collect-gentxs --home $NODE_3_HOME

echo "# Verify genesis"
verim-cosmosd validate-genesis --home $NODE_3_HOME

cp $NODE_3_HOME/config/genesis.json $NODE_0_HOME/config/
cp $NODE_3_HOME/config/genesis.json $NODE_1_HOME/config/
cp $NODE_3_HOME/config/genesis.json $NODE_2_HOME/config/


echo "################################ Anna (at least one participant) Updates address book of her node. It will alow nodes to connect to each other."
peers="$NODE_0_ID@node0:26656,$NODE_1_ID@node1:26656,$NODE_2_ID@node2:26656,$NODE_3_ID@node3:26656"
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_0_HOME/config/config.toml
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_1_HOME/config/config.toml
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_2_HOME/config/config.toml
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_3_HOME/config/config.toml


echo "################################ (any participant, optional) Sets minimal gas prices"
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node0/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node1/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node2/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node3/config/app.toml
