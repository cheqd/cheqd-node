# Using cheqd Cosmos CLI to manage accounts

## Overview

A `cheqd-node` instance can be controlled and configured using the [cheqd Cosmos CLI](README.md).

This document contains the commands for account management.

## Account-related commands in cheqd CLI

### Querying account balances

#### Command

```bash
cheqd-noded query bank balances <address>
```

#### Example

```bash
cheqd-noded query bank balances cheqd1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v
```

### Transferring tokens

#### Command

```bash
cheqd-noded tx bank send <from> <to-address> <amount> --node <url> --chain-id <chain> --fees <fee>
```

#### Arguments

* `from` can be either key alias or address. If it's an address, corresponding key should be in keychain.

#### Example

```bash
$ cheqd-noded tx bank send alice 
cheqd10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000ncheq --node http://localhost:26657 --chain-id cheqd --fees 50ncheq
```
