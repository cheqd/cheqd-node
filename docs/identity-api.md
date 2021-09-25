# Client-app Identity APIs

## Overview

This page describes how identity domain transactions need to be implemented by client-side applications/libraries such as [`cheqd-sdk`](https://github.com/cheqd/cheqd-sdk) \(forked from [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools)\).

Details on how identity transactions are defined is available in [ADR 002: Identity entities and transactions](../architecture/adr-list/adr_002_identity_transactions.md).

### Base write flow

1. **Build a request** _Example_: `build_create_did_request(id, verkey, alias)`
2. **Sign the request using DID key** _Example_:  `indy_crypto_sign(did, verkey)`
3. **Build a transaction with the request from previous step** _Example_: `build_tx(pool_alias, pub_key, builded_request, account_number, account_sequence, max_gas, max_coin_amount, denom, timeout_height, memo)`
4. **Sign the transaction** _Example_: `cheqd_keys_sign(wallet_handle, key_alias, tx)`. 
5. **Broadcast a signed transaction** _Example_: `broadcast_tx_commit(pool_alias, signed)`.

#### Response format

```text
  Response {
   check_tx: TxResult {
      code: 0,
      data: None,
      log: "",
      info: "",
      gas_wanted: 0,
      gas_used: 0,
      events: [
      ],
      codespace: ""
   },
   deliver_tx: TxResult {
      code: 0,
      data: Some(Data([...])),
      log: "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1pvnjjy3vz0ga6hexv32gdxydzxth7f86mekcpg\"},{\"key\":\"sender\",\"value\":\"cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"amount\",\"value\":\"100cheq\"}]}]}]",
      info: "",
      gas_wanted: 0,
      gas_used: 0,
      events: [...], 
      codespace: ""
   },
   hash: "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E",
   height: 353
}
```

**`hash`** : Transaction hash

**`height`**: Ledger height

## DID transactions

### Create DID

#### cheqd-sdk function

`build_create_did_request(id, verkey, alias)`

#### Request format

```text
CreateDidRequest 
{
    "data": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "DID for Alice"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `id` \(base58-encoded string\): Target DID as base58-encoded string for 16 or 32 byte DID value
* `verkey` \(base58-encoded string, possibly starting with "~"; optional\): Target verification key. It can start with "~", which means that it is an abbreviated `verkey` and should be 16 bytes long when decoded. Otherwise, it's a full `verkey` which should be 32 bytes long when decoded.
* `alias` \(string; optional\)

#### Response format

```text
CreateDidResponse {
    "key": "did:GEzcdDLhCpGCYRHW82kjHd" 
}
```

* `key`\(string\): A unique key is used to store this DID in a state

#### Response validation

* `CreateDidRequest` must be signed by the DID from `id` field. It means that this DID must be an owner of this DID transaction.

### Update DID

#### cheqd-sdk function

`build_update_did_request(id, verkey, alias)`

#### Request format

```text
UpdateDidRequest 
{
    "data": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "DID for Alice"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `id` \(base58-encoded string\): Target DID as base58-encoded string for 16 or 32 byte DID value.
* `verkey` \(base58-encoded string, possibly starting with "~"; optional\): Target verification key. It can start with "~", which means that it is an abbreviated `verkey` and should be 16 bytes long when decoded. Otherwise, it's a full `verkey` which should be 32 bytes long when decoded.
* `alias` \(string; optional\).

#### Response format

```text
UpdateDidResponse {
    "key": "did:GEzcdDLhCpGCYRHW82kjHd" 
}
```

* `key`\(string\): A unique key is used to store this DID in a state

#### Response validation

* A transaction with `id` from `UpdateDidRequest`must already be in a ledger created by `CreateDidRequest`
* `UpdateDidRequest` must be signed by the DID from `id` field. It means that this DID must be an owner of this DID transaction.

### Get DID

#### cheqd-sdk function

`build_query_get_did(id)`

* `id` \(base58-encoded string\): Target DID as base58-encoded string for 16 or 32 byte DID value.

#### Request format

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`: Path for RPC endpoint for cheqd pool
* `data`: Query with an entity key from a state. String `did:<id>` encoded to bytes
* `height`: Ledger height \(size\). `None` for auto calculation
* `prove`:  Boolean value. `True` for getting state proof in a pool response

#### Response format

```text
QueryGetDidResponse{
        "did": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "DID for Alice"
             },
}
```

## ATTRIB transactions

### Create ATTRIB

#### cheqd-sdk function

`build_create_attrib_request(did, raw)`

#### Request format

```text
CreateAttribRequest 
{
    "data": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `did` \(base58-encoded string\): Target DID as base58-encoded string for 16 or 32 byte DID value.
* `raw` \(JSON; mutually exclusive with `hash` and `enc`\): Raw data represented as JSON, where the key is attribute name and value is attribute value.

#### Response format

```text
CreateAttribResponse {
    "key": "attrib:GEzcdDLhCpGCYRHW82kjHd" 
}
```

* `key`\(string\): A unique key is used to store these attributes in a state

#### Response validation

* A DID transaction with `id` from `UpdateAttribRequest`must already be in a ledger created by `CreateDidRequest`
* `CreateAttribRequest` must be signed by the DID from `did` field. It means that this DID must be an owner of this ATTRIB transaction.

### Update ATTRIB

#### cheqd-sdk function

`build_update_attrib_request(id, raw)`

#### Request format

```text
UpdateAttribRequest 
{
    "data": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `did` \(base58-encoded string\): Target DID as base58-encoded string for 16 or 32 byte DID value.
* `raw` \(JSON; mutually exclusive with `hash` and `enc`\): Raw data represented as JSON, where the key is attribute name and value is attribute value.

#### Response format

```text
UpdateAttribResponse {
        "key": "attrib:GEzcdDLhCpGCYRHW82kjHd" 
}
```

* `key`\(string\): A unique key is used to store these attributes in a state

#### Response validation

* A DID transaction with `id` from `UpdateAttribRequest`must already be in a ledger created by `CreateDidRequest`
* `UpdateAttribRequest` must be signed by  DID from `did` field. It means that this DID must be an owner of this ATTRIB transaction.

### Get ATTRIB

#### cheqd-sdk function

`build_query_get_attrib(did)`

* `did` \(base58-encoded string\) Target DID as base58-encoded string for 16 or 32 byte DID value.

#### Request format

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`: Path for RPC endpoint for cheqd pool
* `data`: Query with an entity key from a state. String `did:<id>` encoded to bytes
* `height`: Ledger height \(size\). `None` for auto calculation
* `prove`:  Boolean value. `True` for getting state proof in a pool response

#### Response format

```text
QueryGetAttribResponse{
        "attrib": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
}
```

## SCHEMA transactions

### Create SCHEMA

#### cheqd-sdk function

`build\_create\_schema\_request\(version, name, attr\_names\)`

#### Request format

```text
CreateSchemaRequest 
{
    "data": {
            "version": "1.0",
            "name": "Degree",
            "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `attrNames`\(array\): Array of attribute name strings \(125 attributes maximum\)
* `name`\(string\): Schema's name string
* `version`\(string\): Schema's version string

#### Response format

```text
CreateSchemaResponse {
        "key": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0" 
}
```

* `key`\(string\): A key is used to store this schema in a state

#### Response validation

* A SCHEMA transaction with DID from `owner` field must already be in a ledger created by `CreateDidRequest`
* `CreateSchemaRequest` must be signed by  DID from `owner` field. 

### Get Schema

#### cheqd-sdk function

`build\_query\_get\_schema\(name, version, owner\)`

* `name`\(string\): Schema's name string
* `version`\(string\): Schema's version string
* `owner` \(string\): Schema's owner did

#### Request format

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`: Path for RPC Endpoint for cheqd pool
* `data`: Query with an entity key from a state. String `schema:<owner>:<name>:<version>` encoded to bytes
* `height`: Ledger height \(size\). `None` for auto calculation;
* `prove`: Boolean value. `True` for getting state proof in a pool response. 

#### Response format

```text
QueryGetSchemaResponse{
        "attrib": {
                "version": "1.0",
                "name": "Degree",
                "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
}
```

## CRED\_DEF

### Create Credential Definition

#### cheqd-sdk function

`build\_create\_cred\_def\_request\(cred\_def, schema\_id, signature\_type, tag\)`

#### Request format

```text
CreateCredDefRequest 
{
    "data": {
                "signatureType": "CL",
                "schema_id": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0",
                "tag": "some_tag",    
                "cred_def": {
                    "primary": ....,
                    "revocation": ....
            },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `cred_def` \(dict\): Dictionary with Credential Definition's data:
  * `primary` \(dict\): Primary credential public key
  * `revocation` \(dict\): Revocation credential public key
* `schema_id` \(string\): Schema's key from a state
* `signatureType` \(string\): Type of the Credential Definition \(that is credential signature\). `CL` \(Camenisch-Lysyanskaya\) is the only supported type now.
* `tag` \(string, optional\): A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

#### Response format

```text
CreateCredDefResponse {
        "key": "cred_def:GEzcdDLhCpGCYRHW82kjHd:schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0:some_tag:CL" 
}
```

* `key`\(string\): A unique key that is used to store this Credential Definition in a state

#### Response validation

* A CRED\_DEF transaction with DID from `owner` field must already be in a ledger created by `CreateDidRequest`
* `CreateCredDefRequest` must be signed by  DID from `owner` field. 

### Get Credential Definition

#### cheqd-sdk function

`build\_query\_get\_cred\_def\(name, version, owner\)`

* `schema_id`\(string\): Schema's key from a state
* `signatureType`\(string\): Type of the Credential Definition \(that is credential signature\). CL \(Camenisch-Lysyanskaya\) is the only supported type now.
* `owner` \(string\): Credential Definition's owner DID
* `tag` \(string, optional\): A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

#### Request format

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`: Path for RPC endpoint for cheqd pool
* `data`: Query with an entity key from a state. String `cred_def:<owner>:<schema_id>:<tag>:<signatureType>` encoded to bytes
* `height`: Ledger height \(size\). `None` for auto calculation
* `prove`: Boolean value. `True` for getting state proof in a pool response. 

#### Response format

```text
QueryGetCredDefResponse{
        "cred_def": {
                "signatureType": "CL",
                "schema_id": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0",
                "tag": "some_tag",    
                "cred_def": {
                    "primary": ....,
                    "revocation": ....
         },
}
```

