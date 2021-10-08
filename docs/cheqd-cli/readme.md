# cheqd Command Line Interface (CLI) guide

## Overview

There are two command line interface (CLI) tools for interacting with a running `cheqd-node` instance:

* **cheqd Cosmos CLI**: This is intended for node operators. Typically for node configuration, setup, and Cosmos keys.
* **[Verifiable Data Registry \(VDR\) Tools](https://gitlab.com/evernym/verity/vdr-tools) CLI**: This is intended for carrying out interactions related to decentralised identity / self-sovereign identity (SSI) functionality.

This document is focussed on providing guidance on ow to use the **cheqd Cosmos CLI**.

## Querying ledger

Typical ledger `query` command looks like this:

```bash
cheqd-noded query <module> <query> <params> --node <url>
```

### Arguments for `query` command

* `--node`: IP address or URL of node to send the request to

### Example of `query` command

```bash
$ cheqd-noded query bank balances cheqd1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v --node http://nodes.testnet.cheqd.network:26657
```

## Submitting transactions

Typical transaction submit (`tx`) command looks like this:

```bash
cheqd-noded tx <module> <tx> <params> --node <url> --chain-id <chain> --fees <fee>
```

### Arguments for `tx` command

* `--node`: IP address or URL of node to send request to
* `--chain-id`: i.e. `cheqd-testnet-2`
* `--fees`: Maximum fee limit that is allowed for the transaction.

### Status codes for `tx` command

Pay attention at return status code. It should be 0 if a transaction is submitted successfully. Otherwise, an error message may be returned.

### Example of `tx` command

```bash
$ cheqd-noded tx bank send alice cosmos10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000cheq --node http://localhost:26657 --chain-id cheqd --fees 100000cheq
```

## Managing account balances

### Querying account balances

Command:

```bash
cheqd-noded query bank balances <address> --node <url>
```

### Example `query` command for querying account balance

```bash
$ cheqd-noded query bank balances cheqd1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v --node http://nodes.testnet.cheqd.network:26657
```

**Transferring tokens:**

Command:

```text
cheqd-noded tx bank send <from> <to-address> <amount> --node <url> --chain-id <chain> --fees <fee>
```

Params:

* `from` can be either key alias or address. If it's an address, corresponding key should be in keychain.

Example:

```text
cheqd-noded tx bank send alice cosmos10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000stake --node http://localhost:26657 --chain-id cheqd --fees 100000cheq
```
