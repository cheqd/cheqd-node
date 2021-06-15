#!/bin/bash

# Generates configurations for 4 nodes. Mostly the same script as `gen_node_configs.sh` but discribes
# the real life flow (keys don't leave the place where they are generated).

set -euox pipefail

CHAIN_ID="verim"
NODE_CONFIGS_DIR="node_configs"

rm -rf $NODE_CONFIGS_DIR
mkdir $NODE_CONFIGS_DIR

# client

verim-noded keys add jack --home $NODE_CONFIGS_DIR/client
verim-noded keys add alice --home $NODE_CONFIGS_DIR/client
verim-noded keys add bob --home $NODE_CONFIGS_DIR/client
verim-noded keys add anna --home $NODE_CONFIGS_DIR/client

# node 0

verim-noded init node0 --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node0
cp -r $NODE_CONFIGS_DIR/client/* $NODE_CONFIGS_DIR/node0

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node0
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node0
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node0
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node0

verim-noded gentx jack 1000000stake --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node0

# node 1

verim-noded init node1 --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node1
cp -r $NODE_CONFIGS_DIR/client/* $NODE_CONFIGS_DIR/node1

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node1
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node1
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node1
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node1

verim-noded gentx alice 1000000stake --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node1

# node 2

verim-noded init node2 --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node2
cp -r $NODE_CONFIGS_DIR/client/* $NODE_CONFIGS_DIR/node2

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node2
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node2
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node2
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node2

verim-noded gentx bob 1000000stake --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node2

# node 3

verim-noded init node3 --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node3
cp -r $NODE_CONFIGS_DIR/client/* $NODE_CONFIGS_DIR/node3

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node3
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node3
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node3
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/node3

verim-noded gentx anna 1000000stake --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/node3

# Collect all validator creation transactions

mkdir $NODE_CONFIGS_DIR/client/config/gentx

cp $NODE_CONFIGS_DIR/node0/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_CONFIGS_DIR/node1/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_CONFIGS_DIR/node2/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_CONFIGS_DIR/node3/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx

# Embed them into genesis

verim-noded init dummy-node --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/client

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client

verim-noded collect-gentxs --home $NODE_CONFIGS_DIR/client
verim-noded validate-genesis --home $NODE_CONFIGS_DIR/client

# Update genesis for all nodes

cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_CONFIGS_DIR/node0/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_CONFIGS_DIR/node1/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_CONFIGS_DIR/node2/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_CONFIGS_DIR/node3/config/

# Find out node ids

id0=$(ls $NODE_CONFIGS_DIR/node0/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id1=$(ls $NODE_CONFIGS_DIR/node1/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id2=$(ls $NODE_CONFIGS_DIR/node2/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id3=$(ls $NODE_CONFIGS_DIR/node3/config/gentx | sed 's/gentx-\(.*\).json/\1/')

# sed in macos requires extra argument
extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    extension='.orig'
fi

# Update address book of the first node
peers="$id0@node0:26656,$id1@node1:26656,$id2@node2:26656,$id3@node3:26656"
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_CONFIGS_DIR/node0/config/config.toml

# Make RPC enpoint available externally
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_CONFIGS_DIR/node0/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_CONFIGS_DIR/node1/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_CONFIGS_DIR/node2/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_CONFIGS_DIR/node3/config/config.toml


# Set gas prices
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_CONFIGS_DIR/node0/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_CONFIGS_DIR/node1/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_CONFIGS_DIR/node2/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_CONFIGS_DIR/node3/config/app.toml
