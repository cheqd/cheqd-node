# Single Node Localnet

## Description

The script to generate configuration for a network of one node and run it locally without docker.

> Warning: this script removes all files in `$HOME/.cheqdnode` directory. Backup it first.

## Prerequisites

* [Starport](https://docs.starport.network/guide/install.html)

## How to run

1. Build cheqd-noded:

   ```text
    starport chain build
   ```

2. Generate network configuration:

   Run: `gen_node_config.sh`.

3. Run single node network:

   Run: `cheqd-noded start`.

## Result

### Nodes

This will setup 1 node listening on the following ports:

* p2p: 26656
* rpc: 26657

You can tests connection to the node using browser: `http://localhost:<rpc_port>`. Example for the first node: `http://localhost:26657`.

### Accounts

Also, there will be 1 key generated and corresponding genesis accounts created for node operator:

* node_operator;

## CLI commands

See [the reference](https://github.com/cheqd/cheqd-node/blob/main/docs/cheqd-cli/README.md) to learn about the most common CLI flows.
