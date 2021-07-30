# Running a new node

This document describes in detail how to configure infrastructure and deploy a new node (observer or validator).

If a new network needs to be initialized, please first follow the instructions for [creating a new network from genesis](how-to-setup-a-new-network.md).

If a new validator needs to be added to the existing network, please refer to [joining existing network](how-to-join-existing-network.md) instruction.

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
- Get `tar` archive with the binary compiled for Ubuntu 20.04 in [releases](https://github.com/cheqd/cheqd-node/releases);
- Get `deb` for Ubuntu 20.04 in [releases](https://github.com/cheqd/cheqd-node/releases);
- Get docker image form [packages](https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-node).

The most preferable way to get `cheqd-node` is to use `.deb` package. Detailed information about it can be found [here](#deb-package-installation.md)  
## Node deployment

Follow these steps to deploy a new node:

1. Setup a server that satisfies [hardware requirements](#hardware-requirements) and [operating system requirements](#operating-system);

    More about hardware requirements can be found [here](https://docs.tendermint.com/master/nodes/running-in-production.html#hardware).

2. In the case of using tarball, put the binary to the location which is in PATH.

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
   
    8.1 In case of using tarball: `cheqd-noded start`
   
     8.2 In case of using `.deb` package:`systemctl start cheqd-noded.service`

9. (optional) Setup sentry nodes for DDOS protection:

    You can read about sentry nodes [here](https://docs.tendermint.com/master/nodes/validators.html).

10. (optional) Setup cosmovisor for automatic updates:

    You can read about sentry nodes [here](https://docs.cosmos.network/master/run-node/cosmovisor.html).

11. (optional) Read other advices about running node in production:

    You can read advices [here](https://docs.tendermint.com/master/nodes/running-in-production.html).

## Getting node info

### Node id

Node id is a part of peer info. To get `node id` run the following command on the node's machine:

```
cheqd-noded tendermint show-node-id
```

### Validator public key

Validator public key is used to promote node to the validator. To get it run the following command on the node's machine:

```
cheqd-noded tendermint show-validator
```

### Peer information

Peer info is used to connect to peers when setting up a new node. It has the following format:

```
<node-id>@<node-url>
```

Example:

```
ba1689516f45be7f79c7450394144711e02e7341@3.13.19.41:26656
```
