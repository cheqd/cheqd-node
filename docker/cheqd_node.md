# cheqd-node docker image

## Description

Debian based docker image that contains the latest version of `cheqd-nonded` executable.

## Prebuilt package

You can find prebuilt package here:

[https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-node](https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-node)

To pull it use:

```text
docker pull ghcr.io/cheqd/cheqd-node:latest
```

## Building

To build the image:

* Go to the repository root;
* Run `docker build -f docker/cheqd_node/Dockerfile -t cheqd-node .`.

Default home directory for `cheqd` user is `/cheqd`. It can be overridden via `CHEQD_HOME_DIR` build argument. Example: `--build-arg CHEQD_HOME_DIR=/home/cheqd`.

Note: If you are using M1 Macbook you should modify the FROM statement in the Dockerfile, should be like this "FROM --platform=linux/amd64 golang:buster as builder "

## Usage

### cheqd-noded

`cheqd-noded` executable is entry point by default.

Usage:

```text
docker run -it --rm cheqd-node <command> <args>
```

### node-runner

Used to initialize configuration files and run a node in one command.

Parameters that can be passed via environment variables:

* `NODE_MONIKER` - node moniker;
* `GENESIS` - base64 encoded content of `genesis.json`;
* `NODE_KEY` - base64 encoded content of `node_key.json`;
* `PRIV_VALIDATOR_KEY` - base64 encoded content of `priv_validator_key.json`;
* `NODE_ARGS` \(optional\) - argument string passed to the `cheqd-noded start` command.

Additional parameters that will be applied via `cheqd-noded configure`:

* `CREATE_EMPTY_BLOCKS`
* `FASTSYNC_VERSION`
* `MIN_GAS_PRICES`
* `RPC_LADDR`
* `P2P_EXTERNAL_ADDRESS`
* `P2P_LADDR`
* `P2P_MAX_PACKET_MSG_PAYLOAD_SIZE`
* `P2P_PERSISTENT_PEERS`
* `P2P_RECV_RATE`
* `P2P_SEED_MODE`
* `P2P_SEEDS`
* `P2P_SEND_RATE`

Usage:

```text
docker run -it --rm --entrypoint node-runner -e NODE_MONIKER=<moniker> -e GENESIS="<content>" -e NODE_KEY="<content>" -e PRIV_VALIDATOR_KEY="<content>" cheqd-node
```
