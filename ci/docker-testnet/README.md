# `verim-testsnet` docker image

## Description

Debian based docker image with the latest version of `verim-nonded` executable and preconfigured network of 4 nodes. Intended for use in CI pipelines.

## Prerequisites

- Build `verim-node` image first. See the [instruction](../docker/README.md).

## Building

To build the image:

- Go to the repository root
- Run `docker build -f ci/docker-testnet/Dockerfile -t verim-testnet .`
