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

### Query DID

Fetch DID Document and metadata associated with a specified DID. The fully-qualified DID must be stored on ledger already, otherwise the request fails.

#### Method call

The only input that the query/get DID method requires is the `id` associated with the DID that needs to be fetched.

```rust
build_query_get_did(id)
```

The request itself is wrapped inside an envelope that requires the following details:

```protobuf
Request 
{
    "path": "<pool-rpc-endpoint>",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

The values required above are:

* `path`: Path for the Tendermint RPC endpoint from which the response was returned.
* `data`: Query with an `key` that corresponds to a specific DID state, encoded to bytes.
* `height`: Block height from which the state is to be fetched. This can be set to `none` for auto-calculation that fetches from latest ledger state, while providing a previously block height can be used to query the DID state at a previous block height / equivalent time.
* `prove`: Boolean value (`true` or `false`) that defines whether the state proof should be fetched along with the response.

#### Example response data

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
