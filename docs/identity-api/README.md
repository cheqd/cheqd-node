# Client-app Identity APIs

## Overview

This page describes how identity domain transactions need to be implemented by client-side applications/libraries such as [`cheqd-sdk`](https://github.com/cheqd/cheqd-sdk) (forked from [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools)).

Details on how identity transactions are defined is available in [ADR 002: Identity entities and transactions](../../architecture/adr-list/adr-002-cheqd-did-method.md).

## Ledger transactions/operations

### Base write flow

1. **Build a request** _Example_: `build_create_did_request(id, verkey, alias)`
2. **Sign the request using DID key** _Example_:  `indy_crypto_sign(did, verkey)`
3. **Build a transaction with the request from previous step** _Example_: `build_tx(pool_alias, pub_key, builded_request, account_number, account_sequence, max_gas, max_coin_amount, denom, timeout_height, memo)`
4. **Sign the transaction** _Example_: `cheqd_keys_sign(wallet_handle, key_alias, tx)`.
5. **Broadcast a signed transaction** _Example_: `broadcast_tx_commit(pool_alias, signed)`.

#### Example output

```protobuf
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
      log: "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cheqd1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cheqd1pvnjjy3vz0ga6hexv32gdxydzxth7f86mekcpg\"},{\"key\":\"sender\",\"value\":\"cheqd1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"amount\",\"value\":\"500000ncheq\"}]}]}]",
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

### Create DID

Used to create a new DID. The unique ID is generated client-side by VDR Tools SDK, but checked for uniqueness on the ledger before being committed.

#### Input parameters

* `id` (string): Target DID as Base58-encoded string with 16 or 32 byte long unique identifier.
* `verkey` (string): All Verification Method key(s) linked to this DID and its DID controller(s). At least one Verification Method key *must* be defined.

#### Method call

The `CreateDidRequest` must be signed by the controller DID(s) and their associated key(s) defined. It is invoked as follows:

```rust
build_create_did_request(id, verkey)
```

If successful, the fully-qualified DID with unique identifier is returned as a `key` string. This `key` is how the DID is uniquely referenced when it is stored in the ledger state.

```rust
CreateDidResponse {
    "key": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue" 
}
```

#### Example response data

```jsonc
{
  "data": {
    "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue", // DID with unique defined in the method call parameters
    "controller": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue"],
    "verificationMethod": [
        {
          "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey", // id#verkey
          "type": "Ed25519VerificationKey2020",
          "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
          "publicKeyMultibase": "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf"
        }
    ],
    "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey"]
  },
  "signatures": {
    "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba"
    // Multiple verification methods and corresponding signatures can be added here
  },
  "metadata": {}
}
```

### Update DID

Update an existing DID on the cheqd network ledger. The fully-qualified DID must be stored on ledger already, otherwise the request fails.

#### Input parameters

* `id` (string): Target DID as Base58-encoded string with 16 or 32 byte long unique identifier.
* `verkey` (string): All Verification Method key(s) linked to this DID and its DID controller(s).
* `versionId` (string): Transaction hash of the last applicable DIDDoc version.

#### Method call

All DID controller(s) from `controller` field must already be stored on ledger, from a previous `CreateDidRequest`.

The DID update request must be signed by DIDs from `controller` field, or if `controller` does not exist, by at least one key from `authentication`.

```rust
build_update_did_request(id, verkey, version_id)
```

If successful, the fully-qualified DID with unique identifier is returned as a `key` string. This `key` is how the DID is uniquely referenced when it is stored in the ledger state.

```rust
UpdateDidRequest {
    "key": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue" 
}
```

#### Example response data

```jsonc
{
  "data": {
    "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
    "controller": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue"],
    "verificationMethod": [
      {
        "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey",
        "type": "Ed25519VerificationKey2020", // external (property value)
        "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
        "publicKeyMultibase": "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf"
      }
    ],
    "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey"],
    "versionId": "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E"
  },
  "signatures": {
    "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba"
    // Multiple verification methods and corresponding signatures can be added here
  },
  "metadata": {}
}
```

### Get DID

#### cheqd-sdk function

`build_query_get_did(id)`

* `id` (base58-encoded string): Target DID as base58-encoded string for 16 or 32 byte DID value.

#### Example response

```protobuf
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

```jsonc
QueryGetDidResponse{
        "did": {
               "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
               "controller": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue"],
               "verificationMethod": [
                  {
                    "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey",
                    "type": "Ed25519VerificationKey2020", // external (property value)
                    "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
                    "publicKeyMultibase": "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf"
                  }
              ],
              "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey"],
        },
        "metadata": {
              "created": "2020-12-20T19:17:47Z",
              "updated": "2020-12-20T19:19:47Z",
              "deactivated": false,
              "versionId": "N22KY2Dyvmuu2PyyqSFKueN22KY2Dyvmuu2PyyqSFKue",
        }
}
```

### Create SCHEMA

#### cheqd-sdk function

`build_create_schema_request(id, controller, version, name, attr_names)`

