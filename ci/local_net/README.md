# Localnet

## Description

The set of scripts to generate configuration for a testnet of four nodes and run it locally in docker-compose.

## Prerequisites

- [Starport](https://docs.starport.network/intro/install.html) 
- docker-compose

## How to run

1.  Build docker image:

    See [the instruction](../docker/README.md)

2. Build verim-noded:

    ```
    starport build
    ```

3. Generate node configurations:

    Run: `gen_node_configs.sh` or `gen_node_configs.sh`.

4. Run docker-compose:

    ```
    docker-compose up
    ```

## Result

This will setup 4 nodes listening on the following ports:

- Node0:
    - p2p: 26656
    - rpc: 26657
- Node1:
    - p2p: 26666
    - rpc: 26667
- Node2:
    - p2p: 26676
    - rpc: 26677
- Node3:
    - p2p: 26686
    - rpc: 26687

You can tests connection to a node using browser: `http://localhost:<rpc_port>`. Example for the fitst node: `http://localhost:26657`.

When connecting using CLI, point path to home directory: `--home localnet/client`. This directory contains keys from genesis acounts.

## Demo commands:

### Show balances

```
vc query bank balances (vc keys show anna -a --home localnet/client) --home localnet/client
```

### Create NYM

```
vc tx verim create-nym "alias" "verkey" "did" "role" --from anna --gas-prices 1token --chain-id verim-node-chain --home localnet/client
```

### List nym

```
vc query verim list-nym --home localnet/client
```
