# Overview
This document includes just steps for writing a DID to the pool, without additional info

## Steps for writing a DID to the mainnet
1. Create a wallet:
   `wallet create cheqd-wallet key`
2. Open a wallet:
   `wallet open cheqd-wallet key`
3. Create a pool connection:
   `cheqd-pool add alias=cheqd-pool-mainnet rpc_address=https://rpc.cheqd.net chain_id=cheqd-mainnet-1`
4. Open a connection to the mainnet:
   `cheqd-pool open alias=cheqd-pool-mainnet`
5. Generate a DID:
   - `did new`
   - copy generated DID
   - `did new did=did:cheqd:mainnet:<copied-DID>`
6. Generate/restore cosmos signature:
   - `cheqd-keys add alias=cheqd-mainnet-key mnemonic`
   - paste your mnemonic after suggestion `Enter value for mnemonic:`
7. Send previously created DID:
   `cheqd-ledger create-did did=<DID_from_step_5> max_coin=<value> max_gas=<value> denom=ncheq memo=memo key_alias=cheqd-mainnet-key`
   You need only setup `max_coin` and `max_gas` parameters, according to your financial politic.