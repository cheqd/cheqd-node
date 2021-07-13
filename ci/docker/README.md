# `verim-node` docker image

## Description

Debian based docker image that contains the latest version of `verim-nonded` executable.

## Building

To build the image:

- Go to the repository root;
- Run `docker build -f ci/docker/Dockerfile -t verim-node .`.

## Usage

### verim-noded

`verim-noded` executable is entry point.

Usage:

```
docker run -it --rm verim-node <command> <args>
```

### node-runner

Used to run a node in one command. The following env variable should be defined:

- `NODE_MONIKER` - node moniker;
- `GENESIS` - base64 encoded content of `genesis.json`;
- `NODE_KEY` - base64 encoded content of `node_key.json`;
- `PRIV_VALIDATOR_KEY` - base64 encoded content of `priv_validator_key.json`;
- `NODE_ARGS` (optional) - argument string passed to the `verim-noded start` command.

Usage:

```
docker run -it --rm --entrypoint node-runner -e NODE_MONIKER=<moniker> -e GENESIS="<content>" -e NODE_KEY="<content>" -e PRIV_VALIDATOR_KEY="<content>" verim-node
```
