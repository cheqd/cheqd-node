# Using cheqd Cosmos CLI to manage keys

## Overview

[cheqd Cosmos CLI](readme.md) can be used manage keys on a node.Keys are closely related to accounts and on-ledger authentication.

Account addresses are on a cheqd node are an encoded version of a public key. Each account is linked with at least one public-private key pair. Multi-sig accounts can have more than one key pair associated with them.

To submit a transaction on behalf of an account, it must be signed with an account's private key.

Cosmos supports [multiple keyring backends](https://docs.cosmos.network/master/run-node/keyring.html) for the storage and management of keys. Each node operator is free to use the key management method they prefer.

Our recommended method is to use the `os` keyring backend, as it is a safe default compared to file-based key management methods.

To use the `os` keyring backend, append `--keyring-backend os` to each command that is related to key management or usage.

### Types of keys on a cheqd node

Each cheqd validator node has at least two keys.

#### Node key

* Default location is `$NODE_HOME/config/node_key.json`
* Used for peer-to-peer communication

#### Validator key

* Default location is `$NODE_HOME/config/priv_validator_key.json`
* Used to sign consensus messages

## Node-related commands in cheqd CLI

### Creating a key

`Mnemonic phrase` and `account address` will be printed. Keep mnemonic safe. This is the only way to restore access to the account if they keyring cannot be recovered.

#### Command

```bash
cheqd-noded keys add <alias>
```

### Restoring a key from backup mnemonic phrase

Allows restoring a key from a previously-created bip39 `mnemonic phrase`.

#### Command

```bash
cheqd-noded keys add --recover <alias>
```

### Listing available keys on a node\*\*

#### Command

```bash
cheqd-noded keys list
```

### Using a key for transaction signing

Most transactions will require you to use `--from <key-alias>` param which is a name or address of private key with which to sign a transaction.

#### Command

```bash
cheqd-noded tx <module> <tx-name> --from <key-alias>
```
