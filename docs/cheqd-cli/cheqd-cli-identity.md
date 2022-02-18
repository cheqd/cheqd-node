# Using cheqd Cosmos CLI to manage DID Documents

## Overview

[cheqd Cosmos CLI](README.md) can be used for creating Decentralised Identifiers (DIDs) and DID Documents, as well as updating existing DIDs based on the [`did:cheqd` DID method](../../architecture/adr-list/adr-002-cheqd-did-method.md).

## Crypto commands in cheqd CLI - EXPERIMENTAL FUNCTIONALITY - Do not use in production!


#### Command for generating public and private part of verification key

```bash
cheqd-noded debug ed25519 random >> keys.txt
```

P.S. it's very important to keep your private verification key in a safe place. 
The simpliest suggestion here is to redirect output to file as shown in bash command above.

#### Response example

```text
{"pub_key_base_64":"MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=","pub_key_multibase_58":"",\
"priv_key_base_64":"FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q=="}
```

#### Convert base64 to multibase58
```bash
cheqd-noded debug encoding base64-multibase58 <public_key_in_base64_representaion>
```

#### Response example
```text
~ cheqd-noded debug encoding base64-multibase58 MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=
z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE
```


## Identity-related commands in cheqd CLI

**Note**: The `--chain-id` and `--node` flags are optional and only required if using the CLI from a machine that is not running an active cheqd node on mainnet/testnet.

### Create a new DID

#### Command

```bash
cheqd-noded tx cheqd create-did <did-document-json> <did-verification-method-id> --ver-key <did-private-key> --from <wallet-name> --gas auto --gas-adjustment 1.2 --gas-prices 25ncheq --chain-id <chain-id> --node <node-rpc-endpoint>
```

#### Sample DID Document JSON

```json
{
  "id": "did:cheqd:<namespace>:<unique-id>",
  "verification_method": [
    {
      "id": "did:cheqd:<namespace>:<unique-id>#<key-alias>",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:cheqd:<namespace>:<unique-id>",
      "public_key_multibase": "<verification-public-key-multibase>"
    }
  ],
  "authentication": [
    "did:cheqd:<namespace>:<unique-id>#<auth-key-alias>"
  ]
}
```

#### Arguments

* `id`: A unique identifier for format `did:cheqd`, conforming to the [cheqd DID method specification](../../architecture/adr-list/adr-002-cheqd-did-method.md).

* `id` feild for `verification_method` section: key identifier in the following format: `did:cheqd:<namespace>:<unique-id>#<key-alias>`. Exapmle: `did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1`.
* `--ver-key`: Base64 encoded ed25519 private key to sign identity message with. A pair for the key from DID Document. \
    Use for testing purposes only because the key will be stored in shell history!
* `--from`: Cosmos account key which will pay fees for the transaction ordering.
* `--node`: IP address or URL of node to send request to
* `--chain-id`: i.e. `cheqd-testnet-4`
* `--fees`: Maximum fee limit that is allowed for the transaction.


#### Example

```bash
cheqd-noded tx cheqd create-did '{"id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                   "verification_method": [{\
                                     "id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1",\
                                     "type": "Ed25519VerificationKey2020",\
                                     "controller": "id:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                     "public_key_multibase": "z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE"\
                                   }],\
                                   "authentication": [\
                                     "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"\
                                   ]\
                                 }' "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1" \
  --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 5000000ncheq
```

BTW, it can be more useful to prepare json file using your favorite editor and after that run the command as:

```bash
cheqd-noded tx cheqd create-did "$(cat json.txt)" --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 5000000ncheq
```

where `json.txt` - file with DID-Doc in json fromat.

### Querying a DID

Allows fetching the DID Document associated with an existing DID on cheqd networks.

#### Command

```bash
cheqd-noded query cheqd did <id> --chain-id <chain-id> --node <node-rpc-endpoint>
```

#### Example

```bash
cheqd-noded query cheqd did did:cheqd:testnet:zJ5EDiiiKWDyo79n --chain-id cheqd-testnet-4 --node http://rpc.testnet.cheqd.network:26657
```


### Update DID Document

#### Arguments

* `DIDDoc_in_JSON`: A string with DID Document in Json format.  `id` is not changeable field and mast be used from creation transaction.
    Base example:
```text
{
  "id": "did:cheqd:<namespace>:<unique-id>",
  "verification_method": [
    {
      "id": "did:cheqd:<namespace>:<unique-id>#<key-alias>",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:cheqd:<namespace>:<unique-id>",
      "public_key_multibase": "<verification-public-key-multibase>"
    }
  ],
  "authentication": [
    "did:cheqd:<namespace>:<unique-id>#<auth-key-alias>"
  ]
}
```

