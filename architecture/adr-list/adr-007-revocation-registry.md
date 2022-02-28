# ADR 007: Revocation registry

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar |
| **ADR Stage** | DRAFT |
| **Implementation Status** | Not Implemented |
| **Start Date** | 2021-09-10 |

## Summary

Issued credentials need to be revocable by their issuers. Revocation needs to be straightforward and fast. Testing of revocation needs to preserve privacy (be non-correlating), and it should be possible to do without contacting the issuer.

## Context

This has obvious use cases for professional credentials being revoked for fraud or misconduct, e.g., a driver’s license could be revoked for criminal activity. However, it’s also important if a credential gets issued in error (e.g., has a typo in it that misidentifies the subject). The latter case is important even for immutable and permanent credentials such as a birth certificate.

In addition, it seems likely that the data inside credentials will change over time (e.g., a person’s mailing address or phone number updates). This is likely to be quite common, revocation can be used to guarantee currency of credential data when it happens. In other words, revocation may be used to force updated data, not just to revoke authorization.

## Decision

### REVOC\_REG\_DEF

Adds a Revocation Registry Definition, that Issuer creates and publishes for a particular Credential Definition. It contains public keys, maximum number of credentials the registry may contain, reference to the Credential Definition, plus some revocation registry specific data.

* **`value` (dict):**

  Dictionary with Revocation Registry Definition's data:

  * **`max_cred_num`** (integer): The maximum number of credentials the Revocation Registry can handle
  * **`tails_hash`** (string): Tails file digest
  * **`tails_location`** (string): Tails file location (URL)
  * **`issuance_type`** (string enum): Defines credential revocation strategy. Can have the following values:
    * `ISSUANCE_BY_DEFAULT`: All credentials are assumed to be issued and active initially, so that Revocation Registry needs to be updated (`REVOC_REG_ENTRY` transaction sent) only when revoking. Revocation Registry stores only revoked credentials indices in this case. Recommended to use if expected number of revocation actions is less than expected number of issuance actions.
    * `ISSUANCE_ON_DEMAND`: No credentials are issued initially, so that Revocation Registry needs to be updated (`REVOC_REG_ENTRY` transaction sent) on every issuance and revocation. Revocation Registry stores only issued credentials indices in this case. Recommended to use if expected number of issuance actions is less than expected number of revocation actions.
  * **`public_keys`** (dict): Revocation Registry's public key

* **`id`** (string): Revocation Registry Definition's unique identifier (a key from state trie is currently used) `owner:cred_def_id:revoc_def_type:tag`
* **`cred_def_id`** (string): The corresponding Credential Definition's unique identifier (a key from state tree is currently used)
* **`revoc_def_type`** (string enum): Revocation Type. `CL_ACCUM` (Camenisch-Lysyanskaya Accumulator) is the only supported type now.
* **`tag`** (string): A unique tag to have multiple Revocation Registry Definitions for the same Credential Definition and type issued by the same DID.

**Note**: REVOC\_REG\_DEF **can** be updated.

#### State format

`(owner, cred_def_id, revoc_def_type, tag) -> {data, tx_hash, tx_timestamp }`

_Request Example_:

```text
{
    "data": {
        "id": "L5AD5g65TDQr1PPHHRoiGf:3:FC4aWomrA13YyvYC1Mxw7:3:CL:14:some_tag:CL_ACCUM:tag1",
        "cred_def_id": "FC4aWomrA13YyvYC1Mxw7:CL:5ZTp9g4SP6t73rH2s8z:some_tag"
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

_Reply Example_:

```text
{}
```

### REVOC\_REG\_ENTRY

The Revocation Registry Entry contains the new accumulator value and issued/revoked indices. This is just a delta of indices, not the whole list. It can be sent each time a new credential is issued/revoked.

* **`value`** (dict):

  Dictionary with revocation registry's data:

  * **`accum`** (string): The current accumulator value
  * **`prev_accum`** (string): The previous accumulator value. It is compared with the current value, and transaction is rejected if they don't match. This is necessary to avoid dirty writes and updates of accumulator.
  * **`issued`** (list of integers): An array of issued indices (may be absent/empty if the type is `ISSUANCE_BY_DEFAULT`). This is delta, and will be accumulated in state.
  * **`revoked`** (list of integers):  An array of revoked indices. This is delta; will be accumulated in state)

* **`revoc_reg_def_id`** (string): The corresponding Revocation Registry Definition's unique identifier (a key from state trie is currently used)
* **`revoc_def_type`** (string enum): Revocation Type. `CL_ACCUM` (Camenisch-Lysyanskaya Accumulator) is the only supported type now.

**Note**: REVOC\_REG\_ENTRY **can** be updated.

#### State format

1. `MARKER_REVOC_REG_ENTRY_ACCUM:revoc_reg_def_id -> {data, tx_hash, tx_timestamp }`
2. `MARKER_REVOC_REG_ENTRY:revoc_reg_def_id -> {data, tx_hash, tx_timestamp }`

_Request Example_:

```text
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

_Reply Example_:

```text
{}
```

## References

* [Hyperledger Indy Credential Revocation Hype](https://hyperledger-indy.readthedocs.io/projects/hipe/en/latest/text/0011-cred-revocation/README.html)

