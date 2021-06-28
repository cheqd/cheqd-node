# Running a new Validator Node

This document describes in details how to configure a validator node, and add it to the existing network.

If a new network needs to be initialized, please follow the Running Genesis Network instructions first. After this more validator nodes can be added by following the instructions from this doc.

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
Current delivery is compiled and tested under `Ubuntu 20.04 LTS` so we recommend using this distribution for now. In future, it will be possible to compile the application for a wide range of operating systems thanks to Go language.

## Deployment steps

New participant:

1. Creates an account:

    - **Generates local keys for his future account:**

        Command `verim-noded keys add <key_name>`

        Example `verim-noded keys add alice`

    - **Asks another member to transfer some tokens:**

        Tokens are used to post transactions. It also used to create a stake for new nodes.

        Another mmber will ask for the address of the new participant. Cosmos account address is a function of public key.

        Use this command to find out your adress and other key information: `verim-noded keys show <key_name>`

2. Initilizes new node and connects it to the network as observer:

    - **Find out `<chain-id>`:**

        Command `verim-noded status --node <remove_node_ip>`
        
        Chain id will showed as `network` property.
        
        Another way is to look into the genesis file. Chain id should be also defined here.

    - **Initializes node config files:**
        
        Command: `verim-noded init <node_name> --chain-id <chain_id>`
        
        Example: `verim-noded init alice-node --chain-id verim-node`
        
    - **Gets genesis:**
        
        Genesis should be published for public networks. If not, you can ask any existing network participant for it.
        
        Location (destination) of the genesis file: `$HOME/.verimnode/config/genesis.json`
        
    - **Updates address book of his node:**
        
        Persistent nodes addresses should be also published publically. If not, you can ask any existing network participant for it.
        
        Open node's config file: `$HOME/.verimnode/config/config.toml`
        
        Search for `persistent_peers` parameter and set it's value to a comma separated list of other participant node addresses.
        
        Format: `<node-0-id>@<node-0-ip>, <node-1-id>@<node-1-ip>, <node-2-id>@<node-2-id>, <node-3-id>@<node-3-id>`.
        
        Domain namaes can be used instead IP adresses.
        
        Example:
        
        ```
        persistent_peers = "d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656, 9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656, abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656, cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656"
        ```

    - **Makes RPC endpoint available externally (optional):**
        
        This step is necessary if you want to allow incoming client applications connections to your node. Otherwise, the node will be accessible only locally. 
        
        Open node configuration file using text editor you prefer: `$HOME/.verimnode/config/config.toml`
        
        Search for `ladr` parameter in `RPC Server Configuration Options` section and replace it's value to `0.0.0.0:26657`
                
        Example: `laddr = "tcp://0.0.0.0:26657"`
        
    - **Start node:**
        
        Command: `verim-noded start`
        
        It's better to use process supervisor like `systemd` to run persistent nodes.
        
3. Promotes the node to the validator:

    - **Post a transaction to te network:**
    
        ```
        verim-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas <amount-gas> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
        ```

        `commission-max-change-rate`, `commission-max-rate` and `commission-rate` may take fraction number as `0.01`

        Use this command to find out `<validator-pubkey>`: `verim-noded tendermint show-validator`. This command **MUST** be run on the node's machine.
        
        Example:
        
        ```
        verim-noded tx staking create-validator --amount 50000000stake --from steward1 --moniker steward1 --chain-id verimnode --min-self-delegation="1" --gas="auto" --gas-prices="1token" --pubkey cosmosvalconspub1zcjduepqpmyzmytdzjhf2fjwttjsrv49t62gdexm2yttpmgzh38p0rncqg8ssrxm2l --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01"
        ```
