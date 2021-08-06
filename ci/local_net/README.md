# Localnet

## Description

The set of scripts to generate configuration for a testnet of four nodes and run it locally in docker-compose.

## Prerequisites

- [Starport](https://docs.starport.network/intro/install.html) 
- docker-compose

## How to run

1.  Build docker image:

    See [the instruction](../docker/README.md).

2. Build cheqd-noded:

    ```
    starport chain build
    ```

3. Generate node configurations:

    Run: `gen_node_configs_demo.sh`.

4. Run docker-compose:

    Run: `run_docker.sh`.

## Result

This will setup 4 nodes listening on the following ports:

- Node0:
    - p2p: 26656
    - rpc: 26657
- Node1:
    - p2p: 26659
    - rpc: 26660
- Node2:
    - p2p: 26662
    - rpc: 26663
- Node3:
    - p2p: 26665
    - rpc: 26666
    
- Client's with configs in `node_configs/client` directory. 

Also, there are 4 accounts for operators:
- operator0;
- operator1;
- operator2;
- operator3;

You can tests connection to a node using browser: `http://localhost:<rpc_port>`. Example for the fitst node: `http://localhost:26657`.

When connecting using CLI, point path to home directory: `--home node_configs/client`. This directory contains keys from genesis acounts.

## Demo commands:

### Show balances

```
cheqd-noded query bank balances $(cheqd-noded keys show operator0 -a --home node_configs/client) --home node_configs/client
```

### Create NYM

```
cheqd-noded tx cheqd create-nym "alias" "verkey" "did" "role" --from operator0 --gas-prices 1cheq --chain-id cheqd --home node_configs/client
```

### List nym

```
cheqd-noded query cheqd list-nym --home node_configs/client
```
