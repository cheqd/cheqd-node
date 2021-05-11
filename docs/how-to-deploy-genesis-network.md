# Deploy genesis network

This document describes in details how to configure a genesis network with any amount of particiatns.

### Hardware requirements
Minimal:
- 1GB RAM
- 25GB of disk space
- 1.4 GHz CPU

Recommended (for highload applications):
- 2GB RAM
- 100GB SSD
- x64 2.0 GHz 2v CPU

### Operating System
Current delivery is compiled and tested under `Ubuntu 20.04 LTS` so we recommend using this distribution for now. In the future, it will be possible to compile the application for a wide range of operating systems thanks to Go language.

## Deployment steps

1 Jack

  1. Generates local key: `verim-cosmosd keys add jack --home $NODE_0_HOME`, where `$NODE_0_HOME="localnet/node0"`

  2. Initializes local node: `verim-cosmosd init node0 --chain-id $CHAIN_ID --home $NODE_0_HOME`, where `$CHAIN_ID="verim-cosmos-chain"`

1.3 Jack adds genesis account: `verim-cosmosd add-genesis-account jack 10000000token,100000000stake --home $NODE_0_HOME`

1.4 Jack generates genesis node transaction: `verim-cosmosd gentx jack 1000000stake --chain-id $CHAIN_ID --home $NODE_0_HOME`

<br/>

2.1 Alice generates local key: `verim-cosmosd keys add alice --home $NODE_1_HOME`, where `$NODE_1_HOME="localnet/node1"`

2.2 Alice initializes local node: `verim-cosmosd init node1 --chain-id $CHAIN_ID --home $NODE_1_HOME`

2.3 Alice gets genesis from Jack. If it is the only machine: `cp $NODE_0_HOME/config/genesis.json $NODE_1_HOME/config`

2.4 Alice gets genesis node transactions form Jack. If it is the only machine:
- `mkdir $NODE_1_HOME/config/gentx`
- `cp $NODE_0_HOME/config/gentx/* $NODE_1_HOME/config/gentx`

2.5 Alice adds genesis account: `verim-cosmosd add-genesis-account alice 10000000token,100000000stake --home $NODE_1_HOME`

<br/>

3.1 Bob generates local key: `verim-cosmosd keys add bob --home $NODE_2_HOME`, where `$NODE_2_HOME="localnet/node2"`

3.2 Bob initializes local node: `verim-cosmosd init node2 --chain-id $CHAIN_ID --home $NODE_2_HOME`

3.3 Bob gets genesis from Alice. If it is the only machine: `cp $NODE_1_HOME/config/genesis.json $NODE_2_HOME/config`

3.4 Bob gets genesis node transactions form Jack. If it is the only machine:
- `mkdir $NODE_2_HOME/config/gentx`
- `cp $NODE_1_HOME/config/gentx/* $NODE_2_HOME/config/gentx`

3.5 Bob adds genesis account: `verim-cosmosd add-genesis-account bob 10000000token,100000000stake --home $NODE_2_HOME`

<br/>

4.1 Anna generates local key: `verim-cosmosd keys add bob --home $NODE_3_HOME`, where `$NODE_3_HOME="localnet/node3"`

4.2 Anna initializes local node: `verim-cosmosd init node3 --chain-id $CHAIN_ID --home $NODE_2_HOME`

4.3 Anna gets genesis from Bob. If it is the only machine: `cp $NODE_2_HOME/config/genesis.json $NODE_3_HOME/config`

4.4 Anna gets genesis node transactions form Jack. If it is the only machine:
- `mkdir $NODE_3_HOME/config/gentx`
- `cp $NODE_2_HOME/config/gentx/* $NODE_3_HOME/config/gentx`

4.5 Anna adds genesis account: `verim-cosmosd add-genesis-account anna 10000000token,100000000stake --home $NODE_3_HOME`

4.6 Eventually Anna adds genesis node transactions into genesis: `verim-cosmosd collect-gentxs --home $NODE_3_HOME`

4.7 Anna verifies genesis: `verim-cosmosd validate-genesis --home $NODE_3_HOME`

4.8 Anna shares her genesis with other nodes. If it is the only machine:
- `cp $NODE_3_HOME/config/genesis.json $NODE_0_HOME/config/`
- `cp $NODE_3_HOME/config/genesis.json $NODE_1_HOME/config/`
- `cp $NODE_3_HOME/config/genesis.json $NODE_2_HOME/config/`

<br/>

5.1 Updates address book of them node. It allows nodes to connect to each other. Every node owner needs to run follow command: `sed -i "s/persistent_peers = \"\"/persistent_peers = \"$NODE_0_ID@node0:26656,$NODE_1_ID@node1:26656,$NODE_2_ID@node2:26656,$NODE_3_ID@node3:26656\"/g" $NODE_#_HOME/config/config.toml`, where $NODE_#_HOME means variable with number of node owner.

5.2 Eventually need to set minimal gas prices (it needs for fees transaction). Every node owner needs to run follow command: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_#_HOME/config/app.toml`, where instead `1token` node can set other price.

<br/>

After this steps the nodes of Jack, Alice, Bob and Anna have same genesis, and they can connect to each other.