#### Example response

```jsonc
{
    "data": {
            "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
            "type": "CL-Schema",
            "controller": ["did:cheqd:mainnet:GEzcdDLhCpGCYRHW82kjHd"]
            "version": "1.0",
            "name": "Degree",
            "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
    "signatures": {
            "did:cheqd:mainnet:GEzcdDLhCpGCYRHW82kjHd#verkey": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba"
            // Multiple verification methods and corresponding signatures can be added here
    },
    "metadata": {}
}
```

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with cheqd DID Method prefix `did:cheqd:<namespace>:` and a resource
type at the end.
- **`type`**: String with a schema type. Now only `CL-Schema` is supported.
- **`attr_names`**: Array of attribute name strings (125 attributes maximum)
- **`name`**: Schema's name string
- **`version`**: Schema's version string
- **`controller`**: DIDs list of strings or only one string of a schema
  controller(s). All DIDs must exist.

#### Response format

```jsonc
CreateSchemaResponse {
        "key": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema" 
}
```

* `key`\(string\): A key is used to store this schema in a state

#### Request validation

* All DIDs from `controller` field must already be in a ledger created by `CreateDidRequest`
* Schema create request must be signed by DIDs from `controller` field. 

### Get Schema

#### cheqd-sdk function

`build_query_get_schema(id)`

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with cheqd DID Method prefix `did:cheqd:<namespace>:` and a resource
  type at the end.

#### Example response

```protobuf
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

```jsonc
QueryGetSchemaResponse{
        "schema": {
            "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
            "type": "CL-Schema",
            "controller": ["did:cheqd:mainnet:GEzcdDLhCpGCYRHW82kjHd"],
            "version": "1.0",
            "name": "Degree",
            "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
        },
}
```

## CRED\_DEF

### Create Credential Definition

#### cheqd-sdk function

`build_create_cred_def_request(cred_def, schema_id, signature_type, tag)`

#### Example response

```jsonc
CreateCredDefRequest 
{
    "data": {   
                "id": "did:cheqd:mainnet:5ZTp9g4SP6t73rH2s8zgmtqdXyT?service=CL-CredDef",
                "type": "CL-CredDef",
                "controller": ["did:cheqd:mainnet:123456789abcdefghi"],
                "schema_id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
                "tag": "some_tag",
                "value": {
                  "primary": "...",
                  "revocation": "..."
                }
            },
    "signatures": {
            "did:cheqd:mainnet:GEzcdDLhCpGCYRHW82kjHd#verkey": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba"
            // Multiple verification methods and corresponding signatures can be added here
    },
    "metadata": {}
}
```

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with Cheqd
  DID Method prefix `did:cheqd:<namespace>:` and a resource
  type at the end.
- **`value`** (dict): Dictionary with Credential Definition's data if
  `signature_type` is `CL`:
  - **`primary`** (dict): Primary credential public key
  - **`revocation`** (dict, optional): Revocation credential public key
- **`schema_id`** (string): `id` of a Schema the credential definition is created
  for.
- **`type`** (string): Type of the credential definition (that is
  credential signature). `CL-CredDef` (Camenisch-Lysyanskaya) is the only
  supported type now. Other signature types are being explored for future releases.
- **`tag`** (string, optional): A unique tag to have multiple public keys for
  the same Schema and type issued by the same DID. A default tag `tag` will be
  used if not specified.
- **`controller`**: DIDs list of strings or only one string of a credential
  definition controller(s). All DIDs must exist.
  

#### Response format

```jsonc
CreateCredDefResponse {
        "key": "did:cheqd:mainnet:5ZTp9g4SP6t73rH2s8zgmtqdXyT?service=CL-CredDef" 
}
```

* `key`(string): A unique key that is used to store this Credential Definition in a state


#### Request validation

* All DIDs from `controller` field must already be in a ledger created by `CreateDidRequest`
* Cred Def create request must be signed by DIDs from `controller` field.

### Get Credential Definition

#### cheqd-sdk function

`build_query_get_cred_def(id)`

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with cheqd DID Method prefix `did:cheqd:<namespace>:` and a resource
  type at the end.
  
#### Example response

```protobuf
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`: Path for RPC endpoint for cheqd pool
* `data`: Query with an entity key from a state. String `cred_def:<owner>:<schema_id>:<tag>:<signature_type>` encoded to bytes
* `height`: Ledger height \(size\). `None` for auto calculation
* `prove`: Boolean value. `True` for getting state proof in a pool response.

#### Response format

```jsonc
QueryGetCredDefResponse{
    "cred_def": {
                "id": "did:cheqd:mainnet:5ZTp9g4SP6t73rH2s8zgmtqdXyT?service=CL-CredDef",
                "type": "CL-CredDef",
                "controller": ["did:cheqd:mainnet:123456789abcdefghi"],
                "schema_id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
                "tag": "some_tag",
                "value": {
                    "primary": "...",// Primary
                    "revocation": "..." // Revocation registry
                }
        },
}
```
