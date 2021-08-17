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

## Discuss this project, get help and support
The [cheqd Community Slack is open for anyone to join](http://cheqd.link/join-cheqd-slack) and is our discussion forum for the open-source community, software developers, and node operators. Please reach out to us for real-time discussions and help.

You can also follow the cheqd on other social channels:
- [Twitter](https://twitter.com/cheqd_io)
- [Telegram](https://t.me/cheqd) (with a separate [announcements-only channel](https://t.me/cheqd_announcements))
- [Medium](https://blog.cheqd.io/)
- [LinkedIn](http://cheqd.link/linkedin)

## Learn more about the Cosmos ecosystem

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Cosmos Discord chat](https://discord.gg/W8trcGV)
