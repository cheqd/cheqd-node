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

NODE_0_HOME="$NODE_CONFIGS_DIR/node0"

verim-noded init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME
cp -r $NODE_CONFIGS_DIR/client/* $NODE_0_HOME

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_0_HOME
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_0_HOME
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_0_HOME
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_0_HOME

verim-noded gentx jack 1000000stake --chain-id $CHAIN_ID --home $NODE_0_HOME

export NODE_0_ID=$(verim-noded tendermint show-node-id --home $NODE_0_HOME)

# node 1

NODE_1_HOME="$NODE_CONFIGS_DIR/node1"

verim-noded init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME
cp -r $NODE_CONFIGS_DIR/client/* $NODE_1_HOME

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_1_HOME
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_1_HOME
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_1_HOME
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_1_HOME

verim-noded gentx alice 1000000stake --chain-id $CHAIN_ID --home $NODE_1_HOME

export NODE_1_ID=$(verim-noded tendermint show-node-id --home $NODE_1_HOME)

# node 2

NODE_2_HOME="$NODE_CONFIGS_DIR/node2"

verim-noded init node2 --chain-id $CHAIN_ID --home $NODE_2_HOME
cp -r $NODE_CONFIGS_DIR/client/* $NODE_2_HOME

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_2_HOME
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_2_HOME
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_2_HOME
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_2_HOME

verim-noded gentx bob 1000000stake --chain-id $CHAIN_ID --home $NODE_2_HOME

export NODE_2_ID=$(verim-noded tendermint show-node-id --home $NODE_2_HOME)

# node 3

NODE_3_HOME="$NODE_CONFIGS_DIR/node3"

verim-noded init node3 --chain-id $CHAIN_ID --home $NODE_3_HOME
cp -r $NODE_CONFIGS_DIR/client/* $NODE_3_HOME

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_3_HOME
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_3_HOME
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_3_HOME
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_3_HOME

verim-noded gentx anna 1000000stake --chain-id $CHAIN_ID --home $NODE_3_HOME

export NODE_3_ID=$(verim-noded tendermint show-node-id --home $NODE_3_HOME)

# Collect all validator creation transactions

mkdir $NODE_CONFIGS_DIR/client/config/gentx

cp $NODE_0_HOME/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_1_HOME/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_2_HOME/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx
cp $NODE_3_HOME/config/gentx/* $NODE_CONFIGS_DIR/client/config/gentx

# Embed them into genesis

verim-noded init dummy-node --chain-id $CHAIN_ID --home $NODE_CONFIGS_DIR/client

verim-noded add-genesis-account jack 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account alice 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account bob 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client
verim-noded add-genesis-account anna 10000000token,100000000stake --home $NODE_CONFIGS_DIR/client

verim-noded collect-gentxs --home $NODE_CONFIGS_DIR/client
verim-noded validate-genesis --home $NODE_CONFIGS_DIR/client

# Update genesis for all nodes

cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_0_HOME/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_1_HOME/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_2_HOME/config/
cp $NODE_CONFIGS_DIR/client/config/genesis.json $NODE_3_HOME/config/


# # sed in macos requires extra argument
# extension=''
# if [[ "$OSTYPE" == "linux-gnu"* ]]; then
#     extension=''
# elif [[ "$OSTYPE" == "darwin"* ]]; then
#     extension='.orig'
# fi


# # Make RPC enpoint available externally
# sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $NODE_0_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26660"/g' $NODE_1_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26663"/g' $NODE_2_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26666"/g' $NODE_3_HOME/config/config.toml

# # Make p2p enpoints different
# sed -i $extension 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26656"/g' $NODE_0_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26659"/g' $NODE_1_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26662"/g' $NODE_2_HOME/config/config.toml
# sed -i $extension 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26665"/g' $NODE_3_HOME/config/config.toml

# # Update address book of the first node
# peers="$id0@node0:26656,$id1@node1:26659,$id2@node2:26662,$id3@node3:26665"
# sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" $NODE_0_HOME/config/config.toml

# # Set gas prices
# sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_0_HOME/config/app.toml
# sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_1_HOME/config/app.toml
# sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_2_HOME/config/app.toml
# sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_3_HOME/config/app.toml
