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

1. Participants must choose <chain_id> for the network.
2. Each participant (one by one):
    
    - Generates local keys for his future account:  
    
        Comman `verim-cosmosd keys add <key_name>`

        Examp `verim-cosmosd keys add alice`
    
    - Initializes node config files:
        
        Command: `verim-cosmosd init <node_name> --chain-id <chain_id>`
        
        Example: `verim-cosmosd init alice-node --chain-id verim-cosmos`
        
    - (Each participatn except the first one) Gets genesis from the previous participant:
        
        Location on the previous participant's machine: `$HOME/.verimcosmos/config/genesis.json`
        
        Destination folder on the current participant's machine: `$HOME/.verimcosmos/config/`
        
    - (Each participatn except the first one) Gets genesis node transactions form the previous participant.
        
        Location on the previous participant's machine: `$HOME/.verimcosmos/config/gentx/`
        
        Destination folder on the current participant's machine: `$HOME/.verimcosmos/config/gentx/`
                
    - Adds a genesis account with his public key:
        
        Command: `verim-cosmosd add-genesis-account <key_name> 10000000token,100000000stake`
        
        Example: `verim-cosmosd add-genesis-account alice 10000000token,100000000stake`
        
    - Generates genesis node transaction:
        
        Command: `verim-cosmosd gentx <key_name> 1000000stake --chain-id <chain_id>`
        
        Example: `verim-cosmosd gentx alice 1000000stake --chain-id verim-cosmos`
        
        **TODO: Node owner should specify gas prices here. Need to research how it works.**
        
3. The last participant:

    - Adds genesis node transactions into genesis:
        
        Command: `verim-cosmosd collect-gentxs`
        
    - Verifies genesis:
        
        Command: `verim-cosmosd validate-genesis`
        
    - Shares his genesis with other nodes:
        
        Location on the last participant's machine: `$HOME/.verimcosmos/config/genesis.json`
        
        Destination folder on the other participant's machines: `$HOME/.verimcosmos/config/`

After this steps:
- Nodes of all participants have same genesis;
- The genesis contains:
    - Accounts of all participants (genesis accounts);
    - Node creation transactions from all participants (genesis nodes).

### Running the network

- Each participant:

    - Shares his node ID and IP with each other:
        
        Command to find out node's id: `verim-cosmosd tendermint show-node-id`. This command **MUST** be run on the machine where node's config files are located.
        
        Node IP is external IP of the node's machine.
        
        Node adress is the combination of IP and ID in the following format: `ID@IP`.
        
        Port is the RPC adress of the node. It can be configured here: `$HOME/.verimcosmos/config/config.toml`. Default value is `26656`.
        
        Node address example: `d45dcc54583d6223ba6d4b3876928767681e8ff6@192.168.0.142:26656`
        
    - Update address book of his node:
        
        Open node's config file: `$HOME/.verimcosmos/config/config.toml`
        
        Search for `persistent_peers` parameter and set it's value to a comma separated list of other participant node addresses.
        
        Format: `<node-0-id>@<node-0-ip>, <node-1-id>@<node-1-ip>, <node-2-id>@<node-2-id>, <node-3-id>@<node-3-id>`.
        
        Domain namaes can be used instead IP adresses.
        
        Example:
        
        ```
        persistent_peers = "d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656, 9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656, abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656, cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656"
        ```
        
    - Makes RPC endpoint available externally (optional):
        
        This step is necessary if you want to allow incoming client applications connections to your node. Otherwise, the node will be accessible only locally. 
        
        Open node configuration file using text editor you prefer: `$HOME/.verimcosmos/config/config.toml`
        
        Search for `ladr` parameter in `RPC Server Configuration Options` section and replace it's value to `0.0.0.0:26657`
                
        Example: `laddr = "tcp://0.0.0.0:26657"`
        
    - Start node:
        
        Command: `verim-cosmosd start`
        
        It's better to use process supervisor like `systemd` to run persistent nodes.


Congratulations the network is running!
