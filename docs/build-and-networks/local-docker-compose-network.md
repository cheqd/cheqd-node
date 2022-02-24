# Docker Compose Based Localnet

## Description

The set of scripts to generate configurations for a network of four nodes and run it locally in docker-compose.

## Prerequisites

* [Starport](https://docs.starport.network/guide/install.html)
* docker-compose

## How to run

1. Build docker image:

    See [the instruction](../setup-and-configure/docker-install.md).

2. Build cheqd-node:

    ```bash
    starport chain build
    ```

3. Generate node configurations:

    Run: `gen_node_configs.sh`.

4. Run docker-compose:

    Run: `run_docker.sh`.

## Result

### Nodes

This will setup 4 nodes listening on the following ports:

* Node0:
  * p2p: 26656
  * rpc: 26657
* Node1:
  * p2p: 26659
  * rpc: 26660
* Node2:
  * p2p: 26662
  * rpc: 26663
* Node3:
  * p2p: 26665
  * rpc: 26666

You can tests connection to a node using browser: `http://localhost:<rpc_port>`. Example for the first node: `http://localhost:26657`.

### Accounts

Also, there will be 4 keys generated and corresponding genesis accounts created for node operators:

* operator0;
* operator1;
* operator2;
* operator3;

When connecting using CLI, point path to home directory of the node corresponding to the operator: `--home network-config/validator-x`.

## CLI commands

See [the cheqd CLI guide](../cheqd-cli/README.md) to learn about the most common CLI flows.
