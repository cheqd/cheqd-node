# cheqd

cheqd is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

## Building node from source

Prerequisites:

- Install [Go](https://golang.org/doc/install)
- Install [Starport](https://docs.starport.network/intro/install.html)

To build the node executable run:

```
starport chain build
```

To look up binary's location run:

```
which cheqd-noded
```

## Building node in docker

Use this [instruction](ci/docker/README.md).

## Running local network using starport

Prerequisites:

- Install [Go](https://golang.org/doc/install)
- Install [Starport](https://docs.starport.network/intro/install.html)

Only the network of one node is supported. To run the network of one node:

```
starport serve
```

`serve` command installs dependencies, builds, initializes and starts your blockchain in development.

Your blockchain in development can be configured with `config.yml`. To learn more see the [reference](https://github.com/tendermint/starport#documentation).

## Running local network using docker

Use this [instruction](ci/local_net/README.md).

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/W8trcGV)

