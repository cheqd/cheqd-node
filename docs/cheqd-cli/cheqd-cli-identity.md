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