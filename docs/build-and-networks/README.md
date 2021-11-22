# Building and testing

## Building `cheqd-node` from source

### Prerequisites

* Install [Go](https://golang.org/doc/install)
* Install [Starport](https://docs.starport.network/guide/install.html)

To build the `cheqd-node` executable run:

```bash
starport chain build
```

To look up binary's location run:

```text
which cheqd-noded
```

## Building node in docker

Use this [instruction](build-and-networks.md).

## Running local network using starport

Prerequisites:

* Install [Go](https://golang.org/doc/install)
* Install [Starport](https://docs.starport.network/guide/install.html)

Only the network of one node is supported. To run the network of one node:

```text
starport serve
```

`serve` command installs dependencies, builds, initializes and starts your blockchain in development.

Your blockchain in development can be configured with `config.yml`. To learn more see the [reference](https://github.com/tendermint/starport#documentation).

## Running local network using docker compose

Use the [Docker Compose instructions](docker-compose.md).