* `id` feild for `verification_method` section: key identifier in the following format: `did:cheqd:<namespace>:<unique-id>#<key-alias>`. Exapmle: `did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1`.
* `--ver-key`: Base64 encoded ed25519 private key to sign identity message with. A pair for the key which was used for the DID creation. \
    Use for testing purposes only because the key will be stored in shell history!
* `--from`: Cosmos account key which will pay fees for the transaction ordering.
* `--node`: IP address or URL of node to send request to
* `--chain-id`: i.e. `cheqd-testnet-4`
* `--fees`: Maximum fee limit that is allowed for the transaction.

#### Command

```bash
cheqd-noded tx cheqd create-did <DIDDoc_in_JSON> <did_verification_method_id> --ver-key <identity_private_key_BASE_64> \
  --from <cosmos_account> --node <url> --chain-id <chain> --fees <fee>
```

#### Example

```bash
cheqd-noded tx cheqd create-did '{"id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                   "verification_method": [{\
                                     "id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1",\
                                     "type": "Ed25519VerificationKey2020",\
                                     "controller": "id:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                     "public_key_multibase": "zCeJfYbiFoUcENEjuxnU9ez6VBZjxavTjSZtHP6y226fp"\
                                   }],\
                                   "authentication": [\
                                     "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"\
                                   ]\
                                 }' "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1" \
  --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 5000000ncheq
```

BTW, it can be more useful to prepare json file using your favorite editor and after that run the command as:

```bash
cheqd-noded tx cheqd create-did "$(cat json.txt)" --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 5000000ncheq
```
Let's the result will be like:
```text
{"pub_key_base_64":"MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=","pub_key_multibase_58":"",\
"priv_key_base_64":"FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q=="}
```
2. Get multibase58 string for inserting it into the DID-doc (`public_key_multibase` field in `verification_method` section)
   
```bash
cheqd-noded debug encoding base64-multibase58 MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=
```

Response will be:

```text
z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE
```

3. We need to insert multibase58 string from step 2 as `public_key_multibase` field in `verification_method` section
   and create a unic identifier for `verification_method`. It should be `<DID>#<some_unic_String>`. Let it be `did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1`.
After those preparations, the base DID-doc will look like:

```text
{
  "id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n",
  "verification_method": [
    {
      "id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1",
      "type": "Ed25519VerificationKey2020",
      "controller": "id:cheqd:testnet:zJ5EDiiiKWDyo79n",
      "public_key_multibase": "z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE"
    }
  ],
  "authentication": [
    "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"
  ]
}
```

And the full command for sending it to the pool will be:

```bash
cheqd-noded tx cheqd create-did '{"id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                   "verification_method": [{\
                                     "id": "did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1",\
                                     "type": "Ed25519VerificationKey2020",\
                                     "controller": "id:cheqd:testnet:zJ5EDiiiKWDyo79n",\
                                     "public_key_multibase": "z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE"\
                                   }],\
                                   "authentication": [\
                                     "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"\
                                   ]\
                                 }' "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1" \
  --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from one_of_our_test_account --node https://rpc.testnet.cheqd.network:443 --chain-id cheqd-testnet-4 --fees 50ncheq
```

4. For checking, that DID was successfully written the next query can be used:

```bash
cheqd-noded query cheqd did "did:cheqd:testnet:zJ5EDiiiKWDyo79n" --output json
```

## Requirements from OS side

Our target OS system is Ubuntu 20.04.
In this case, for running demo flow we can use variants: Virtualbox or docker. 
For example, let it be a docker image, cause it's the most fastest way to start playing.
The next command can help:

```bash
docker run -it --rm -u cheqd --entrypoint bash ghcr.io/cheqd/cheqd-node:0.4.0
```

The next step is to restore operator's keys:

```bash
cheqd-noded keys add <cheqd-operator-name> --recover --keyring-backend test
```

where, `cheqd-operator-name` it's name of alias for storing your keys locally, whatever you want.

### Exampe of working with test account.

For example, for test purposes let's create a key with alias `operator`:

```text
~ docker run -it --rm -u cheqd ghcr.io/cheqd/cheqd-node:0.4.0 keys add operator --keyring-backend test

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

- operator address: `cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh`
- mnemonic phrase ( 24 words ):
`crawl field same drill indoor olympic tank lamp range olive announce during pact idea fall canal sauce film attend response mammal bounce stable suffer`

Having this mnemonic phrase the user is able to restore their keys whenever they want. For continue playing a user needs to run:

```text
~ docker run -it --rm -u cheqd --entrypoint bash ghcr.io/cheqd/cheqd-node:0.4.0
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

And after that all the commands from the [flow](#demo-flow-for-sending-did-to-the-testnet) can be called.

P.S. the case with `docker` can be used only for demonstration purposes, cause after closing the container all the data will be lost.
For production purposes, maybe it would be great to have an image with Ubuntu 20.04 and operator's keys inside.
