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
    - Makes RPC endpoint available externally (optional, allows cliens to connect to the node):
        - Command: `sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "<external-ip-address>"/g' $HOME/.verimcosmos/config/config.toml`
        - Example: `sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $HOME/.verimcosmos/config/config.toml`
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
    - Shares his node's ID@IP with each other:
        - Find out own id-node: `verim-cosmosd tendermint show-node-id`
        - Node IP matches to `[rpc] laddr` field in `$HOME/.verimcosmos/config/config.toml`
    -  Updates address book of own node:
        - Command:
            ```
            sed -i "s/persistent_peers = \"\"/persistent_peers = \"
            <node-0-id>@<node-0-ip>,
            <node-1-id>@<node-1-ip>,
            <node-2-id>@<node-2-id>,
            <node-3-id>@<node-3-id>\"/g" $HOME/.verimcosmos/config/config.toml
            ```
        - Example:
            ```
            sed -i "s/persistent_peers = \"\"/persistent_peers = \"d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656,9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656,abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656,cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656\"/g" $HOME/.verimcosmos/config/config.toml
            ```
    - // TODO: Research 5.2 Eventually need to set minimal gas prices (it needs for fees transaction). Every node owner needs to run follow command: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "1token"/g' $NODE_#_HOME/config/app.toml`, where instead `1token` node can set other price.
    - Start node:
        - `verim-cosmosd start`
        - It's better to use process supervisor like `systemd`.


Congrats!
