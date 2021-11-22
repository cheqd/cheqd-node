# Using cheqd Cosmos CLI to manage a node

## Overview

A `cheqd-node` instance can be controlled and configured using the [cheqd Cosmos CLI](README.md).

This document contains the commands for node operators that relate to node management, configuration, and status.

## Node-related commands in cheqd CLI

### Starting a node

```bash
cheqd-noded start
```

### Getting node ID / node address

Node ID or node address is a part of peer info. It's calculated from node's `pubKey` as `hex(address(nodePubKey))`. To get `node id` run the following command on the node's machine:

```bash
cheqd-noded tendermint show-node-id
```

### Get validator address

Validator address is a function of validator's public key. To get `bech32` encoded validator address run this command on node's machine:

```bash
$ cheqd-noded tendermint show-address
cheqdvalcons1sg4azh7qwk6akm0eadkgvgq2kegtzksr09a685
```

There are several ways to get hex-encoded validator address:

1. Convert from bech32

   ```bash
   cheqd-noded keys parse <bech-32-encoded-address>
   ```

2. Query node using CLI:

   ```bash
   cheqd-noded status
   ```

   Look for `"ValidatorInfo":{"Address":"..."}`

### Getting validator public key

Validator public key is used in `create-validator` transactions. To get `bech32` encoded validator public key, run the following command on the node's machine:

```bash
$ cheqd-noded tendermint show-validator

{"@type":"/cosmos.crypto.ed25519.PubKey","key":"y8v/nsf+VFCnJ7c9ZM/C4tUMnWKHhU+K+B82B+5vUZg="}
```

### Sharing peer information

Peer info is used to connect to peers when setting up a new node. It has the following format:

```bash
<node-id>@<node-url>
```

Example:

```bash
ba1689516f45be7f79c7450394144711e02e7341@3.13.19.41:26656
```

Using this information other participants will be able to join your node.
