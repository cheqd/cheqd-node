# Using cheqd Cosmos CLI to manage DID Documents

## Overview

[cheqd Cosmos CLI](README.md) can be used for creating, updating DIDDocs on the ledger side


## Identity-related commands in cheqd CLI

### Querying DID Document

#### Command

```bash
cheqd-noded query cheqd did <id> --node http://localhost:26657
```

#### Example

```bash
cheqd-noded query cheqd did did:cheqd:testnet:zJ5EDiiiKWDyo79n --node http://nodes.testnet.cheqd.network:26657
```

### Create DID Document

#### Arguments

* `DIDDoc_in_JSON`: A string with a new DID Document in Json format.
    Base example:
    ```
    {
      "id": "<DID>",
      "verification_method": [{
        "id": "<KEY_ID>'",
        "type": "Ed25519VerificationKey2020",
        "controller": "<DID>",
        "public_key_multibase": "<ALICE_VER_PUB_MULTIBASE_58>"
      }],
      "authentication": [
        "<KEY_ID>"
      ]
    }
    ```
* `did_verification_method_id`: key identifier(`<KEY_ID>`) in the following format: `<DID>#<key-alias>`. Exapmle: `did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1`.
* `--ver-key`: Base64 encoded ed25519 private key to sign identity message with. A pair for the key from DID Document. \
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
                                     "public_key_multibase": "z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE"\
                                   }],\
                                   "authentication": [\
                                     "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"\
                                   ]\
                                 }' "id:cheqd:testnet:zJ5EDiiiKWDyo79n#key1" \
  --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" \
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 50ncheq
```

### Update DID Document

#### Arguments

* `DIDDoc_in_JSON`: A string with DID Document in Json format.  `id` is not changeable field and mast be used from creation transaction.
    Base example:
    ```
    {
      "id": "<DID>",
      "verification_method": [{
        "id": "<KEY_ID>'",
        "type": "Ed25519VerificationKey2020",
        "controller": "<DID>",
        "public_key_multibase": "<ALICE_VER_PUB_MULTIBASE_58>"
      }],
      "authentication": [
        "<KEY_ID>"
      ]
    }
    ```
* `did_verification_method_id`: key identifier in the following format: `<DID>#<key-alias>`. Exapmle: `did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1`.
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
  --from my_account --node http://nodes.testnet.cheqd.network:26657 --chain-id cheqd-testnet-4 --fees 50ncheq
```

## Crypto commands in cheqd CLI - EXPERIMENTAL FUNCTIONALITY - Do not use in production!


#### Command for generating public and private part of verification key

```bash
cheqd-noded debug ed25519 random
```

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


## Demo flow for sending DID to the testnet

As an example, let it be DID `did:cheqd:testnet:zJ5EDiiiKWDyo79n`

1. Generate verification key:
```bash
cheqd-noded debug ed25519 random
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

- operator address: `address: cheqd1vjuh4fjkcq0c02qullrt27z822gpn06sah2elh`
- mnemonic phrase ( 24 words ):
`crawl field same drill indoor olympic tank lamp range olive announce during pact idea fall canal sauce film attend response mammal bounce stable suffer`

Having this mnemonic phrase the usr is able to restore their keys whenever they want. For continue playing a user needs to run:

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

And after that all the commands from the flow can be called.

P.S. the case with `docker` can be used only for demonstration purposes, cause after closing the container all the data will be lost.
For production purposes, maybe it would be great to have an image with Ubuntu 20.04 and operator's keys inside.