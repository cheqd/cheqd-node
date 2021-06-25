# `verim-testsnet` docker image

## Description

Debian based docker image with the latest version of `verim-nonded` executable and preconfigured network of 4 nodes. Intended for use in CI pipelines.

## Prerequisites

- Build `verim-node` image first. See the [instruction](../docker/README.md).

## Building

To build the image:

- Go to the repository root
- Run `docker build -f ci/docker-testnet/Dockerfile -t verim-testnet .`

## Running

- Run `docker run -it --rm -p "26657:26657" -p "26660:26660" -p "26663:26663" -p "26666:26666" verim-testnet`
- RPC posta are exposed on the folowing ports:
  - node_0: `26657`
  - node_1: `26660`
  - node_2: `26663`
  - node_3: `26666`
- Try to connect to any node in your browser, for instance: ``
