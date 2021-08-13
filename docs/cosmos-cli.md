# cheqd cosmos cli

There are two cli tools for cheqd:

- Cosmos SDK based: for infrastructure management;
- VDR tools based: for identity management.

This document describes common workflows for cheqd cosmos cli.


## Managing keys

Keys are closely related to accounts and on ledger authentication.

Account address is a properly encoded hash of public key. It means that each account is connected with one key (multisig accounts are exceptions).

To submit a transaction on behalf of an account it must be signed with account's private key.

It's highly recommended to add `--keyring-backend os` to each command that is related to key management or usage.

__Key creation:__

```
cheqd-noded keys add <alias>
```

`Mnemonic phrase` and `account address` will be printed. Kee mnemonic safe. This is the only wa to restore access to the account.

__Key restoring:__

```
cheqd-noded keys add <alias> --recover <mnemonic>
```

__Keys listing:__

```
cheqd-noded keys list
```

__Using a key for transaction signing:__

Most transactions will require you to use `--from <key-alias>` param which is a name or address of private key with which to sign tx.

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

Extra arguments:
- `--node` - ip address or url of node to connect.

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
- `--node` - ip address or url of node to connect;
- `--chain-id` - i.e. `cheqd-testnet` or `cheqd-mainnet`;
- `--fees` - max fee to pay along with transaction.

Status code:

Pay attention at return status code. It should be 0 if a transaction is submitted successfully. Otherwise error message is returned.

## Managing NYMs

__Creating a NUM:__

Command:

```
cheqd-noded tx cheqd create-nym <alias> <verkey> <did> <role>  --from <key-alias> --node <url> --chain-id <chain> ---fee <fee>
```

Example:

```
cheqd-noded tx cheqd create-nym "alias" "verkey" "did" "role"  --chain-id cheqd --from alice --node http://localhost:26657 --chain-id cheqd ---fee 100000cheq
```

Id of the created NYM will be returned.

__Querying a NYM by id:__

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
