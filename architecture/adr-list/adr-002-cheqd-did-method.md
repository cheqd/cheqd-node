# ADR 002: cheqd DID method, identity entities, and transactions

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar, Brent Zundel, Ankur Banerjee |
| **ADR Stage** | PROPOSED |
| **Implementation Status** | Implementation in progress |
| **Start Date** | 2021-09-23 |

## Summary

This ADR defines the cheqd DID method and describes the identity entities, queries, and transaction types for the cheqd network.

[Decentralized identifiers](https://www.w3.org/TR/did-core) (DIDs) are a type of identifier that enables verifiable, decentralized digital identity. A DID refers to any subject (for example, a person, organization, thing, data model, abstract entity, and so on) as determined by the controller of the DID.

The identity entities and transactions for the cheqd network are designed to support usage scenarios and functionality currently supported by [Hyperledger Indy](https://github.com/hyperledger/indy-node).

## Context

### Rationale for baselining against Hyperledger Indy

Hyperledger Indy is a verifiable data registry (VDR) built for DIDs with a strong focus on privacy-preserving techniques. It is one of the most widely-adopted self-sovereign identity (SSI) blockchain ledgers; most notably by the [Sovrin network](https://sovrin.org/overview/).



However, as the core of Hyperledger Indy's architecture was designed before the [W3C DID specification](https://www.w3.org/TR/did-core/) started to be defined, the [Indy DID Method](https://hyperledger.github.io/indy-did-method/) (`did:indy`) has aspects that are not fully-compliant with latest specifications.

Hyperledger Indy contains the following
[identity domain transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md):

1. `NYM`
2. `ATTRIB`
3. `SCHEMA`
4. `CRED_DEF`
5. `REVOC_REG_DEF`
6. `REVOC_REG_ENTRY`

Our aim is to bring the functionality enabled by these transactions into
`cheqd-node` to allow the use cases of existing SSI projects that work with
Hyperledger Indy to be supported by the cheqd network.

We define transaction that differ from the ones listed above, but
which enable equivalent support for privacy-respecting SSI use cases. The
differences stem primarily from our conformance to
[Decentralized Identifiers v1.0](https://www.w3.org/TR/did-core/) and the
capabilities of the underlying
[Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk) that we use.

## Decision

The cheqd DID Method will conform to [Decentralized Identifiers v1.0](https://github.com/w3c/did-core) with the goal of maximizing interoperability with other compatible tools and projects.

### Syntax

As with all DIDs, cheqd identifiers begin with the string `did:`.

#### DID Method Name
The `method-name` is the string `cheqd:`.

#### cheqd DID Method Specific Identifier

The cheqd DID `method-specific-id` is made up of two components, a `namespace`
and a `namespace-id`:

- **namespace**: A string that identifies the name of the primary cheqd ledger
(e.g., "mainnet"), followed by a `:`. The namespace string may optionally have a
secondary ledger name prefixed by a `:` following the primary name. If there is
no secondary ledger element, the DID resides on the primary ledger ("mainnet"),
else it resides on the secondary ledger. By convention, the primary is a
production ledger while the secondary ledgers are non-production ledgers (e.g.
testnet) associated with the primary ledger.
- **Namespace Identifier**: A self-certified identifier unique to the given
cheqd DID ledger namespace.

The components are assembled as follows:

`did:cheqd:<namespace>:<namespace identifier>`

Some examples of `did:cheqd` method identifiers are:

- A DID written to the cheqd mainnet ledger: `did:cheqd:mainnet:7Tqg6BwSSWapxgUDm9KKgg`
- A DID written to the cheqd testnet ledger: `did:cheqd:testnet:6cgbu8ZPoWTnR5Rv5JcSMB`

### Operations

#### General structure of transaction requests

All identity requests will have the following format:

```jsonc
{
  "data": { "<request data for writing a transaction to the ledger>" },
  "signatures": {
      "verification method id": "signature"
      // Multiple verification methods and corresponding signatures can be added here
    },
  "metadata": {}
}
```

- **`data`**: Data requested to be written to the ledger, specific for each
request type.
- **`signatures`**: `data`and `metadata` should be signed by all `controller`
private keys. This field contains a dict there key's id from
`DIDDoc.authentication` is a key, and the signature is a value. The `signatures`
must contains signatures from all controllers. And every controller should sign
all fields excluding `signatures` using at least one key from
`DIDDoc.authentication`.
- **`requestId`**: String with unique identifier. Unix timestamp is recommended.
Needed for a reply protection.
- **`metadata`**: Dictionary with additional metadata fields. Empty for now.
This fields provides extensibility in the future, e.g., it can contain
`protocolVersion` or other relevant metadata associated with a request.

  
#### `DID` transactions

[Decentralized Identifiers \(DIDs\) are a W3C specification](https://www.w3.org/TR/did-core/)
for identifiers that enable verifiable, decentralized digital identity.

DIDDoc format conforms to
[DIDDoc spec](https://www.w3.org/TR/did-core/#representations).
The request can be used for creation of new DIDDoc, setting, and rotation of
verification key.

##### DIDDoc

1. **`id`**: Target DID as base58-encoded string for 16 or 32 byte DID value
with Cheqd DID Method prefix `did:cheqd:<namespace>:<namespace identifier>:`.
2. **`controller`** (optional): A list of fully qualified DID strings or one
string. Contains one or more DIDs who can update this DIDdoc. All DIDs must
exist.
3. **`verificationMethod`** (optional): A list of Verification Methods
4. **`authentication`** (optional): A list of Verification Methods or strings
with key aliases
5. **`assertionMethod`** (optional): A list of Verification Methods or strings
with key aliases
6. **`capabilityInvocation`** (optional): A list of Verification Methods or
strings with key aliases
7. **`capabilityDelegation`** (optional): A list of Verification Methods or
strings with key aliases
8. **`keyAgreement`** (optional): A list of Verification Methods or strings
with key aliases
9. **`service`** (optional): A set of Service Endpoint maps
10. **`alsoKnownAs`** (optional): A list of strings. A DID subject can have
multiple identifiers for different purposes, or at different times. The
assertion that two or more DIDs refer to the same DID subject can be made using
the `alsoKnownAs` property.
11. **`@context`** (optional): A list of strings with links or JSONs for
describing specifications that this DID Document is following to.

For Example:

```jsonc
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/suites/ed25519-2020/v1"
  ],
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "verificationMethod": [
    {
      "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#authKey1",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf"
    },
    {
      "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#capabilityInvocationKey",
      "type": "Ed25519VerificationKey2020", // external (property value)
      "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
      "publicKeyMultibase": "z4BWwfeqdp1obQptLLMvPNgBw48p7og1ie6Hf9p5nTpNN"
    }
  ],
  "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#authKey1"],
  "capabilityInvocation": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#capabilityInvocationKey"],
}
```

##### Verification Method

1. **`id`** (string): A string with format `<DIDDoc-id>#<key-alias>`
2. **`controller`**: A list of fully qualified DID strings or one string. All
DIDs must exist.
3. **`type`** (string)
4. **`publicKeyJwk`** (`map[string,string]`, optional): A map representing a
JSON Web Key that conforms to [RFC7517](https://tools.ietf.org/html/rfc7517).
See definition of `publicKeyJwk` for additional constraints.
5. **`publicKeyMultibase`** (optional): A base58-encoded string that conforms to
a [MULTIBASE](https://datatracker.ietf.org/doc/html/draft-multiformats-multibase-03)
encoded public key.
**Note**: Verification Method cannot contain both `publicKeyJwk` and
`publicKeyMultibase` but must contain at least one of them.

For Example:

```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#key-0",
  "type": "JsonWebKey2020",
  "controller": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "publicKeyJwk": {
    "kty": "OKP",
    // external (property name)
    "crv": "Ed25519",
    // external (property name)
    "x": "VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ"
    // external (property name)
  }
}
```

##### Service

1. **`id`** (string): The value of the id property MUST be a URI conforming to
[RFC3986](https://www.rfc-editor.org/rfc/rfc3986). A conforming producer MUST
NOT produce multiple service entries with the same ID. A conforming consumer
MUST produce an error if it detects multiple service entries with the same ID.
It has a follow formats: `<DIDDoc-id>#<service-alias>` or `#<service-alias>`.
2. **`type`** (string): The service type and its associated properties SHOULD be
registered in the DID Specification Registries
[DID-SPEC-REGISTRIES](https://www.w3.org/TR/did-spec-registries/)
3. **`serviceEndpoint`** (strings): A string that conforms to the rules of
[RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs, a map, or a set
composed of a one or more strings that conform to the rules of
[RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs and/or maps.

For Example:

```jsonc
{
  "id":"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#linked-domain",
  "type": "LinkedDomains",
  "serviceEndpoint": "https://bar.example.com"
}
```

##### Create `DID`

If there is no DID entry on the ledger with the specified DID (`DID.id`), it is
considered as a creation request for a new DID.

**Note**: The field `signatures`(from `WriteRequest`) must contain signatures
from all new controllers.

##### Update `DID`

If there is a DID entry on the ledger with the specified DID (`DID.id`), then
this considered a request for updating an existing DID.

For updating `versionId` from `UpdateDIDRequest` should be filled by a
transaction hash of the previous DIDDoc version. It is needed for a replay
  protection.

**Note**: The field `signatures`(from `WriteRequest`) must contain signatures
from all old controllers and all new controllers.

##### Resolve DID

DIDDoc State format:

`"diddoc:<id>" -> {DIDDoc, DidDocumentMetadata, txHash, txTimestamp }`

`didDocumentMetadata` is created by the node after transaction ordering and
before adding it to a State.

DID Document Metadata:

1. **`created`** (string): Formatted as an XML Datetime normalized to UTC
00:00:00 and without sub-second decimal precision. For example:
2020-12-20T19:17:47Z.
2. **`updated`** (string): The value of the property MUST follow the same
formatting rules as the created property. The `updated` field is null if an
Update operation has never been performed on the DID document. If an updated
property exists, it can be the same value as the created property when the
difference between the two timestamps is less than one second.
3. **`deactivated`** (strings): If DID has been deactivated, DID document
metadata MUST include this property with the boolean value true. By default
`false`.
4. **`versionId`** (strings): Contains transaction hash of the current DIDDoc
version.

For Example:
```jsonc
{
  "created": "2020-12-20T19:17:47Z",
  "updated": "2020-12-20T19:19:47Z",
  "deactivated": false,
  "versionId": "N22KY2Dyvmuu2PyyqSFKueN22KY2Dyvmuu2PyyqSFKue",
}
```

TODO: describe the resolution process

##### Deactivate DID

TODO: describe the deactivation process

#### `SCHEMA`

This transaction is used to create a Schema associated with credentials.

It is not possible to update an existing Schema, to ensure the original schema
used to issue any credentials in the past are always available.

If a Schema evolves, a new schema with a new version or name needs to be created.

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with Cheqd
DID Method prefix `did:cheqd:<namespace>:<namespace identifier>:` and a resource
type at the end.
- **`type`**: String with a schema type. Now only `CL-Schema` is supported.
- **`attrNames`**: Array of attribute name strings (125 attributes maximum)
- **`name`**: Schema's name string
- **`version`**: Schema's version string
- **`controller`**: DIDs list of strings or only one string of a schema
controller(s). All DIDs must exist.

`SCHEMA` entity transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue/schema",
  "type": "CL-Schema",
  "controller": "did:cheqd:mainnet:IK22KY2Dyvmuu2PyyqSFKu",  // Schema Issuer DID
  "version": "1.0",
  "name": "Degree",
  "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
}
```

Don't store Schema DIDDoc in the State.

Schema URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue`

Schema Entity URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue/schema` 

`SCHEMA` DID Document transaction format: 
```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "controller": "did:cheqd:mainnet:IK22KY2Dyvmuu2PyyqSFKu", // Schema Issuer DID
  "service":[
    {
      "id": "cheqd-schema1",
      "type": "CL-Schema",
      "serviceEndpoint": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?resource=true"
    }
  ]
}
```

**Note**: `SCHEMA` **cannot** be updated

**`SCHEMA` State format:**

- `"schema:<id>" -> {SchemaEntity, txHash, txTimestamp}`

`id` example: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue`

[Link to DidDocumentMetadata description](#resolve-did)

#### `CRED_DEF`

Adds a Credential Definition (in particular, public key), which is created by an
Issuer and published for a particular Credential Schema.

It is not possible to update Credential Definitions. If a Credential Definition
needs to be evolved (for example, a key needs to be rotated), then a new
Credential Definition needs to be created for a new Issuer DIDdoc.
Credential Definitions is added to the ledger in as verification method for
Issuer DIDDoc

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with Cheqd
DID Method prefix `did:cheqd:<namespace>:<namespace identifier>:` and a resource
type at the end.
- **`value`** (dict): Dictionary with Credential Definition's data if
`signature_type` is `CL`:
  - **`primary`** (dict): Primary credential public key
  - **`revocation`** (dict, optional): Revocation credential public key
- **`schemaId`** (string): `id` of a Schema the credential definition is created
for.
- **`signatureType`** (string): Type of the credential definition (that is
credential signature). `CL-Sig-Cred_def` (Camenisch-Lysyanskaya) is the only
supported type now. Other signature types are being explored for future releases.
- **`tag`** (string, optional): A unique tag to have multiple public keys for
the same Schema and type issued by the same DID. A default tag `tag` will be
used if not specified.
- **`controller`**: DIDs list of strings or only one string of a credential
definition controller(s). All DIDs must exist.

`CRED_DEF` entity transaction format:

```jsonc
{
    "id": "<cred_def_url>",
    "type": "CL-CredDef",
    "controller": "did:cheqd:mainnet:123456789abcdefghi",
    "schemaId": "did:cheqd:mainnet:5ZTp9g4SP6t73rH2s8zgmtqdXyT/schema",
    "tag": "some_tag",
    "value": {
      "primary": "...",
      "revocation": "..."
    }
}
```

Don't store Schema DIDDoc in the State.

CredDef URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue`

CredDef Entity URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue/credDef`

`CRED_DEF` DID Document transaction format:
```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "controller": "did:cheqd:mainnet:IK22KY2Dyvmuu2PyyqSFKu", // CredDef Issuer DID
  "service":[
    {
      "id": "cheqd-cred-def",
      "type": "CL-CredDef",
      "serviceEndpoint": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue/credDef"
    }
  ]
}
```

`CRED_DEF` state format:

`"credDef:<id>" -> {CredDefEntity, txHash, txTimestamp}`

**Note**: `CRED_DEF` **cannot** be updated.

### Security Considerations

TODO: add security considerations

### Privacy Considerations

TODO: add privacy considerations

### Rationale and Alternatives

#### Changes from Indy 

##### Rename `NYM` transactions to `DID` transactions

[**NYM** is the term used by Hyperledger Indy](https://hyperledger-indy.readthedocs.io/projects/node/en/latest/transactions.html#nym) for [Decentralized Identifiers \(DIDs\)](https://www.w3.org/TR/did-core/)
that are created on ledger. A DID is typically the identifier that is associated
with a specific organisation issuing/managing SSI credentials.

Our proposal is to change the term `NYM` in transactions to `DID`, which would
make understanding the context of a transaction easier to understand. This
change will bring transactions better in-line with World Wide Web Consortium
\(W3C\) terminology.

#### Remove `role` field from DID transaction

Hyperledger Indy is a public-permissioned distributed ledger. As `cheqd-node`
is based on a public-permissionless network based on the
[Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk), the `role`
type is no longer necessary.

#### Dropping `ATTRIB` transactions

`ATTRIB` was originally used in Hyperledger Indy to add document content similar
to DID Documents (DIDDocs). The cheqd DID method replaces this by implementing
DIDDocs for most transaction types.


#### Schema options not used
##### Option 2

Schema URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue#<schema_entity_id>`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "schema":[
    {
      "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#schema1",
      "type": "CL-Schema",
      "controller": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
      "value": {
                "version": "1.0",
                "name": "Degree",
                "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
              },
    },
  ]
}
```

##### Option 3

Schema URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "schema": {
              "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
              "type": "CL-Schema",
              "controller": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
              "value": {
                "version": "1.0",
                "name": "Degree",
                "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
              },
            },
}
```

##### Option 4

Schema URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue#<schema_entity_id>`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
  "schema":[
              {
                "id": "cheqd-schema",
                "type": "CL-Schema",
                "schemaRef": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue?resource=true"
              }
          ]
}
```

`SCHEMA` State format:

- `"schema:<id>" -> {SchemaEntity, txHash, txTimestamp}`

`id` example: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue`

[Link to DidDocumentMetadata description](#resolve-did)

#### Cred Def options not used

##### Option 2

Store inside Issuer DID Document

CredDef URL: `did:cheqd:N22KY2Dyvmuu2PyyqSFKue#<cred_def_entity_id>`

`CRED_DEF` DID Document transaction format:
```jsonc
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/suites/jws-2020/v1",
    "https://w3id.org/security/suites/ed25519-2020/v1"
  ],
  "id": "did:cheqd:mainnet:123456789abcdefghi",
  "verificationMethod": [{
    "id": "passport-keys",
    "type": "CL-Sig-Cred_def",
    "controller": "did:cheqd:mainnet:123456789abcdefghi",
    "schemaId": "did:cheqd:mainnet:5ZTp9g4SP6t73rH2s8zgmtqdXyT",
    "tag": "some_tag",
    "value": {
      "primary": "...",
      "revocation": "..."
    }
  }]

}
```

`CRED_DEF` state format:

Stored inside [DIDDoc](#resolve-did)

## Consequences

### Backward Compatibility

- `cheqd-node` [release v0.1.17](https://github.com/cheqd/cheqd-node/releases/tag/v0.1.17) and earlier had a transaction type called `NYM` which would allow writing/reading a unique identifier on ledger. However, this `NYM` state was not fully defined as a DID method and did not contain DID Documents that resolved when the DID identifier was read. This `NYM` transaction type is deprecated and the data written to cheqd testnet with legacy states will not be retained.

### Positive

### Negative

### Neutral

## References

- [Hyperledger Indy Identity transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)
- [W3 DID Spec](https://www.w3.org/TR/did-core/)

## Unresolved questions

- Describe the resolution process for DID
- Describe the deactivation process
- Add security considerations 
- Add privacy considerations