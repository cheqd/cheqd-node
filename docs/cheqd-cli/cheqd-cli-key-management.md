# Using cheqd Cosmos CLI to manage keys

## Overview

[cheqd Cosmos CLI](README.md) can be used manage keys on a node.Keys are closely related to accounts and on-ledger authentication.

Account addresses are on a cheqd node are an encoded version of a public key. Each account is linked with at least one public-private key pair. Multi-sig accounts can have more than one key pair associated with them.

To submit a transaction on behalf of an account, it must be signed with an account's private key.

Cosmos supports [multiple keyring backends](https://docs.cosmos.network/master/run-node/keyring.html) for the storage and management of keys. Each node operator is free to use the key management method they prefer.

By default, the `cheqd-noded` binary is configured to use the `os` keyring backend, as it is a safe default compared to file-based key management methods.

For test networks or local networks, this can be overridden to the `test` keyring backend which is less secure and uses a file-based key storage mechanism where the keys are stored un-encrypted. To use the `test` keyring backend, append `--keyring-backend test` to each command that is related to key management or usage.

### Types of keys on a cheqd node

Each cheqd validator node has at least two keys.

#### Node key

* Default location is `$HOME/config/node_key.json`
* Used for peer-to-peer communication

#### Validator key

* Default location is `$HOME/config/priv_validator_key.json`
* Used to sign consensus messages

## Node-related commands in cheqd CLI

### Creating a key

When a new key is created, an **account address** and a **mnemonic backup phrase **will be printed. Keep mnemonic safe. This is the only way to restore access to the account if they keyring cannot be recovered.

#### Command

```bash
cheqd-noded keys add <alias>
```

### Restoring a key from backup mnemonic phrase

Allows restoring a key from a previously-created BIP39 **mnemonic phrase**.

#### Command

```bash
cheqd-noded keys add --recover <alias>
```

### Listing available keys on a node

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


### Example of working with test account


For example, for test purposes let's create a key with alias `operator`:


```bash
$ docker run -it --rm -u cheqd ghcr.io/cheqd/cheqd-node:0.4.0 keys add operator
Enter keyring passphrase:
Re-enter keyring passphrase:

- name: operator
  type: local
  address: cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0N73o8ke8bp/7c7PgsRjHGddjHvk0USHwq+RDzwwE0t"}'
  mnemonic: ""


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

crawl field same drill indoor olympic tank lamp range olive announce during pact idea fall canal sauce film attend response mammal bounce stable suffer
```

  

The main bullets here:

  

*   operator address: `cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh`

mnemonic phrase ( 24 words ):

*   `crawl field same drill indoor olympic tank lamp range olive announce during pact idea fall canal sauce film attend response mammal bounce stable suffer`

  

Having this mnemonic phrase the user is able to restore their keys whenever they want. For continue playing a user needs to run:

  

```bash
$ docker run -it --rm -u root --entrypoint bash ghcr.io/cheqd/cheqd-node:0.4.0
... apt install ca-certificates
... su cheqd

cheqd@8c3f88f653ab:~$ cheqd-noded keys add operator --recover --keyring-backend test
> Enter your bip39 mnemonic
crawl field same drill indoor olympic tank lamp range olive announce during pact idea fall canal sauce film attend response mammal bounce stable suffer

- name: operator
  type: local
  address: cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0N73o8ke8bp/7c7PgsRjHGddjHvk0USHwq+RDzwwE0t"}'
  mnemonic: ""

cheqd@8c3f88f653ab:~$ cheqd-noded keys list --keyring-backend test
- name: operator
  type: local
  address: cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0N73o8ke8bp/7c7PgsRjHGddjHvk0USHwq+RDzwwE0t"}'
  mnemonic: ""

cheqd@8c3f88f653ab:~$
```

  

As you can see, the recovered address is the same as was created before.

  

And after that all the commands from the tutorial above can be called.

  

P.S. the case with `docker` can be used only for demonstration purposes, cause after closing the container all the data will be lost.

For production purposes, maybe it would be great to have an image with Ubuntu 20.04 and operator's keys inside.