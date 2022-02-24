# Building `cheqd-noded`

## Building from source

Prerequisites:

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

## Building in docker

Use this [instruction](../setup-and-configure/docker-install.md).

# Running a network

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

## Running local network in single docker image

Use the [Docker localnet instructions](local-docker-network.md).

## Running local network using docker compose

Use the [Docker Compose localnet instructions](local-docker-compose-network.md).
