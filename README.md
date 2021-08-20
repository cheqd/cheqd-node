# cheqd

cheqd is a purpose-built network for decentralised identit built using [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Tendermint](https://github.com/tendermint/tendermint).

## Overview
Getting started as a node operator on the cheqd network requires the following steps:

1. Install a node (on hosting platform of your choice)
2. Configure a node to join a specific network (currently only a _testnet_
3. Use node functionality and upgrade after installation

## Installation
Follow the instructions to [set up a new node](docs/setting-up-a-new-node.md). This document covers:
1. Minimum system requirements and pre-requisities
2. Installation process using Debian (.deb) package, binary, and Docker
3. Fetching basic node information after installation

There are packaged releases available for node installation, depending on the method you prefer.

### Install using (.deb) package releases
* Make sure you've followed through how to [set up a new node](docs/setting-up-a-new-node.md)
* Undertand an [overview of what the Debian (.deb) package configures](docs/deb-package-overview.md)
* If you already have an existing node installed using the .deb package, follow these steps on how to [upgrade your node using the .deb package](docs/deb-package-upgrade.md)

### Install using Docker

## Network Configuration
  * [Setting up a new validator](docs/setting-up-a-new-validator.md)
  * [Setting up a new network](docs/setting-up-a-new-network.md)

## Usage
Once installed, **cheqd-node** can be controlled using the [cheqd Cosmo CLI reference guide](docs/cosmos-cli.md). At present, this supports configuring accounts and writing Decentralized Identifier (DID) entries on ledger. We plan on expand functionality in future releases.

## Development
**cheqd-node** is created with [Starport](https://github.com/tendermint/starport). If you want to build a node from source or contribute to the code, please read our guide to [building and testing](docs/building-and-testing.md).


## Community

### Discuss the cheqd project, get help, and support
The [cheqd Community Slack is open for anyone to join](http://cheqd.link/join-cheqd-slack) and is our discussion forum for the open-source community, software developers, and node operators. Please reach out to us for real-time discussions and help.

You can also follow the cheqd on other social channels:
- [Twitter](https://twitter.com/cheqd_io)
- [Telegram](https://t.me/cheqd) (with a separate [announcements-only channel](https://t.me/cheqd_announcements))
- [Medium](https://blog.cheqd.io/)
- [LinkedIn](http://cheqd.link/linkedin)

### Learn more about the wider Cosmos ecosystem

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Cosmos Discord chat](https://discord.gg/W8trcGV)
