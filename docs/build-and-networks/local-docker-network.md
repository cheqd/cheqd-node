# cheqd-testsnet docker image

## Description

Debian based docker image with the `cheqd-nonded` executable inside preconfigured to run a network of 2 nodes. Intended for use in CI pipelines.

## Prebuilt package

You can find prebuilt package here:

[https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-testnet](https://github.com/cheqd/cheqd-node/pkgs/container/cheqd-testnet)

To pull it use:

```text
docker pull ghcr.io/cheqd/cheqd-testnet:latest
```

## Prerequisites

* Build `cheqd-node` image first. See the [instruction](../setup-and-configure/docker-install.md).

## Building

To build the image:

* Go to the repository root
* Run `docker build -f tests/networks/docker-localnet/Dockerfile -t cheqd-testnet .`

## Running

* Run `docker run -it --rm -p "26657:26657" -p "1317:1317" cheqd-testnet`
* RPC apis are exposed on the folowing ports:
  * node\_0: `26657`
* Try to connect to any node in your browser, for instance: `http://localhost:26657/`

