# cheqd: Node Documentation

**cheqd** is a purpose-built network for decentralised identity.

`cheqd-node` is the server/node portion of the cheqd network stack, built using [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Tendermint](https://github.com/tendermint/tendermint).

## Quick start for joining cheqd mainnet

Getting started as a node operator on the cheqd network mainnet requires the following steps:

1. [Install the `cheqd-node` software](docs/setup-and-configure/readme.md) on a hosting platform of your choice.
2. When you have a node successfully installed, please fill out our [**node operator onboarding form**](http://cheqd.link/mainnet-onboarding). You will need to have the following details on hand to fill out the form:
   1. Node ID for your node
   2. IP address / DNS record that points to the node \(if you're using an IP address, a static IP is recommended\)
   3. Peer-to-peer \(P2P\) connection port \(defaults to `26656`\)
   4. Validator account address (begins with `cheqd`)
   5. Moniker (Nickname/moniker that is set for your mainnet node)
3. Once you have received or purchased your tokens, [promote your node to a validator](docs/setup-and-configure/configure-new-validator.md).
4. If successfully configured, your node would become the latest validator on the cheqd Mainnet! Say hi to the other node operators on our [**cheqd Community Slack**](http://cheqd.link/join-cheqd-slack)

Any time you have questions or need support, feel free to reach out through the [ask for help](https://cheqd-community.slack.com/archives/C02AQ9UK4HY) channel and either the cheqd team or one of your fellow node operators will be happy to offer some guidance. 

## Usage

Once installed, `cheqd-node` can be controlled using the [cheqd Cosmos CLI guide](docs/cheqd-cli/readme.md).

### Currently supported functionality

* Basic token functionality for holding and transferring tokens to other accounts on the same network
* Creating, managing, and configuring accounts and keys on a cheqd node
* Staking and participating in public-permissionless governance
* Governance framework for public-permissionless self-sovereign identity networks
* DID method specification

### Upcoming functionality

A non-exhaustive list of future planned functionality \(not necessarily in order of priority\) is highlighted below:

* Creating and querying DIDDocs
* Creating and managing Verifiable Credentials anchored to DIDs on cheqd mainnet

We plan on adding new functionality rapidly and on a regular basis. We will be sharing regular updates through our **Live Product Updates** page which includes our product roadmap, release notes, node operator FAQs and more. We welcome feedback on our [cheqd Community Slack](http://cheqd.link/join-cheqd-slack) workspace.

## Building from source

`cheqd-node` is created with [Starport](https://github.com/tendermint/starport). If you want to build a node from source or contribute to the code, please read our guide to [building and testing](https://github.com/cheqd/cheqd-node/tree/f74ec3e0ad08adcf2e4173de80dbd9442edc337e/docs/building-and-testing.md).

### Creating a local network

If you are building from source, or otherwise interested in running a local network, we have [instructions on how to set up a new network](https://github.com/cheqd/cheqd-node/tree/f74ec3e0ad08adcf2e4173de80dbd9442edc337e/docs/setting-up-a-new-network.md) for development purposes.

## Community

The [**cheqd Community Slack**](http://cheqd.link/join-cheqd-slack) is our chat channel for the open-source community, software developers, and node operators.

Please reach out to us there for discussions, help, and feedback on the project.

## Bug Reporting & New Feature Requests 

If you notice anything not behaving how you expected, or would like to make a suggestion / request for a new feature, please submit a **bug_report** or **feature_request**  by creating a [**New Issue**](https://github.com/cheqd/cheqd-node/issues/new/choose) and selecting the relevant template. 


### Social media

Follow the cheqd team on our social channels for news, announcements, and discussions.

* [@cheqd\_io](https://twitter.com/cheqd_io) on Twitter
* [@cheqd](https://t.me/cheqd) on Telegram \(with a separate [announcements-only channel](https://t.me/cheqd_announcements)\)
* [Medium](https://blog.cheqd.io/) blog
* [LinkedIn](http://cheqd.link/linkedin)

