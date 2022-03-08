# Setting up a new cheqd node

## Context

This document describes how to use install and configure a new instance of `cheqd-node` from pre-built packages and adding it to an existing network (such as the cheqd testnet) as an observer or validator.

For other scenarios, please see [setting up a new network from scratch](../build-and-networks/manual-network-setup.md) and [building `cheqd-node` from source](../build-and-networks/README.md).

## Pre-requisites

### Hardware requirements

For most nodes, the RAM/vCPU requirements are relatively static and do not change over time. However, the disk storage space needs to grow as the chain grows and will evolve over time.

It is recommended to provide the disk storage as an expandable volume/partition that is mounted on your node configuration data path (the default is under `/home/cheqd`) so that it can be expanded independent of the root volume.

Extended information on [recommended hardware requirements is available in Tendermint documentation](https://docs.tendermint.com/v0.35/nodes/running-in-production.html#hardware). The figures below have been updated from the default Tendermint recommendations to account for current cheqd network chain size, real-world usage accounting for requests nodes need to handle, etc.

#### Minimum specifications

* 2 GB RAM
* x64 1.4 GHz 1 vCPU (or equivalent)
* 120 GB of disk space

#### Recommended specifications

* 4 GB RAM
* x64 2.0 GHz 2 vCPU (or equivalent)
* 150 GB SSD

### Operating system

Our [packaged releases](https://github.com/cheqd/cheqd-node/releases) are currently compiled and tested for `Ubuntu 20.04 LTS`, which is the recommended operating system for installation using Debian package or binaries.

For other operating systems, we recommend using [pre-built Docker image releases for `cheqd-node`](https://github.com/orgs/cheqd/packages?repo_name=cheqd-node).

We plan on supporting other operating systems in the future, based on demand for specific platforms by the community.

### Storage volumes

We recommend using a storage path that can be kept persistent and restored/remounted (if necessary) for the configuration, data, and log directories associated with a node. This allows a node to be restored along with configuration files such as node keys and for the node's copy of the ledger to be restored without triggering a full chain sync.

The default directory location for `cheqd-node` installations is `$HOME/.cheqdnode`, which computes to `/home/cheqd/.cheqdnode` when [using the Debian package installer](debian/README.md). Custom paths can be defined if desired.

### Ports

To function properly, `cheqd-node` requires two types of ports to be configured. Depending on the setup, you may also need to configure firewall rules to allow the correct ingress/egress traffic.

Node operators should ensure there are no existing services running on these ports before proceeding with installation.

#### P2P port

The P2P port is used for peer-to-peer communication between nodes. This port is used for your node to discover and connect to other nodes on the network. It should allow traffic to/from any IP address range.

* By default, the P2P port is set to `26656`.
* Inbound TCP connections on port `26656` (or your custom port) should be allowed from *any* IP address range.
* Outbound TCP connections must be allowed on *all* ports to *any* IP address range.
* The default P2P port can be changed in `$HOME/.cheqdnode/config/config.toml`.

Further details on [how P2P settings work is defined in Tendermint documentation](https://docs.tendermint.com/v0.35/nodes/configuration.html#p2p-settings).

#### RPC port

The RPC port is intended to be used by client applications to interact with the node.

During node configuration for cheqd testnet, the RPC port is needed to transfer the initial balance of `cheq` tokens required for staking.

Beyond this stage, it is up to a node operator whether they want this port to be exposed to the public internet.

The [RPC endpoints for a node](https://docs.tendermint.com/master/rpc/) provide REST, JSONRPC over HTTP, and JSONRPC over WebSockets. These API endpoints can provide useful information for node operators, such as healthchecks, network information, validator information etc.

* By default, the RPC port is set to `26657`
* Inbound and outbound TCP connections should be allowed from destinations desired by the node operator. The default is to allow this from any IPv4 address range.
* TLS for the RPC port can also be setup separately. Currently, TLS setup is not automatically carried out in the install process described below.
* The default RPC port can be changed in `$HOME/.cheqdnode/config/config.toml`.

### Sentry nodes (optional)

Tendermint allows more complex setups in production, where the ingress/egress to a validator node is [proxied behind a "sentry" node](https://docs.tendermint.com/v0.35/nodes/validators.html#setting-up-a-validator).

While this setup is not compulsory, node operators with higher stakes or a need to have more robust network security may consider setting up a sentry-validator node architecture.

## Installing and configuring a cheqd node

Follow the guide for your preferred installation method:

* [Debian package install](debian/deb-package-install.md)
* [Docker install](docker-install.md)
* [Binary install](binary-install.md)

[Configure your node as a validator](../validator-guide/README.md) after successful installation.

## Further information

* Tendermint documentation has [best practices for running a Cosmos node in production](https://docs.tendermint.com/v0.35/nodes/running-in-production.html).
* [Сosmovisor could be used for automatic upgrades](https://docs.cosmos.network/master/run-node/cosmovisor.html); however in our testing so far this method has not been reliable and is therefore currently not recommended.
