# Overview

This document includes just steps for writing a DID to the pool, without additional info.

## Pre-requirements for environment

All the commands are prepared for using inside Ubuntu 20.04. It can be a Virtual machine, for example VirtualBox or docker.
Maybe, the fastest way - to choose docker, like:
`docker run -it --rm ubuntu:20.04 bash`
and after that follow the installation flow. 

For installation process of `vdrtools-cli` [this instruction](install.md) can help

## Steps for writing a DID to the testnet

1. Create a wallet:
   `wallet create cheqd-wallet key`
2. Open a wallet:
   `wallet open cheqd-wallet key`
3. Create a pool connection:
   `cheqd-pool add alias=cheqd-pool-testnet rpc_address=https://rpc.testnet.cheqd.network/ chain_id=cheqd-testnet-4`
4. Open a connection to the testnet:
   `cheqd-pool open alias=cheqd-pool-testnet`
5. Generate a DID:
   - `did new`
   - copy generated DID
   - `did new did=did:cheqd:testnet:<copied-DID>`
6. Generate/restore cosmos signature:
   - `cheqd-keys add alias=cheqd-testnet-key mnemonic`
   - paste your mnemonic after suggestion `Enter value for mnemonic:`
7. Send previously created DID:
   `cheqd-ledger create-did did=<DID_from_step_5> max_coin=<value> max_gas=<value> denom=ncheq memo=memo key_alias=cheqd-testnet-key`
   You need only setup `max_coin` and `max_gas` parameters, according to your financial politic.