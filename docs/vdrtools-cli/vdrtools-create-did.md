# How to create DID using VDR tools
In general, we need to create to main things, wallet and did inside it.
Also, all the next commands shoud be typed inside the CLI. For make it possible, just run in terminal:
```
vdrtools-cli
```
## Create wallet
For wallet creation we just need:
```
indy> wallet create <name-of-wallet> key
```
For example, sesson for creating a wallet with name `cheqd-wallet` looks, like:
```
indy> wallet create cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been created
indy>
```
where `key` - it's a password for your wallet

## Open wallet
For all operations with items in the wallet we need to open it before. For make it possible, please call:
```
indy> wallet open <name-of wallet> key
```
For example, for our previously created wallet we need to call:
```
indy> wallet open cheqd-wallet key
Enter value for key:

Wallet "cheqd-wallet" has been opened
```

## Create key for cosmos signature
One of signatures that required for validation on the cheqd-node side is cosmos's signature.
For adding a signing key we need to create/restore from mnemonic by using:
```
indy> cheqd-keys add alias=<name-of-key>
```
or in case of restoring from mnemonic it will be:
```
indy> cheqd-keys add alias=<name-of-key> mnemonic
```
and there will be a suggestion to insert the full mnemonic phrase in silent mode (it will not be shown).
For example, for creating a new random key:
```
cheqd-wallet:indy> cheqd-keys add alias=cheqd-key
Random key has been added.
Account_id is: "cheqd1dguvqgrlkn5sts2zfaqmlk4t6gdhv6e8puvkge"
**Important** write this mnemonic phrase in a safe place.
Mnemonic phrase is:
electric marine palace chaos open review friend left convince pupil spoon cigar brain mass cake bronze potato suspect answer pig common alert ice choose
```

## Create and connect to the pool
The next step which is needed for establishing connection to the server it's a command, name `cheqd-pool`.

### Create connection to the network
In general it's some kind of configuration, which can be reusable in the future. The main command here:
```
cheqd-pool add alias=<alias-value> rpc_address=<rpc_address-value> chain_id=<chain_id-value>
```
For example, for local network it can be:
```
indy> cheqd-pool add alias=cheqd-pool rpc_address=http://127.0.0.1:26657 chain_id=cheqd
Pool "cheqd-pool" has been created "{"alias":"cheqd-pool","rpc_address":"http://127.0.0.1:26657","chain_id":"cheqd"}"
```

### Connect to the pool
```
cheqd-pool open alias=<alias-value>
```
For connecting to the previously created instance it would be:
```
indy> cheqd-pool open alias=cheqd-pool
Pool "cheqd-pool" has been opened
```

