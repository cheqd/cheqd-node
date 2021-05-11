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

### Generating genesis file

0. Participants must agree on `<chain_id>`
1. Each participant (one by one):
    - Generates local keys for his future account:  
        - Command: `verim-cosmosd keys add <key_name>`
        - Example: `verim-cosmosd keys add alice`
    - Initializes local node:
        - Command: `verim-cosmosd init <node_name> --chain-id <chain_id>`
        - Example: `verim-cosmosd init alice-node --chain-id verim-cosmos`
    - (Each participatn except the first one) Gets genesis from the previous participant:
        - Location on the previous participant's machine: `$HOME/.verimcosmos/config/genesis.json`
        - Destination folder on the current participant's machine: `$HOME/.verimcosmos/config/`
    - (Each participatn except the first one) Gets genesis node transactions form the previous participant.
        - Location on the previous participant's machine: `$HOME/.verimcosmos/config/gentx/`
        - Destination folder on the current participant's machine: `$HOME/.verimcosmos/config/gentx/`
        - TODO: We can specify `minimum-gas-prices` and some other price releted flags. Need to find out how it works.
    - Adds a genesis account with his public key:
        - Command: `verim-cosmosd add-genesis-account <key_name> 10000000token,100000000stake`
        - Example: `verim-cosmosd add-genesis-account alice 10000000token,100000000stake`
    - Generates genesis node transaction:
        - Command: `verim-cosmosd gentx <key_name> 1000000stake --chain-id <chain_id>`
        - Example: `verim-cosmosd gentx alice 1000000stake --chain-id verim-cosmos`
2. The last participant:
    - Adds genesis node transactions into genesis:
        - Command: `verim-cosmosd collect-gentxs`
    - Verifies genesis:
        - Command: `verim-cosmosd validate-genesis`
    - Shares his genesis with other nodes:
        - Location on the previous participant's machine: `$HOME/.verimcosmos/config/genesis.json`
        - Destination folder on the current participant's machine: `$HOME/.verimcosmos/config/`

After this steps the nodes of all participants have same genesis, and they can connect to each other.

### Running the network

- Each of them:
    - Shares his node's ID@IP with each other
        - How to get this information?
    - 5.1 Updates address book of them node. It allows nodes to connect to each other. Every node owner needs to run follow command: `sed -i "s/persistent_peers = \"\"/persistent_peers = \"$NODE_0_ID@node0:26656,$NODE_1_ID@node1:26656,$NODE_2_ID@node2:26656,$NODE_3_ID@node3:26656\"/g" $NODE_#_HOME/config/config.toml`, where $NODE_#_HOME means variable with number of node owner.
    - // TODO: Research 5.2 Eventually need to set minimal gas prices (it needs for fees transaction). Every node owner needs to run follow command: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_#_HOME/config/app.toml`, where instead `1token` node can set other price.
    - Start node:
        - `verim-cosmosd start`
        - It's better to use process superviser like `systemd`.


Congrats!
