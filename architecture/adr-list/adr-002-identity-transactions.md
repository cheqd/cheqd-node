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

Our proposal is to change the term `NYM` in transactions to `DID`, which would make understanding the context of a transaction easier to understand. This change will bring transactions better in-line with World Wide Web Consortium \(W3C\) terminology.

#### Remove `role` field from DID transaction

Hyperledger Indy is a public-permissioned distributed ledger. As `cheqd-node` is based on a public-permissionless network based on the [Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk), the need for having a `role` type is not necessary.

_**Note**: Hyperledger Indy also contains other transaction types beyond the ones listed above, but these are currently not in scope for implementation in `cheqd-node`. They will be considered for inclusion later in the product roadmap._

## Decision

### cheqd DID Method

The `did:cheqd` method DID has four components that are concatenated to make a W3C DID specification conformant identifier. The components are:

- **DID**: the hardcoded string `did:` to indicate the identifier is a DID
- **`cheqd` DID method**: the hardcoded string `cheqd:` indicating that the identifier uses the `cheqd` DID method specification.
- **Namespace**: A string that identifies the name of the primary cheqd ledger ("mainnet"), followed by a `:`. The namespace string may optionally have a secondary ledger name prefixed by a `:` following the primary name. If there is no secondary ledger element, the DID resides on the primary ledger ("mainnet"), else it resides on the secondary ledger. By convention, the primary is a production ledger while the secondary ledgers are non-production ledgers (e.g. testnet) associated with the primary ledger.
- **Namespace Identifier**: A self-certified identifier unique to the given cheqd DID ledger namespace. To be self-certifying, the identifier must be derived from the initial `verkey` associated with the identifier.

The components are assembled as follows:

`did:cheqd:<namespace>:<namespace identifier>`

Some examples of `did:cheqd` method identifiers are:

- A DID written to the cheqd mainnet ledger:
  `did:cheqd:mainnet:7Tqg6BwSSWapxgUDm9KKgg`
- A DID written to the cheqd testnet ledger:
  `did:cheqd:testnet:6cgbu8ZPoWTnR5Rv5JcSMB`


### General structure of transaction requests

All identity requests will have the following format:

```json
{
    "data": { <request data for writing a transaction to the ledger> },
    "creators": [<identifier>, ...],
    "signatures": [
        <public_key>: <signature>,
      ],
    "metadata": {}
}
```

* **`data`**: Data requested to be written to the ledger, specific for each request type.
* **`creators`**: Creators identifiers (DID) list for this entity. There should be a new DIDs or an existing DIDs, for existing entities.
* **`signatures`**: `data` should be signed by all `creators` private key. This field contains a dict there creator's public key is a key, and the signature is a value.
* **`metadata`**: Dictionary with additional metadata fields. Empty for now. This fields provides extensibility in the future, e.g., it can contain `protocolVersion` or other relevant metadata associated with a request.

## List of transactions and details

### `DID` transactions

