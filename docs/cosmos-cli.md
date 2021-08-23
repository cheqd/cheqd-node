# cheqd Command Line Interface (CLI)

There are two CLI tools for cheqd:

- Cosmos CLI: This is intended for node operators. Typically for node configuration, setup, and Cosmos keys.
- [Verifiable Data Registry (VDR) Tools](https://gitlab.com/evernym/verity/vdr-tools) CLI: This is intended for carrying out interactions related to decentralised identity / self-sovereign identity (SSI) functionality.

This document describes common workflows for cheqd Cosmos CLI.


## Managing keys

Keys are closely related to accounts and on-ledger authentication.

Account address is a properly encoded hash of public key. It means that each account is connected with at least one key. (MultiSig accounts are exceptions.)

To submit a transaction on behalf of an account, it must be signed with account's private key.

It's highly recommended to add `--keyring-backend os` to each command that is related to key management or usage. Cosmos supports [multiple keyring backends](https://docs.cosmos.network/v0.43/run-node/keyring.html), so each node operator is free to use the method they choose. `os` is a safe default to use.

__Creating a key__

```
cheqd-noded keys add <alias>
```

`Mnemonic phrase` and `account address` will be printed. Keep mnemonic safe. This is the only way to restore access to the account if they keyring cannot be recovered.

__Restoring a key from backup mnemonic key__

```
cheqd-noded keys add <alias> --recover <mnemonic>
```

__Listing available keys on a node__

```
cheqd-noded keys list
```

__Using a key for transaction signing__

Most transactions will require you to use `--from <key-alias>` param which is a name or address of private key with which to sign a transaction.

```
cheqd-noded tx <module> <tx-name> --from <key-alias>
```

## Querying ledger

Typical ledger query command looks like this:

```
cheqd-noded query <module> <query> <params> --node <url>
```

Example:

```
cheqd-noded query bank balances cosmos1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v --node http://nodes.testnet.cheqd.network:26657
```

Arguments:
- `--node` - IP address or URL of node to send the request to

## Submitting transactions

Typical transaction submit command looks like this:

```
cheqd-noded tx <module> <tx> <params> --node <url> --chain-id <chain> ---fee <fee>
```

Example:

```
cheqd-noded tx bank send alice cosmos10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000cheq --node http://localhost:26657 --chain-id cheqd ---fee 100000cheq
```

Extra arguments:
- `--node` - IP address or URL of node to send request to
- `--chain-id` - i.e. `cheqd-testnet`
- `--fees` - Max fee limit that is allowed for the transaction. 

Status code:

Pay attention at return status code. It should be 0 if a transaction is submitted successfully. Otherwise, an error message may be returned.

## Managing NYMs

[**NYM** is the term used by Hyperledger Indy](https://hyperledger-indy.readthedocs.io/projects/node/en/latest/transactions.html#nym) for Decentralized Identifiers (DIDs) that are created on ledger. A DID is typically the identifier that is associated with a specific organisation issuing/managing SSI credentials.

For the sake of explaining with similar concepts to current Hyperledger Indy implementations, on the `cheqd-testnet` these are still called NYMs.

Transactions to add a DID to the ledger are called NYM transactions.

Future releases of `cheqd-node` are likely to replace the NYM terminology with DID for better understanding.

__Creating a NYM:__

Command:

```
cheqd-noded tx cheqd create-nym <alias> <verkey> <did> <role>  --from <key-alias> --node <url> --chain-id <chain> ---fee <fee>
```

Example:

```
cheqd-noded tx cheqd create-nym "alias" "verkey" "did" "role"  --chain-id cheqd --from alice --node http://localhost:26657 --chain-id cheqd ---fee 100000cheq
```

ID of the created NYM will be returned.

__Querying a NYM by ID:__

Command:

```
cheqd-noded query cheqd show-nym <id>  --node <url>
```

Example:

```
cheqd-noded query cheqd show-nym 0 --node http://localhost:26657
```

__Listing on-chain NYMs:__

Command:

```
cheqd-noded query cheqd list-nym  --node <url>
```

Example:

```
cheqd-noded query cheqd list-nym --node http://localhost:26657
```

## Managing account balances

__Querying account balances:__

Command:

```
cheqd-noded query bank balances <address> --node <url>
```

Example:

```
cheqd-noded query bank balances cosmos1lxej42urme32ffqc3fjvz4ay8q5q9449f06t4v --node http://nodes.testnet.cheqd.network:26657
```

__Transferring tokens:__

Command:

```
cheqd-noded tx bank send <from> <to-address> <amount> --node <url> --chain-id <chain> ---fee <fee>
```

Params:
- `from` can be either key alias or address. If it's an address, corresponding key should be in keychain.

Example:

```
cheqd-noded tx bank send alice cosmos10dl985c76zanc8n9z6c88qnl9t2hmhl5rcg0jq 10000stake --node http://localhost:26657 --chain-id cheqd ---fee 100000cheq
```

## Managing node

cheqd node has two keys:
- Node key:
  - Default location is `$NODE_HOME/config/node_key.json`
  - Used for p2p communication
- Validator key:
  - Default location is `$NODE_HOME/config/priv_validator_key.json`
  - Used to sign consensus messages

__Running node__


```
cheqd-noded start
```

__Getting node address (node ID)__

Node ID or node address is a part of peer info. It's calculated from node's `pubKey` as `hex(address(nodePubKey))`. To get `node id` run the following command on the node's machine:

```
cheqd-noded tendermint show-node-id
```

__Getting validator address__

Validator address is a function of validator's public key. To get `bech32` encoded validator address run this command on node's machine:

```
cosmosvalcons1l43yqtdjcvyj65vnp29ly8u8yyau92q0ptzdp0
```

There are several ways to get hex encoded validator address:

1. Convert from bech32

    ```
    cheqd-noded keys parse <bech-32-encoded-address>
    ```

2. Query node using CLI:

    ```
    cheqd-noded tendermint show-address --node <node-prc-url>
    ```

    Look for `"ValidatorInfo":{"Address":"..."}`.

__Getting validator public key__

Validator public key is used in `create-validator` transactions. To get `bech32` encoded validator public key, run the following command on the node's machine:

```
cheqd-noded tendermint show-validator
```

__Sharing peer information__

Peer info is used to connect to peers when setting up a new node. It has the following format:

```
<node-id>@<node-url>
```

Example:

```
ba1689516f45be7f79c7450394144711e02e7341@3.13.19.41:26656
```

Using this information other participants will be able to join your node.
