# Using cheqd Cosmos CLI for token transactions

## Overview

A `cheqd-node` instance can be controlled and configured using the [cheqd Cosmos CLI](README.md).

This document contains the commands for reading and writing token transactions.

## Token-related transaction commands in cheqd CLI

### Querying ledger

#### Command

```bash
cheqd-noded query <module> <query> <params> --node <url>
```

#### Arguments

* `--node`: IP address or URL of node to send the request to

#### Example

```bash
$ cheqd-noded query bank balances 

cheqd1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v --node http://nodes.testnet.cheqd.network:26657
```

### Submitting transactions

#### Command

```bash
cheqd-noded tx <module> <tx> <params> --node <url> --chain-id <chain> --fees <fee>
```

#### Arguments

* `--node`: IP address or URL of node to send request to
* `--chain-id`: i.e. `cheqd-testnet-2`
* `--fees`: Maximum fee limit that is allowed for the transaction.

#### Status codes

Pay attention at return status code. It should be 0 if a transaction is submitted successfully. Otherwise, an error message may be returned.

#### Example

```bash
$ cheqd-noded tx bank send alice 

cheqd10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000ncheq --node http://localhost:26657 --chain-id cheqd --fees 50ncheq
```
