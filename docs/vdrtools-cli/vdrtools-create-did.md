# How to create DID using VDR tools

In general, we need to create to main things, wallet and did inside it.
Also, all the next commands shoud be typed inside the CLI. For make it possible, just run in terminal:

```bash
vdrtools-cli
```

## Create wallet

For wallet creation we just need:

```bash
indy> wallet create <name-of-wallet> key
```

For example, sesson for creating a wallet with name `cheqd-wallet` looks, like:

```text
indy> wallet create cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been created
indy>
```

where `key` - it's a password for your wallet

## Open wallet

For all operations with items in the wallet we need to open it before. For make it possible, please call:

```bash
indy> wallet open <name-of wallet> key
```

For example, for our previously created wallet we need to call:

```text
indy> wallet open cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been opened
```

## Create key for cosmos signature

One of signatures that required for validation on the cheqd-node side is cosmos's signature.
For adding a signing key we need to create/restore from mnemonic by using:

```bash
indy> cheqd-keys add alias=<name-of-key>
```

or in case of restoring from mnemonic it will be:

```bash
indy> cheqd-keys add alias=<name-of-key> mnemonic
```

and there will be a suggestion to insert the full mnemonic phrase in silent mode (it will not be shown).
For example, for creating a new random key:

```text
cheqd-wallet:indy> cheqd-keys add alias=cheqd-key
Random key has been added.
Account_id is: "cheqd1dguvqgrlkn5sts2zfaqmlk4t6gdhv6e8puvkge"
**Important** write this mnemonic phrase in a safe place.
Mnemonic phrase is:
electric marine palace chaos open review friend left convince pupil spoon cigar brain mass cake bronze potato suspect answer pig common alert ice choose
```

## Generate DID

For generating simple DID, like just an identificator can be used the next commmand:

```bash
indy> did new
```

As example, the result can be:

```text
cheqd_pool(cheqd-pool):cheqd-wallet:indy> did new
Did "C9mR4KH6Mb7FWsCjsAfnVo" has been created with "~T3TuBAkJPgFyre8bVWuRjF" verkey
```

If full DID with `testnet` or `mainnet` service is needed the next steps are useful:

1. `did new`
2. `did new did=<prefix>:<previously_created_DID>`
   
In this case after the step 2 we can use this DID for sending to the cheqd network.
For example:

```text
cheqd_pool(cheqd-pool):cheqd-wallet:indy> did new
Did "C9mR4KH6Mb7FWsCjsAfnVo" has been created with "~T3TuBAkJPgFyre8bVWuRjF" verkey
cheqd_pool(cheqd-pool):cheqd-wallet:indy> did new did=did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo
Did "did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo" has been created with "FrVpdaHkumCBYW91a1j5gLQwDSr65S72KErTqhoXcL65" verkey
cheqd_pool(cheqd-pool):cheqd-wallet:indy>
```

where the first created DID was `C9mR4KH6Mb7FWsCjsAfnVo`.

## Create and connect to the pool

The next step which is needed for establishing connection to the server it's a command, name `cheqd-pool`.

### Create connection to the network

In general it's some kind of configuration, which can be reusable in the future. The main command here:

```bash
cheqd-pool add alias=<alias-value> rpc_address=<rpc_address-value> chain_id=<chain_id-value>
```

For example, for local network it can be:

```text
indy> cheqd-pool add alias=cheqd-pool rpc_address=http://127.0.0.1:26657 chain_id=cheqd
Pool "cheqd-pool" has been created "{"alias":"cheqd-pool","rpc_address":"http://127.0.0.1:26657","chain_id":"cheqd"}"
```

### Connect to the pool

```bash
cheqd-pool open alias=<alias-value>
```

For connecting to the previously created instance it would be:

```text
indy> cheqd-pool open alias=cheqd-pool
Pool "cheqd-pool" has been opened
```

## Send a DID to the pool

After the whole previos steps we are able to send a DID to the pool, which we connected previously.
For sending DID the next command can be used:

```bash
cheqd-ledger create-did did=<did-value> key_alias=<key_alias-value> max_coin=<max_coin-value> max_gas=<max_gas-value> denom=<denom-value> [memo=<memo-value>] [simulate=<simulate-value>]
```

where:

- `did-value` - previously created DID. In our examples it was `did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo`
- `key-alias` - previously created/restored key. In our example it was `cheqd-key`
- `denom`     - for cheqd pool the default value is `ncheq`
- `max_coin` and `max_gas` - financial parameters for the pool
  
For example:

```text
cheqd_pool(cheqd-pool):cheqd-wallet:indy> cheqd-ledger create-did did=did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo max_coin=250000000 max_gas=10000000 denom=ncheq memo=memo key_alias=cheqd-key
Abci-info request result "{"response":{"app_version":"1","data":"cheqd-node","last_block_app_hash":[121,54,48,74,111,116,97,74,68,57,72,102,55,102,65,112,49,78,105,50,98,88,89,48,122,109,48,57,74,111,99,66,77,65,50,65,108,107,99,101,116,100,65,61],"last_block_height":"2100","version":"0.4.0"}}"
Response from ledger: "{\"id\":\"did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo\"}"
```

## Get DID from the pool

For checking that DID was written successfully `get-did` subcommand can be used:

```bash
cheqd-ledger get-did did=<did-value>
```

As example for previously created DID:

```text
cheqd_pool(cheqd-pool):cheqd-wallet:indy> cheqd-ledger get-did did=did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo

DID info: {"did":{"id":"did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo","controller":["did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo"],"verification_method":[{"id":"did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo#verkey","type":"Ed25519VerificationKey2020","controller":"did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo","public_key_multibase":"zFrVpdaHkumCBYW91a1j5gLQwDSr65S72KErTqhoXcL65"}],"authentication":["did:cheqd:testnet:C9mR4KH6Mb7FWsCjsAfnVo#verkey"]},"metadata":{"created":"2022-02-16 09:28:53.0549393 +0000 UTC","updated":"2022-02-16 09:28:53.0549393 +0000 UTC","deactivated":false,"version_id":"mRg6zKRgu+EXksfulQsay4icEpCNaMmgUQfxIyuXoCc="}}
```
