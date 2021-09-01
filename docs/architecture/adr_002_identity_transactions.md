# ADR-002: Identity transactions

## Status
PROPOSED | Not Implemented

## Summary

Identity entities for cheqd ledger according [Hyperledger Indy Identity transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)

## Context

Hyperledger Indy contains domain transactions NYM, ATTRIB, SCHEMA, CLAIM_DEF, REVOC_REG_DEF, REVOC_REG_ENTRY
We need to move them to Cheqd-node to allow SSI systems work with Cheqd.
But we need to make some changes in the previous domain model.
1) Rename NYM transaction to DID transaction
2) Remove `role` field from DID transaction because Cheqd is based on poof-of-stake permissionless model. The concept of role is mot applicable in the new context.


## Decision

All identity requests have the follow format:
```
{
    "data": { <request data for creating a transaction to the ledger> },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

- `data` - request data for putting it to ledger, specific for each request type.
- `owner` - owner identifier(DID) for this entity. New one or old one for existing entity.
- `signature` - data signed be `owner` private key.
- `metadata` - dictionary with additional metadata. Empty for now but can contain `protocolVersion` in future or other fields. 


Transactions list:

### DID

Decentralized identifiers (DIDs) are a new type of identifier that enables verifiable, decentralized digital identity. (see [specification](https://www.w3.org/TR/did-core/)).

The request can be used for
creation of new DIDs, setting, and rotation of verification key.

- `dest` (base58-encoded string):

  Target DID as base58-encoded string for 16 or 32 byte DID value.
  It may differ from `identifier` metadata field, where `identifier` is the DID of the submitter.
  If they are equal (in permissionless case), then the transaction must be signed by the newly created `verkey`.

  *Example*: `identifier` is a DID of an Endorser creating a new DID, and `dest` is a newly created DID.

- `verkey` (base58-encoded string, possibly starting with "~"; optional):

  Target verification key as base58-encoded string. It can start with "~", which means that
  it's abbreviated verkey and should be 16 bytes long when decoded, otherwise it's a full verkey
  which should be 32 bytes long when decoded. If not set, then either the target identifier
  (`dest`) is 32-bit cryptonym CID (this is deprecated), or this is a user under guardianship
  (doesn't own the identifier yet).

- `alias` (string; optional):

  DID's alias.

  ```
      {
        "alias": "Alice did",
        "dest": "GEzcdDLhCpGCYRHW82kjHd",
        "verkey": "~HmUWn928bnFT6Ephf65YXv"
      }
  ```
<b>Update DID:</b>

If there is no DID transaction with the specified DID (`dest`), then it can be considered as the creation of a new DID.

If there is a DID transaction with the specified DID (`dest`),  then this is update of existing DID.
In this case, we can specify only the values we would like to override. All unspecified values remain the same.
So, if key rotation needs to be performed, the owner of the DID needs to send a NYM request with
`dest`, `verkey` only. `alias` will stay the same.

<b>Note:</b>
Fields `dest` and `owner` should have the same value.


*Request Example*:
```
CreateDidRequest
{
    "data": {
              "alias": "Alice did",
              "dest": "GEzcdDLhCpGCYRHW82kjHd",
              "verkey": "~HmUWn928bnFT6Ephf65YXv"
             }
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

*Reply Example*:
```
CreateDidResponce
{}
```

DID transaction format:
```
Did
{
    "alias": "Alice did",
    "dest": "GEzcdDLhCpGCYRHW82kjHd",
    "verkey": "~HmUWn928bnFT6Ephf65YXv"
}
```

### ATTRIB

Adds or updates an attribute to a DID record.

- `dest` (base58-encoded string):

  Target DID as base58-encoded string for 16 or 32 byte DID value.

- `raw` (json; mutually exclusive with `hash` and `enc`):

  Raw data is represented as json, where the key is attribute name and value is attribute value.


*Request Example*:
```
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

*Reply Example*:
```
CreateAttribResponse
{}
```

### SCHEMA
Adds Claim's schema.

It's not possible to update an existing Schema.
So, if the Schema needs to be evolved, a new Schema with a new version or name needs to be created.

- `data`:

  Dictionary with Schema's data:

    - `attr_names`: array of attribute name strings (125 attributes maximum)
    - `name`: Schema's name string
    - `version`: Schema's version string

*Request Example*:
```
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

*Reply Example*:
```
{}
```

### CLAIM_DEF
Adds a claim definition (in particular, public key), that Issuer creates and publishes for a particular Claim Schema.

It's not possible to update `data` in existing Claim Def.
So, if a Claim Def needs to be evolved (for example, a key needs to be rotated), then
a new Claim Def needs to be created by a new Issuer DID (`owner`).


- `claim_def` (dict):

  Dictionary with Claim Definition's data:

    - `primary` (dict): primary claim public key
    - `revocation` (dict): revocation claim public key

- `ref` (string):

  Sequence number of a Schema transaction the claim definition is created for.

- `signature_type` (string):

  Type of the claim definition (that is claim signature). `CL` (Camenisch-Lysyanskaya) is the only supported type now.

- `tag` (string, optional):

  A unique tag to have multiple public keys for the same Schema and type issued by the same DID.
  A default tag `tag` will be used if not specified.

*Request Example*:
```
{
    "data": {
        "signature_type": "CL",
        "ref": 10,
        "tag": "some_tag",    
        "claim_def": {
            "primary": ....,
            "revocation": ....
        }
    },
    
    "owner": "L5AD5g65TDQr1PPHHRoiGf",
    "signature": "5ZTp9g4SP6t73rH2s8zgmtqdXyTuSMWwkLvfV1FD6ddHCpwTY5SAsp8YmLWnTgDnPXfJue3vJBWjy89bSHvyMSdS",
    "metadata": {}
}
```

*Reply Example*:
```
{}
```

### REVOC_REG_DEF
Adds a Revocation Registry Definition, that Issuer creates and publishes for a particular Claim Definition.
It contains public keys, maximum number of credentials the registry may contain, reference to the Claim Def, plus some revocation registry specific data.

- `value` (dict):

  Dictionary with revocation registry definition's data:

    - `max_cred_num` (integer): a maximum number of credentials the Revocation Registry can handle
    - `tails_hash` (string): tails' file digest
    - `tails_location` (string): tails' file location (URL)
    - `issuance_type` (string enum): defines credentials revocation strategy. Can have the following values:
        - `ISSUANCE_BY_DEFAULT`: all credentials are assumed to be issued initially, so that Revocation Registry needs to be updated (REVOC_REG_ENTRY txn sent) only when revoking. Revocation Registry stores only revoked credentials indices in this case. Recommended to use if expected number of revocation actions is less than expected number of issuance actions.
        - `ISSUANCE_ON_DEMAND`: no credentials are issued initially, so that Revocation Registry needs to be updated (REVOC_REG_ENTRY txn sent) on every issuance and revocation. Revocation Registry stores only issued credentials indices in this case. Recommended to use if expected number of issuance actions is less than expected number of revocation actions.
    - `public_keys` (dict): Revocation Registry's public key

- `id` (string): Revocation Registry Definition's unique identifier (a key from state trie is currently used)
- `cred_def_id` (string): The corresponding Credential Definition's unique identifier (a key from state trie is currently used)
- `revoc_def_type` (string enum): Revocation Type. `CL_ACCUM` (Camenisch-Lysyanskaya Accumulator) is the only supported type now.
- `tag` (string): A unique tag to have multiple Revocation Registry Definitions for the same Credential Definition and type issued by the same DID.

*Request Example*:
```
{
    "data": {
        "id": "L5AD5g65TDQr1PPHHRoiGf:3:FC4aWomrA13YyvYC1Mxw7:3:CL:14:some_tag:CL_ACCUM:tag1",
        "cred_def_dd": "FC4aWomrA13YyvYC1Mxw7:3:CL:14:some_tag"
        "revoc_def_type": "CL_ACCUM",
        "tag": "tag1",
        "value": {
            "max_cred_num": 1000000,
            "tails_hash": "6619ad3cf7e02fc29931a5cdc7bb70ba4b9283bda3badae297",
            "tails_location": "http://tails.location.com",
            "issuance_type": "ISSUANCE_BY_DEFAULT",
            "public_keys": {},
        },
    },
    
    "owner": "L5AD5g65TDQr1PPHHRoiGf",
    "signature": "5ZTp9g4SP6t73rH2s8zgmtqdXyTuSMWwkLvfV1FD6ddHCpwTY5SAsp8YmLWnTgDnPXfJue3vJBWjy89bSHvyMSdS",
    "metadata": {}
}
```

*Reply Example*:
```
{}
```

### REVOC_REG_ENTRY
The RevocReg entry containing the new accumulator value and issued/revoked indices. This is just a delta of indices, not the whole list. So, it can be sent each time a new claim is issued/revoked.

- `value` (dict):

  Dictionary with revocation registry's data:

    - `accum` (string): the current accumulator value
    - `prev_accum` (string): the previous accumulator value; it's compared with the current value, and txn is rejected if they don't match; it's needed to avoid dirty writes and updates of accumulator.
    - `issued` (list of integers): an array of issued indices (may be absent/empty if the type is ISSUANCE_BY_DEFAULT); this is delta; will be accumulated in state.
    - `revoked` (list of integers):  an array of revoked indices (delta; will be accumulated in state)

- `revoc_reg_def_id` (string): The corresponding Revocation Registry Definition's unique identifier (a key from state trie is currently used)
- `revoc_def_type` (string enum): Revocation Type. `CL_ACCUM` (Camenisch-Lysyanskaya Accumulator) is the only supported type now.

*Request Example*:
```
{
    "data": {
            "revoc_reg_def_id": "L5AD5g65TDQr1PPHHRoiGf:3:FC4aWomrA13YyvYC1Mxw7:3:CL:14:some_tag:CL_ACCUM:tag1"
            "revoc_def_type": "CL_ACCUM",
            "value": {
                "accum": "accum_value",
                "prev_accum": "prev_acuum_value",
                "issued": [],
                "revoked": [10, 36, 3478],
            },
    },
    "owner": "L5AD5g65TDQr1PPHHRoiGf",
    "signature": "5ZTp9g4SP6t73rH2s8zgmtqdXyTuSMWwkLvfV1FD6ddHCpwTY5SAsp8YmLWnTgDnPXfJue3vJBWjy89bSHvyMSdS",
    "metadata": {}
}
```

*Reply Example*:
```
{}
```

## References

- [Hyperledger Indy Identity transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)
- [W3 DID Spec](https://www.w3.org/TR/did-core/)