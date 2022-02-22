# Tutorial for creating a DID with VDR tools

## Overview

The purpose of this document is to outline how someone can create a DID on the cheqd network using [VDR tools.](https://gitlab.com/evernym/verity/vdr-tools)

This tutorial uses VDR Tools CLI to send a DID. The official page:

[https://gitlab.com/evernym/verity/vdr-tools](https://gitlab.com/evernym/verity/vdr-tools)


## Pre-requisites

In order to create a DID using the instructions outlined in this tutorial, you must be using Ubuntu 20.04 terminal. You'll find all the information required to setup your Ubuntu 20.04 terminal at the end of this tutorial.

Please ensure you are running the correct version of testnet. You can check which is the current version of testnet [here](https://rpc.testnet.cheqd.network/abci_info?).

### Ensure you have the correct OS environment

If you don't currently have Ubuntu 20.04 installed on your machine you can use VirtualBox or [Docker](#docker-setup)

### Install VDR Tools on Ubuntu 20.04

After the preparation of the Ubuntu 20.04 terminal please make sure you have installed VDR tools.

In general this process is described on the main page of VDR tools [repository](https://gitlab.com/evernym/verity/vdr-tools#installing) but to be short let's make the next steps inside Ubuntu 20.04:

#### 1\. Install additional packages:

```bash
$ apt install curl libsodium23 libzmq5 libncursesw5-dev -y
```

#### 2\. Download and install `libvdrtools`:

```bash
$ curl https://gitlab.com/evernym/verity/vdr-tools/-/package_files/27311917/download --output libvdrtools_0.8.4-focal_amd64.deb && dpkg -i libvdrtools_0.8.4-focal_amd64.deb
```

#### 3\. Download and install `vdrtools-cli`:

```bash
$ curl https://gitlab.com/evernym/verity/vdr-tools/-/package_files/27311922/download --output vdrtools-cli_0.8.4-focal_amd64.deb && dpkg -i vdrtools-cli_0.8.4-focal_amd64.deb
```

Once this is installed we're now ready to write a DID to the cheqd network.

## Writing a DID to the cheqd network

All the command described below are supposed to be used from VDR Tools CLI.

For running VDR Tools terminal just type:

```bash
$ vdrtools-cli
```

The next steps are about how to send a DID.

#### 1\. Create a wallet

The command:

```bash
wallet create cheqd-wallet key 
```

Example:

```bash
indy> wallet create cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been created
```

where `Enter value for key` is a suggestion to create a secure password.

#### 2\. Open your wallet: 

Next use the following command to open the wallet created in the previous step

```bash
wallet open cheqd-wallet key
```

Within this command you will need to enter the password you have just created.

Example:

```bash
indy> wallet open cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been opened
```

#### 3\. Create a pool connection: 

Now we're ready to create a pool connection. This is a configuration to allow to connection to the of your wallet to the pool.

Use this command:

```bash
cheqd-pool add alias=cheqd-pool-testnet rpc_address=https://rpc.testnet.cheqd.network/ chain_id=cheqd-testnet-4
```

Example:

```bash
cheqd-wallet:indy> cheqd-pool add alias=cheqd-pool-testnet rpc_address=https://rpc.testnet.cheqd.network/ chain_id=cheqd-testnet-4
Pool "cheqd-pool-testnet" has been created "{"alias":"cheqd-pool-testnet","rpc_address":"https://rpc.testnet.cheqd.network/","chain_id":"cheqd-testnet-4"}"
```

#### 4\. Open a connection to the testnet: 

Next we need to setup a direct connection between the users machine and the pool on testnet.

Use this command:

```bash
cheqd-pool open alias=cheqd-pool-testnet
```

Example:

```bash
cheqd-wallet:indy> cheqd-pool open alias=cheqd-pool-testnet
Pool "cheqd-pool-testnet" has been opened
```

#### 5\. Generate a DID:

Now the connection has been made, we're ready to create a DID. In this case we'll generate a DID first locally. This is because the current functionality of VDR tools only allows for the creation of a unique identifier as opposed to a DID.

To do this we'll use need to use a 3 step process using 2 different commands

1.  First use `did new`
2.  Then copy the newly generated DID
3.  Finally, paste your new DID into the allocated space in the templated command `did new did=did:cheqd:testnet:<copied-DID>`

Example:

```bash
cheqd_pool(cheqd-pool-testnet):cheqd-wallet:indy> did new
Did "X9cNs1g3nMvQ77j7E8rPRr" has been created with "~4zHub2htDfEuV7wCsNcEgg" verkey

cheqd_pool(cheqd-pool-testnet):cheqd-wallet:indy> did new did=did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr
Did "did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr" has been created with "7db7n4CNBbYFRcYEjoacMzSbVEqF2ToSovtvzBoovH5B" verkey
```

#### 6\. Generate/restore cosmos signature:

Now we'll need to restore the Cosmos keys from the mnemonic. The mnemonic required here was provided during node operator account creation.

2-steps command:

*   `cheqd-keys add alias=cheqd-testnet-key mnemonic`
*   paste your mnemonic after suggestion `Enter value for mnemonic:`

Example:

```bash
cheqd_pool(cheqd-pool-testnet):cheqd-wallet:indy> cheqd-keys add alias=cheqd-testnet-key mnemonic
Enter value for mnemonic:

Key for account: cheqd1waszljnwadcuynjg6wc9rhhw8ega20nc97muc6
has been restored from mnemonic
```

#### 7\. Send previously created DID: 

We can now use the VDR tools to send our newly created DID using the following command:

```bash
cheqd-ledger create-did did=<DID_from_step_5> max_coin=<value> max_gas=<value> denom=ncheq memo=memo key_alias=cheqd-testnet-key
```

Note: You only need to setup `max_gas` and `max_coin` parameters, according to your financial policy. Defaults are `max_gas=200000 max_coin=5000000`

Example:

```bash
cheqd_pool(cheqd-pool-testnet):cheqd-wallet:indy> cheqd-ledger create-did did=did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr denom=ncheq memo=memo key_alias=cheqd-testnet-key max_gas=200000 max_coin=5000000
Abci-info request result "{"response":{"app_version":"1","data":"cheqd-node","last_block_app_hash":[74,43,67,81,106,98,53,116,108,84,99,43,107,97,52,119,43,49,103,83,82,73,89,69,114,85,55,105,77,48,70,53,109,79,120,69,66,69,104,119,109,109,107,61],"last_block_height":"495420","version":"0.4.0"}}"
Response from ledger: "{\"id\":\"did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr\"}"
```

#### 8\. Check the previously sent DID

Finally we can now check that the DID was successfully written to the ledger.

The command:

```bash
cheqd-ledger get-did did=<DID_from_step_5>
```

Example:

```bash
cheqd_pool(cheqd-pool-testnet):cheqd-wallet:indy> cheqd-ledger get-did did=did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr
DID info: {"did":{"id":"did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr","controller":["did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr"],"verification_method":[{"id":"did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr#verkey","type":"Ed25519VerificationKey2020","controller":"did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr","public_key_multibase":"z7db7n4CNBbYFRcYEjoacMzSbVEqF2ToSovtvzBoovH5B"}],"authentication":["did:cheqd:testnet:X9cNs1g3nMvQ77j7E8rPRr#verkey"]},"metadata":{"created":"2022-02-21 14:06:50.162709363 +0000 UTC","updated":"2022-02-21 14:06:50.162709363 +0000 UTC","deactivated":false,"version_id":"Gf0tfbMUpHHR+z8j7KdXLInlAi4T51F6qgzxUyi9h+M="}}
```

## Docker setup

!!! It's very important. Please, take it into account, that using docker it's just a suggestion to go through the tutorial as the fastest way. Or you have to make your data saved after docker container stopped.

### 1\. Start a container

```bash
$ docker run -it --rm -u root ubuntu:20.04 bash
```

### 2\. Update the package storage

```bash
# apt update
```
