# Creating a new network from genesis

This document describes in details how to configure a genesis network with any amount of participants.

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

Current delivery is compiled and tested for `Ubuntu 20.04 LTS` so we recommend using this distribution for now. In the future, it will be possible to compile the application for a wide range of operating systems thanks to the Go language.

## Deployment steps

### Generating genesis file

1. Participants must choose <chain_id> for the network.
2. Each participant (one by one) should:
    
    - **Generates local keys for the future account:**
    
        Command `verim-noded keys add <key_name>`

        Example `verim-noded keys add alice`
    
    - **Initializes node config files:**
        
        Command: `verim-noded init <node_name> --chain-id <chain_id>`
        
        Example: `verim-noded init alice-node --chain-id verim-node`
        
    - **(Each participant except the first one) Gets genesis from the previous participant:**
        
        Location on the previous participant's machine: `$HOME/.verimnode/config/genesis.json`
        
        Destination folder on the current participant's machine: `$HOME/.verimnode/config/`
        
    - **(Each participant except the first one) Gets genesis node transactions from the previous participant:**
        
        Location on the previous participant's machine: `$HOME/.verimnode/config/gentx/`
        
        Destination folder on the current participant's machine: `$HOME/.verimnode/config/gentx/`
                
    - **Adds a genesis account with a public key:**
        
        Command: `verim-noded add-genesis-account <key_name> 10000000token,100000000stake`
        
        Example: `verim-noded add-genesis-account alice 10000000token,100000000stake`
        
    - **Generates genesis node transaction:**
        
        Command: `verim-noded gentx <key_name> 1000000stake --chain-id <chain_id>`
        
        Example: `verim-noded gentx alice 1000000stake --chain-id verim-node`
        
        **TODO: Node owner should specify gas prices here. Need to research how it works.**
        
3. The last participant:

    - **Adds genesis node transactions into genesis:**
        
        Command: `verim-noded collect-gentxs`
        
    - **Verifies genesis:**
        
        Command: `verim-noded validate-genesis`
        
    - **Shares genesis with other nodes:**
        
        Location on the last participant's machine: `$HOME/.verimnode/config/genesis.json`
        
        Destination folder on the other participant's machines: `$HOME/.verimnode/config/`

After this steps:
- Nodes of all participants have the same genesis;
- The genesis contains:
    - Accounts of all participants (genesis accounts);
    - Node creation transactions from all participants (genesis nodes).

### Running the network

- Each participant:

    - **Shares his node ID and IP with the others:**
        
        Command to find out node's id: `verim-noded tendermint show-node-id`. This command **MUST** be run on the machine where node's config files are located.
        
        Node IP is external IP of the node's machine.
        
        Node adress is the combination of IP and ID in the following format: `ID@IP`.
        
        Port is the RPC adress of the node. It can be configured here: `$HOME/.verimnode/config/config.toml`. Default value is `26656`.
        
        Node address example: `d45dcc54583d6223ba6d4b3876928767681e8ff6@192.168.0.142:26656`
        
    - **Update address book of the node:**
        
        Open node's config file: `$HOME/.verimnode/config/config.toml`
        
        Search for `persistent_peers` parameter and set it's value to a comma separated list of other participant node addresses.
        
        Format: `<node-0-id>@<node-0-ip>, <node-1-id>@<node-1-ip>, <node-2-id>@<node-2-id>, <node-3-id>@<node-3-id>`.
        
        Domain names can be used instead of IP adresses.
        
        Example:
        
        ```
        persistent_peers = "d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656, 9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656, abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656, cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656"
        ```
        
    - **Makes RPC endpoint available externally (optional):**
        
        This step is necessary if you want to allow incoming client applications connections to your node. Otherwise, the node will be accessible only locally. 
        
        Open node configuration file using the text editor that you prefer: `$HOME/.verimnode/config/config.toml`
        
        Search for `ladr` parameter in `RPC Server Configuration Options` section and replace it's value to `0.0.0.0:26657`
                
        Example: `laddr = "tcp://0.0.0.0:26657"`
        
    - **Start node:**
        
        Command: `verim-noded start`
        
        It's better to use process supervisor like `systemd` to run persistent nodes.


Congratulations the network is running!