[Decentralized Identifiers \(DIDs\) are a W3C specification](https://www.w3.org/TR/did-core/) for identifiers that enable verifiable, decentralized digital identity.

DIDDoc format conforms to [DIDDoc spec]().
The request can be used for creation of new DIDDoc, setting, and rotation of verification key.

#### DIDDoc

1. **`id`**: Target DID as base58-encoded string for 16 or 32 byte DID value. In the ledger, we store only an identifier without specifying the method and namespace.
2. **`controller`** (optional): A list of base58-encoded identifier strings.
3. **`verificationMethod`** (optional): A list of Verification Methods
4. **`authentication`** (optional): A list of Verification Methods or strings with key aliases
5. **`assertionMethod`** (optional): A list of Verification Methods or strings with key aliases
6. **`capabilityInvocation`** (optional): A list of Verification Methods or strings with key aliases
7. **`capabilityDelegation`** (optional): A list of Verification Methods or strings with key aliases
8. **`service`** (optional): A set of Service Endpoint maps
9.  **`@context`** (optional): A list of strings

**Example:**
```json
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/suites/ed25519-2020/v1"
  ],
  "id": "N22KY2Dyvmuu2PyyqSFKue",
  "authentication": [
    {
      "id": "N22KY2Dyvmuu2PyyqSFKue#z6MkecaLyHuYWkayBDLw5ihndj3T1m6zKTGqau3A51G7RBf3",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf"
    }
  ],
  "capabilityInvocation": [
    {
      "id": "N22KY2Dyvmuu2PyyqSFKue#z6MkhdmzFu659ZJ4XKj31vtEDmjvsi5yDZG5L7Caz63oP39k",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "z4BWwfeqdp1obQptLLMvPNgBw48p7og1ie6Hf9p5nTpNN"
    }
  ],
  "capabilityDelegation": [
    {
      "id": "N22KY2Dyvmuu2PyyqSFKue#z6Mkw94ByR26zMSkNdCUi6FNRsWnc2DFEeDXyBGJ5KTzSWyi",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "zHgo9PAmfeoxHG8Mn2XHXamxnnSwPpkyBHAMNF3VyXJCL"
    }
  ],
  "assertionMethod": [
    {
      "id": "N22KY2Dyvmuu2PyyqSFKue#z6MkiukuAuQAE8ozxvmahnQGzApvtW7KT5XXKfojjwbdEomY",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "z5TVraf9itbKXrRvt2DSS95Gw4vqU3CHAdetoufdcKazA"
    }
  ]
 }
```


#### Verification Method

1. **`id`** (string): A string with format `<DIDDoc-id>#<key-alias`
2. **`controller`**: A list of base58-encoded identifier strings.
3. **`type`** (string)
4. **`publicKeyJwk`** (`map[string,string]`, optional): A map representing a JSON Web Key that conforms to [RFC7517](https://tools.ietf.org/html/rfc7517). See definition of `publicKeyJwk` for additional constraints.
5. **`publicKeyMultibase`** (optional): A base58-encoded string that conforms to a [MULTIBASE](https://datatracker.ietf.org/doc/html/draft-multiformats-multibase-03) encoded public key.

**Example:**
```json
{
  "id": "N22KY2Dyvmuu2PyyqSFKue#key-0",
  "type": "JsonWebKey2020",
  "controller": "N22KY2Dyvmuu2PyyqSFKue",
  "publicKeyJwk": {
    "kty": "OKP", // external (property name)
    "crv": "Ed25519", // external (property name)
    "x": "VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ" // external (property name)
}
```


#### Service

1. **`id`** (string): The value of the id property MUST be a URI conforming to [RFC3986](https://www.rfc-editor.org/rfc/rfc3986). A conforming producer MUST NOT produce multiple service entries with the same ID. A conforming consumer MUST produce an error if it detects multiple service entries with the same ID. It has a follow format: `<DIDDoc-id>#<service-alias>`

2. **`type`** (string): The service type and its associated properties SHOULD be registered in the DID Specification Registries [DID-SPEC-REGISTRIES](https://www.w3.org/TR/did-spec-registries/).

3. **`serviceEndpoint`** (strings): A string that conforms to the rules of [RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs, a map, or a set composed of a one or more strings that conform to the rules of [RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs and/or maps.

**Example:**

 ```json
"service": [{
  "id":"N22KY2Dyvmuu2PyyqSFKue#linked-domain",
  "type": "LinkedDomains",
  "serviceEndpoint": "https://bar.example.com"
}]
```

#### **Update DID**

If there is no DID transaction with the specified DID \(`dest`\), it is considered as a creation request for a new DID.

If there is a DID transaction with the specified DID \(`dest`\), then this is update of existing DID. In this case, we can specify only the values we would like to override. All unspecified values remain the same. E.g., if a key rotation needs to be performed, the owner of the DID needs to send a DID transaction request with `dest`, `verkey` only. `alias` will stay the same.

**Note:** Fields `dest` and `owner` should have the same value.

#### State format

`id -> {encode(data, creators), tx_hash, tx_timestamp }`

### SCHEMA

This transaction is used to create a Schema associated with credentials.

It is not possible to update an existing Schema, to ensure the original schema used to issue any credentials in the past are always available.

If a Schema evolves, a new schema with a new version or name needs to be created.

* **`data`**: Dictionary with Schema's data:
    * **`id`**: DID as base58-encoded string for 16 or 32 byte DID value.
    * **`attr_names`**: Array of attribute name strings (125 attributes maximum)
    * **`name`**: Schema's name string
    * **`version`**: Schema's version string

#### SCHEMA transaction format

```json
{
  "id": "N22KY2Dyvmuu2PyyqSFKue",
  "version": "1.0",
  "name": "Degree",
  "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
}
```

**Note:** SCHEMA **cannot** be updated

#### State format

`id -> {encode(data, creators), tx_hash, tx_timestamp }`

### CRED_DEF

Adds a Credential Definition (in particular, public key), which is created by an Issuer and published for a particular Credential Schema.

It is not possible to update `data` in existing Credential Definitions. If a Credential Definition needs to be evolved \(for example, a key needs to be rotated\), then a new Credential Definition needs to be created by a new Issuer DID \(`owner`\).

* **`id`**: DID as base58-encoded string for 16 or 32 byte DID value.
* **`value`** \(dict\): Dictionary with Credential Definition's data if `signature_type` is `CL`:
  * **`primary`** (dict): Primary credential public key
  * **`revocation`** (dict, optional): Revocation credential public key
* **`ref`** (string): Hash of a Schema transaction the credential definition is created for.
* **`signature_type`** (string): Type of the credential definition \(that is credential signature\). `CL` \(Camenisch-Lysyanskaya\) is the only supported type now. Other signature types are being explored for future releases.
* **`tag`** (string, optional): A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

#### CRED_DEF transaction format:

```json
{
  "id": "N22KY2Dyvmuu2PyyqSFKue",
  "signature_type": "CL",
  "schema_id": "5ZTp9g4SP6t73rH2s8zgmtqdXyT",
  "tag": "some_tag",    
  "value": {
      "primary": ....,
      "revocation": ....
  }
}
```

**Note**: CRED_DEF **cannot** be updated.

#### State format

`id -> {encode(data, creators), tx_hash, tx_timestamp }`

## References

* [Hyperledger Indy Identity transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)
* [W3 DID Spec](https://www.w3.org/TR/did-core/)

