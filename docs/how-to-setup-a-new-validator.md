# Running a new Validator Node

This document describes in detail how to configure a validator node, and add it to the existing network.

If a new network needs to be initialized, please first follow the instructions for [creating a new network from genesis](how-to-deploy-genesis-network.md). After this, more validator nodes can be added by following the instructions from this document.

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

### Binary distribution

There are several ways to get binary:

- Compile from source code - [instruction](../README.md);
- Get `tar` archive with the binary compiled for Ubuntu 20.04 in [releases](https://github.com/cheqd-id/cheqd-node/releases); <-- Recommended
- Get `deb` for Ubuntu 20.04 in [releases](https://github.com/cheqd-id/cheqd-node/releases);
- Get docker image form [packages](https://github.com/cheqd-id/cheqd-node/pkgs/container/cheqd-node).

## Node deployment

Follow these steps to deploy a new node:

1. Setup a server that satisfies [hardware requirements](#hardware-requirements) and [operating system requirements](#operating-system);

    More about hardware requirements can be found [here](https://docs.tendermint.com/master/nodes/running-in-production.html#hardware).

2. Get the binary using one of the [described ways](#binary-distribution);

    It's recommended to put the binary to the location which is in PATH.

    Example:

    ```
    cp cheqd-noded /usr/bin
    ```

3. Initialize node config files:
        
    Command: `cheqd-noded init <node_name>`
    
    Example: `cheqd-noded init alice-node`
        
4. Set genesis:
        
    Genesis should be published for public networks. If not, you can ask any existing network participant for it.
    
    Location (destination) of the genesis file: `$HOME/.cheqdnode/config/genesis.json`
        
5. Set persistent peers:
        
    Persistent nodes addresses should also be published publically. If not, you can ask any existing network participant for it.
    
    Open node's config file: `$HOME/.cheqdnode/config/config.toml`
    
    Search for `persistent_peers` parameter and set it's value to a comma separated list of other participant node addresses.
    
    Format: `<node-0-id>@<node-0-ip>, <node-1-id>@<node-1-ip>, <node-2-id>@<node-2-id>, <node-3-id>@<node-3-id>`.
    
    Domain names can be used instead of IP adresses.
    
    Example:
    
    ```
    persistent_peers = "d45dcc54583d6223ba6d4b3876928767681e8ff6@node0:26656, 9fb6636188ad9e40a9caf86b88ffddbb1b6b04ce@node1:26656, abbcb709fb556ce63e2f8d59a76c5023d7b28b86@node2:26656, cda0d4dbe3c29edcfcaf4668ff17ddcb96730aec@node3:26656"
    ```

6. (optional) Make RPC endpoint available externally:
     
    This step is necessary if you want to allow incoming client application connections to your node. Otherwise, the node will be accessible only locally. 

    Open the node configuration file using the text editor that you prefer: `$HOME/.cheqdnode/config/config.toml`

    Search for `ladr` parameter in `RPC Server Configuration Options` section and replace it's value to `0.0.0.0:26657`
        
    Example: `laddr = "tcp://0.0.0.0:26657"`

7. Configure firewall rules:

    Allow incoming tcp connections on the P2P port - `26656` by default.

    If you made RPC endpoint available externally, allow incoming tcp connections on the RPC port - `26657` by default.

    Allow all outgoing tcp connections for P2P communication. You can restrict port to the default P2P port `26656` but your node will not be able to connect to nodes with non default P2P port in this case.

8. Start node:

    Command: `cheqd-noded start`

    It's highly recommended to use a process supervisor like `systemd` to run persistent nodes.

9. (optional) Setup sentry nodes for DDOS protection:

    You can read about sentry nodes [here](https://docs.tendermint.com/master/nodes/validators.html).

10. (optional) Setup cosmovisor for automatic updates:

    You can read about sentry nodes [here](https://docs.cosmos.network/master/run-node/cosmovisor.html).

11. (optional) Read other advices about running node in production:

    You can read advices [here](https://docs.tendermint.com/master/nodes/running-in-production.html).

## Network configuration

Follow these steps to promote the deployed node to a validator:

1. Create an account:

    - **Generate local keys for the future account:**

        Command: `cheqd-noded keys add <key_name>`

        Example: `cheqd-noded keys add alice`

    - **Ask another member to transfer some tokens:**

        Tokens are used to post transactions. It also used to create a stake for new nodes.

        Another mmber will ask for the address of the new participant. Cosmos account address is a function of the public key.

        Use this command to find out your adress and other key information: `cheqd-noded keys show <key_name>`

3. Promote the node to the validator:

    - **Post a transaction to the network:**
    
        ```
        cheqd-noded tx staking create-validator --amount <amount-staking> --from <key-name> --chain-id <chain-id> --min-self-delegation <min-self-delegation> --gas <amount-gas> --gas-prices <price-gas> --pubkey <validator-pubkey> --commission-max-change-rate <commission-max-change-rate> --commission-max-rate <commission-max-rate> --commission-rate <commission-rate>
        ```

        `commission-max-change-rate`, `commission-max-rate` and `commission-rate` may take fraction number as `0.01`

        Use this command to find out the `<validator-pubkey>`: `cheqd-noded tendermint show-validator`. This command **MUST** be run on the node's machine.
        
        Example:
        
        ```
        cheqd-noded tx staking create-validator --amount 50000000stake --from steward1 --moniker steward1 --chain-id cheqdnode --min-self-delegation="1" --gas="auto" --gas-prices="1token" --pubkey cosmosvalconspub1zcjduepqpmyzmytdzjhf2fjwttjsrv49t62gdexm2yttpmgzh38p0rncqg8ssrxm2l --commission-max-change-rate="0.02" --commission-max-rate="0.02" --commission-rate="0.01"
        ```
