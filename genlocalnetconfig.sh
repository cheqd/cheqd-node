#!/bin/bash

set -euox pipefail

CHAIN_ID="verimcosmos"

rm -rf localnet
mkdir localnet

# client

verim-cosmosd keys add jack --home localnet/client
verim-cosmosd keys add alice --home localnet/client
verim-cosmosd keys add bob --home localnet/client
verim-cosmosd keys add anna --home localnet/client

# node 0

verim-cosmosd init node0 --chain-id $CHAIN_ID --home localnet/node0
cp -r localnet/client/* localnet/node0

verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home localnet/node0
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home localnet/node0
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home localnet/node0
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home localnet/node0

verim-cosmosd gentx jack 1000000stake --chain-id $CHAIN_ID --home localnet/node0

# node 1

verim-cosmosd init node1 --chain-id $CHAIN_ID --home localnet/node1
cp -r localnet/client/* localnet/node1

verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home localnet/node1
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home localnet/node1
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home localnet/node1
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home localnet/node1

verim-cosmosd gentx alice 1000000stake --chain-id $CHAIN_ID --home localnet/node1

# node 2

verim-cosmosd init node2 --chain-id $CHAIN_ID --home localnet/node2
cp -r localnet/client/* localnet/node2

verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home localnet/node2
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home localnet/node2
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home localnet/node2
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home localnet/node2

verim-cosmosd gentx bob 1000000stake --chain-id $CHAIN_ID --home localnet/node2

# node 3

verim-cosmosd init node3 --chain-id $CHAIN_ID --home localnet/node3
cp -r localnet/client/* localnet/node3

verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home localnet/node3
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home localnet/node3
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home localnet/node3
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home localnet/node3

verim-cosmosd gentx anna 1000000stake --chain-id $CHAIN_ID --home localnet/node3

# Collect all validator creation transactions

mkdir localnet/client/config/gentx

cp localnet/node0/config/gentx/* localnet/client/config/gentx
cp localnet/node1/config/gentx/* localnet/client/config/gentx
cp localnet/node2/config/gentx/* localnet/client/config/gentx
cp localnet/node3/config/gentx/* localnet/client/config/gentx

# Embed them into genesis

verim-cosmosd init dummy-node --chain-id $CHAIN_ID --home localnet/client

verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home localnet/client
verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home localnet/client
verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home localnet/client
verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home localnet/client

verim-cosmosd collect-gentxs --home localnet/client
verim-cosmosd validate-genesis --home localnet/client

# Update genesis for all nodes

cp localnet/client/config/genesis.json localnet/node0/config/
cp localnet/client/config/genesis.json localnet/node1/config/
cp localnet/client/config/genesis.json localnet/node2/config/
cp localnet/client/config/genesis.json localnet/node3/config/

# Find out node ids

id0=$(ls localnet/node0/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id1=$(ls localnet/node1/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id2=$(ls localnet/node2/config/gentx | sed 's/gentx-\(.*\).json/\1/')
id3=$(ls localnet/node3/config/gentx | sed 's/gentx-\(.*\).json/\1/')

# sed in macos requires extra argument
extension=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    extension=''
elif [[ "$OSTYPE" == "darwin"* ]]; then
    extension='.orig'
fi

# Update address book of the first node
peers="$id0@node0:26656,$id1@node1:26656,$id2@node2:26656,$id3@node3:26656"
sed -i $extension "s/persistent_peers = \"\"/persistent_peers = \"$peers\"/g" localnet/node0/config/config.toml

# Make RPC enpoint available externally
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node0/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node1/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node2/config/config.toml
sed -i $extension 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' localnet/node3/config/config.toml


# Set gas prices
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node0/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node1/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node2/config/app.toml
sed -i $extension 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' localnet/node3/config/app.toml
