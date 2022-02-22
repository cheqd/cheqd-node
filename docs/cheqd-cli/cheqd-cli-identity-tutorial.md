# Tutorial for creating a DID + DIDDoc with cheqd CLI

## Overview

The purpose of this document is to outline how someone can create a DIDDoc on the cheqd network.

  

This tutorial uses cheqd node CLI to send a DIDDoc.

  

## Pre-requisites


In order to create a DIDDoc using the instructions outlined in this tutorial, you must be using Ubuntu 20.04 terminal. You'll find all the information required to setup your Ubuntu 20.04 terminal at the end of this tutorial.

  

If you don't currently have Ubuntu 20.04 installed on your machine you can use VirtualBox or [Docker](#requirements-from-os-side)

  

Please ensure you are running the correct version of testnet. You can check which is the current version of testnet [here](https://rpc.testnet.cheqd.network/abci_info?).


## How to send a DIDDoc to the testnet

  
### 1\. Generate verification key:

  

First we'll need to generate a verification key:

  

```bash
$ cheqd-noded debug ed25519 random >> keys.txt
```

  

The result should look like the following:

  

```bash
$ cat keys.txt
{"pub_key_base_64":"MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=","pub_key_multibase_58":"",\
"priv_key_base_64":"FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEpC+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q=="}
```

  

### 2\. Get multibase58 string


It needs for inserting it into the DID-doc (`public_key_multibase` field in `verification_method` section)

  

```bash
$ cheqd-noded debug encoding base64-multibase58 <pub_key_base_64>
```

  

Based on the working example in this tutorial the result will be:

  

```bash
$ cheqd-noded debug encoding base64-multibase58 MnrTheU+vCrN3W+WMvcpBXYBG6D1HrN5usL1zS6W7/k=
z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE
```

  

And the response will be:

  

```text
z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE
```

  

### 3\. Create unique-id for our DID

  

To create a `unique-id` for our DID we can use first 32 symbols of `multibase58` representation of our public key as \`unique-id\`.

For example, we can truncate previous one:

```bash
$ printf '%.32s\n' `cheqd-noded debug encoding base64-multibase58 <pub_key_base_64>`
```

  

The result for our example will be `z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe` , so let's use it as our `unique-id` in our DIDDoc.

  

### 4\. Compile DIDDoc

  

Next we can compile our DIDDoc.

  

Copy and paste the template below into your terminal. We will add additional required information into the blank fields `<xxxxx>` in the next steps.

  

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
  ],
  "service": [{
    "id":"did:cheqd:<namespace>:<unique-id>#<service-key>",
    "type": "LinkedDomains",
    "serviceEndpoint": "<URI-to-object>"
  }]
}
```

Within this template we will be required to enter a number of fields

  

Where:

*   `<namespace>` - for now it can `testnet` or `mainnet` . For this flow we use `testnet`
*   `<unique-id>` - identifier, created on the [step](#3-create-unique-id-for-our-did)
*   `<key-alias>` - a key alias for the verification method identifier
*   `<verification-public-key-multibase>` - result of this [step](#2-get-multibase58-string)
*   `<auth-key-alias>` - alias of authentication key.
*   `<service-key>` - alias for service. More information [here](../../architecture/adr-list/adr-002-cheqd-did-method.md#service)
*   `<URI-to-object>` - URI according to the spec [rfc3986](https://www.rfc-editor.org/rfc/rfc3986) 

In our example:

*   `did:cheqd:<namespace>:<unique-id>` - would be `did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe`
*   `did:cheqd:<namespace>:<unique-id>#<key-alias>` - `did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe#key1`
*   `<verification-public-key-multibase>` - key from this [step](#2-get-multibase58-string). As result `z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE`
*   `did:cheqd:<namespace>:<unique-id>#<auth-key-alias>` - let it be:

`did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe#key1`

  

#### After these preparations, the base DIDDoc will look like:

  

```json
{
  "id": "did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe",
  "verification_method": [
    {
      "id": "did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe#key1",
      "type": "Ed25519VerificationKey2020",
      "controller": "did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe",
      "public_key_multibase": "z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe3Nmod35uua9TE"
    }
  ],
  "authentication": [
    "did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe#key1"
  ],
  "service": [{
    "id":"did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe#linked-domain",
    "type": "LinkedDomains",
    "serviceEndpoint": "https://bar.example.com"
  }]
}
```

  

We recommend you store this DIDDoc in a separate file, like `diddoc.json` and inject it while running the command for sending.

  

### 5\. Send DIDDoc to the pool

  

Now that we have our DIDDoc prepared we can send it to the pool.

  

We can use the following command to send the DIDDoc:

  

```bash
$ cheqd-noded tx cheqd create-did "$(cat diddoc.json)" "did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1" --ver-key "FxaJOy4HFoC2Enu1SizKtU0L+hmBRBAEp+B4TopfQoyetOF5T68Ks3db5Yy9ykFdgEboPUes3m6wvXNLpbv+Q==" --from <alias-to-cosmos-key>  --node https://rpc.testnet.cheqd.network:443 --chain-id cheqd-testnet-4 --fees 5000000ncheq
```

  

Where:

*   `"did:cheqd:testnet:zJ5EDiiiKWDyo79n#key1"` is the `id` of `verification_method` section
*   `--ver-key` - is from `keys.txt` file from the [step](#1-generate-verification-key), `priv_key_base_64` field.
*   `--from` - should be an alias of your cosmos keys.

  

After you execute the command you will receive `"code": 0"`if the DID was successfully written to the ledger. We can do a full query to check this as well.
In case of other error codes field `raw_logs` can help with figuring out the case. For example:

```text
"code":1201,"data":"","raw_log":"failed to execute message; message index: 0: id:cheqd:testnet:fcbarcelona: DID Doc not found"
``` 
  

### 6\. Check that DID was successfully written to the ledger

  

Finally, to check that the DID was successfully written we can use the following query:

  

```bash
$ cheqd-noded query cheqd did "<identifier-of-your-DIDDoc>" --node https://rpc.testnet.cheqd.network:443
```

  

where:

*   `<identifier-of-your-DIDDoc>` - identifier with template `"did:cheqd:<namespace>:<unique-id>"` and `<unique-id>` is from [step](#3-create-unique-id-for-our-did)

  

In our example:

```bash
$ cheqd-noded query cheqd did "did:cheqd:testnet:z4Q41kvWsd1JAuPFBff8Dti7P6fLbPZe" --node https://rpc.testnet.cheqd.network:443
```


**Congratulations! You've created your first, of many, DIDDoc on cheqd!**

  

## Requirements from OS side

  

Our target OS system is Ubuntu 20.04.

  

In this case, for running demo flow we can use variants: Virtualbox or docker.

For example, let it be a docker image, cause it's the most fastest way to start playing.

The next command can help:

  

```bash
$ docker run -it --rm -u root --entrypoint bash ghcr.io/cheqd/cheqd-node:0.4.0
```

  

After that, we need to install needed package for process SSL certificates:

  

```bash
# apt update && apt install ca-certificates -y
```

  

Also, it can help to setup your favourite editor, for example `vim` :

  

```bash
# apt install vim -y 
```

  

The next step is to change user to `cheqd` and restore operator's keys:

  

```bash
# su cheqd
```

  

```bash
$ cheqd-noded keys add <cheqd-operator-name> --recover
```

  

where, `cheqd-operator-name` it's name of alias for storing your keys locally, whatever you want.
