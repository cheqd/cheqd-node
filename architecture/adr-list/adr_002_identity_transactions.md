# ADR 002: Identity entities and transactions

## Status

| Category | Status |
| :--- | :--- |
| **ADR Stage** | PROPOSED |
| **Implementation Status** | Not Implemented |

## Summary

This ADR summarises the identity entities, queries, and transaction types for the cheqd network. These recommendations are based on the design patterns currently used by [Hyperledger Indy](https://github.com/hyperledger/indy-node), a blockchain built for self-sovereign identity \(SSI\).

## Context

Hyperledger Indy contains the following [identity domain transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md):

1. `NYM`
2. `ATTRIB`
3. `SCHEMA`
4. `CRED_DEF`
5. `REVOC_REG_DEF`
6. `REVOC_REG_ENTRY`

Our aim is to replicate similar transactions on `cheqd-node` to allow existing SSI software designed to work with Hyperledger Indy to be compatible with the cheqd network.

### Changes proposed from existing Hyperledger Indy transactions

We have assessed the existing Hyperledger Indy transactions and recommend the following changes to be made.

#### Rename `NYM` transactions to `DID` transactions

[**NYM** is the term used by Hyperledger Indy](https://hyperledger-indy.readthedocs.io/projects/node/en/latest/transactions.html#nym) for [Decentralized Identifiers \(DIDs\)](https://www.w3.org/TR/did-core/) that are created on ledger. A DID is typically the identifier that is associated with a specific organisation issuing/managing SSI credentials.

For the sake of explaining with similar concepts to current Hyperledger Indy implementations, on the `cheqd-testnet` these are still called NYMs. Transactions to add a DID to the cheqd network ledger are currently called NYM transactions.

Our proposal is to change the term `NYM` in transactions to `DID`, which would make understanding the context of a transaction easier to understand. This change will bring transactions better in-line with World Wide Web Consortium \(W3C\) terminology.

#### Remove `role` field from DID transaction

Hyperledger Indy is a public-permissioned distributed ledger. As `cheqd-node` is based on a public-permissionless network based on the [Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk), the need for having a `role` type is not necessary.

_**Note**: Hyperledger Indy also contains other transaction types beyond the ones listed above, but these are currently not in scope for implementation in `cheqd-node`. They will be considered for inclusion later in the product roadmap._

## Decision

### General structure of transaction requests

All identity requests will have the following format:

```text
{
    "data": { <request data for writing a transaction to the ledger> },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* **`data`**: Data requested to be written to the ledger, specific for each request type.
* **`owner`**: Owner identifier \(DID\) for this entity. This could be a new DID or an existing DID, for existing entities.
* **`signature`**: `data` should be signed by `owner` private key.
* **`metadata`**: Dictionary with additional metadata fields. Empty for now. This fields provides extensibility in the future, e.g., it can contain `protocolVersion` or other relevant metadata associated with a request.

## List of transactions and details

### `DID` transactions

[Decentralized Identifiers \(DIDs\) are a W3C specification](https://www.w3.org/TR/did-core/) for identifiers that enable verifiable, decentralized digital identity.

The request can be used for creation of new DIDs, setting, and rotation of verification key.

1. **`dest` \(base58-encoded string\):**

   Target DID as base58-encoded string for 16 or 32 byte DID value. It may differ from `identifier` metadata field, where `identifier` is the DID of the submitter. If they are equal \(in permissionless case\), then the transaction must be signed by the newly created `verkey`.

   _Example_: `identifier` is a DID of an Endorser creating a new DID, and `dest` is a newly created DID.

2. **`verkey` \(base58-encoded string, possibly starting with "~"; optional\):**

   Target verification key as base58-encoded string. It can start with "~", which means that it's abbreviated `verkey` and should be 16 bytes long when decoded, otherwise it's a full `verkey` which should be 32 bytes long when decoded. If not set, then either the target identifier \(`dest`\) is 32-bit cryptonym CID \(this is deprecated\), or this is a user under guardianship \(doesn't own the identifier yet\).

3. **`alias` \(string; optional\):**

   Alias for the DID

   ```text
       {
         "alias": "Alice DID",
         "dest": "GEzcdDLhCpGCYRHW82kjHd",
         "verkey": "~HmUWn928bnFT6Ephf65YXv"
       }
   ```

#### **Update DID**

If there is no DID transaction with the specified DID \(`dest`\), it is considered as a creation request for a new DID.

If there is a DID transaction with the specified DID \(`dest`\), then this is update of existing DID. In this case, we can specify only the values we would like to override. All unspecified values remain the same. E.g., if a key rotation needs to be performed, the owner of the DID needs to send a DID transaction request with `dest`, `verkey` only. `alias` will stay the same.

**Note:** Fields `dest` and `owner` should have the same value.

#### State format

`dest -> {(alias, dest, verkey), last_tx_hash, last_update_timestamp }`

_Request Example_:

```text
CreateDidRequest
{
    "data": {
              "alias": "Alice DID",
              "dest": "GEzcdDLhCpGCYRHW82kjHd",
              "verkey": "~HmUWn928bnFT6Ephf65YXv"
             }
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

_Reply Example_:

```text
CreateDidResponse
{}
```

_DID transaction format:_

```text
Did
{
    "alias": "Alice DID",
    "dest": "GEzcdDLhCpGCYRHW82kjHd",
    "verkey": "~HmUWn928bnFT6Ephf65YXv"
}
```

### ATTRIB

Adds a new Attribute or updates an existing Attribute to a DID record.

* **`dest` \(base58-encoded string\):**

  Target DID as base58-encoded string for 16 or 32 byte DID value.

* **`raw` \(json; mutually exclusive with `hash` and `enc`\):**

  Raw data is represented as JSON, where the key is attribute name and value is attribute value.

**Note:** ATTRIB **can** be updated

#### State format

 `dest -> {(dest, raw), last_tx_hash, last_update_timestamp }`

_Request Example_:

```text
CreateAttribRequest
{
    "data": {
              "dest": "N22KY2Dyvmuu2PyyqSFKue",
              "raw": "{"name": "Alice"}",
             }
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

_Reply Example_:

```text
CreateAttribResponse
{}
```

### SCHEMA

This transaction is used to create a Schema associated with credentials.

It is not possible to update an existing Schema, to ensure the original schema used to issue any credentials in the past are always available.

If a Schema evolves, a new schema with a new version or name needs to be created.

* **`data`:**

  Dictionary with Schema's data:

  * **`attr_names`**: Array of attribute name strings \(125 attributes maximum\)
  * **`name`**: Schema's name string
  * **`version`**: Schema's version string

**Note:** SCHEMA **cannot** be updated

#### State format

`(version, name, owner) -> {(version, name, attr_names), tx_hash, tx_timestamp }`

_Request Example_:

```text
{
    "data": {
            "version": "1.0",
            "name": "Degree",
            "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
    "owner": "L5AD5g65TDQr1PPHHRoiGf",
    "signature": "5ZTp9g4SP6t73rH2s8zgmtqdXyTuSMWwkLvfV1FD6ddHCpwTY5SAsp8YmLWnTgDnPXfJue3vJBWjy89bSHvyMSdS",
    "metadata": {}
}
```

_Reply Example_:

```text
{}
```

### CRED\_DEF 

Adds a Credential Definition \(in particular, public key\), which is created by an Issuer and published for a particular Credential Schema.

It is not possible to update `data` in existing Credential Definitions. If a Credential Definition needs to be evolved \(for example, a key needs to be rotated\), then a new Credential Definition needs to be created by a new Issuer DID \(`owner`\).

* **`cred_def` \(dict\):**

  Dictionary with Credential Definition's data:

  * **`primary`** \(dict\): Primary credential public key
  * **`revocation`** \(dict\): Revocation credential public key

* **`ref` \(string\):**

  Hash of a Schema transaction the credential definition is created for.

* **`signature_type` \(string\):**

  Type of the credential definition \(that is credential signature\). `CL` \(Camenisch-Lysyanskaya\) is the only supported type now. Other signature types are being explored for future releases.

* **`tag` \(string, optional\):**

  A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

**Note**: CRED\_DEF **cannot** be updated.

#### State format

`(owner, signature_type, ref, tag) -> {(primary, revocation), tx_hash, tx_timestamp }`

_Request Example_:

```text
{
    "data": {
        "signature_type": "CL",
        "ref": 5ZTp9g4SP6t73rH2s8zgmtqdXyT,
        "tag": "some_tag",    
        "cred_def": {
            "primary": ....,
            "revocation": ....
        }
    },

    "owner": "L5AD5g65TDQr1PPHHRoiGf",
    "signature": "5ZTp9g4SP6t73rH2s8zgmtqdXyTuSMWwkLvfV1FD6ddHCpwTY5SAsp8YmLWnTgDnPXfJue3vJBWjy89bSHvyMSdS",
    "metadata": {}
}
```

_Reply Example_:

```text
{}
```

## References

* [Hyperledger Indy Identity transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)
* [W3 DID Spec](https://www.w3.org/TR/did-core/)

