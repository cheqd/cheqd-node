# ADR 008: cheqd DIDDoc resources: Schemas and Credential Definitions

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar, Brent Zundel, Ankur Banerjee |
| **ADR Stage** | DRAFT |
| **Implementation Status** | Draft |
| **Start Date** | 2021-09-23 |

## Summary

In this ADR we look only at CL (Camenisch-Lysyanskaya) schema and credential definition transactions that are needed for verifying issued Verifiable Credentials.

This ADR will define how Verifiable Credential schemas can be represented through the use of a DID URL, which when dereferenced, fetches the credential schemas a resource. 

The identity entities and transactions for the cheqd network are designed to support usage scenarios and functionality currently supported by [Hyperledger Indy](https://github.com/hyperledger/indy-node).


## Context

Hyperledger Indy is a verifiable data registry (VDR) built for DIDs with a strong focus on privacy-preserving techniques. It is one of the most widely-adopted SSI blockchain ledgers. Most notably, Indy is used by the [Sovrin Network](https://sovrin.org/overview/).

### Identity-domain transaction types in Hyperledger Indy

Our aim is to support the functionality enabled by [identity-domain transactions in by Hyperledger Indy](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md) into `cheqd-node`. This will partly enable the goal of allowing use cases of existing SSI networks on Hyperledger Indy to be supported by the cheqd network.

The following identity-domain transactions from Indy were considered:

1. `NYM`: Equivalent to "DIDs" on other networks
2. `ATTRIB`: Payload for DID Document generation
3. `SCHEMA`: Schema used by a credential
4. `CRED_DEF`: Credential definition by an issuer for a particular schema
5. `REVOC_REG_DEF`: Credential revocation registry definition
6. `REVOC_REG_ENTRY`: Credential revocation registry entry

Revocation registries for credentials are not covered under the scope of this ADR. This topic is discussed separately in [ADR 007: **Revocation registry**](adr-007-revocation-registry.md) as there is ongoing research by the cheqd project on how to improve the privacy and scalability of credential revocations.

### Resolving DID vs Dereferencing DID

Before diving into the specific architecture of cheqd's schemas, it is imporntant to understand the difference between resolving a DID and dereferencing a DID URL. 

When you resolve a DID, a DID Document is returned. For example, resolving: "did:cheqd:example1234" would return the full DID Document associated with the specific DID "did:cheqd:example1234". 

Dereferecing a DID URL is slightly different. When you dereference a DID URL, you are parsing the URL for specific actions, such as to take a certain path, highlight a specific fragment, or query a specific resource. 

For example, "did:cheqd:example1234?service=ExampleSchema" can be dereferenced. In this case it will query the service of type "ExampleSchema" within the DID Document and will return the resource specified at the Service Endpoint within the DID Document. It is this type of architecture that cheqd uses in this ADR to fetch schemas. 

## Decision

### Schema

The most important concept to understand for this architecture is that each Schema will have its own DID Document, which is able to be resolved or dereferenced.  

Different results will be returned for a resolved or dereferenced Schema DID URL, this will be explained below.

#### Creating a Schema

The first stage of using this Schema architecture is creating the Schema as a transaction on ledger. 

This transaction is used to create a Schema associated with credentials:

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with cheqd DID Method prefix `did:cheqd:<namespace>:` and a resource
type at the end.
- **`type`**: String with a schema type. Currently only `CL-Schema` is supported.
- **`attrNames`**: Array of attribute name strings (125 attributes maximum)
- **`name`**: Schema's name string
- **`version`**: Schema's version string
- **`controller`**: DIDs list of strings or only one string of a schema controller(s). All DIDs must exist.

An example of this schema transaction, written in JSON would be:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
  "type": "CL-Schema",
  "controller": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu",  // Schema Issuer DID
  "version": "1.0",
  "name": "Degree",
  "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
}
```
This creates an artifact, that when dereferenced will return the appropriate fields for the schema. 

#### Schema DID Document URL

This is an example of a Schema's DID Document:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue", // Schema's public DID
  "controller": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu", // Schema issuer's DID
  "service":[
    {
      "id": "cheqd-schema", 
      "type": "CL-Schema", // What is queried in the service
      "serviceEndpoint": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema" // the Resource that is returned (Schema Entity)
    }
  ]
}
```

**`SCHEMA` State format:**

- `"schema:<id>" -> {SchemaEntity, txHash, txTimestamp}`

`id` example: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue`

**Note**

This DID Document will be returned if the Schema DID Document URL, did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue, is **Resolved**

If the Schema's specific Entity, did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema, if attempted to be **Resolved**, the Resolver will **Dereference** the URL and will return the specific schema, found within the "service" section of the DID Document: did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema.

#### Schema Entity URL

If the Schema Entity URL is fetched through dereferencing, the following information is returned:

  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-Schema",
  "type": "CL-Schema",
  "controller": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu",  // Schema Issuer DID
  "version": "1.0",
  "name": "Degree",
  "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]

#### Updating a Schema

It is not possible to update an existing Schema using this architecture. This is because there are no verification methods specified for the DID Document. Therefore, it is a persistent Schema to ensure the original schema used to issue any credentials in the past are always available.

If a Schema evolves, a new Schema with a new version or name needs to be created.


#### Credential Definition

[TODO: explain that a Cred Def is simply an additional property inside of
the Issuer's DID Doc]

Adds a Credential Definition (in particular, public key), which is created by an
Issuer and published for a particular Credential Schema.

It is not possible to update Credential Definitions. If a Credential Definition
needs to be evolved (for example, a key needs to be rotated), then a new
Credential Definition needs to be created for a new Issuer DIDdoc.
Credential Definitions is added to the ledger in as verification method for
Issuer DIDDoc

- **`id`**: DID as base58-encoded string for 16 or 32 byte DID value with Cheqd
DID Method prefix `did:cheqd:<namespace>:` and a resource
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
    "controller": "did:cheqd:mainnet-1:123456789abcdefghi",
    "schemaId": "did:cheqd:mainnet-1:5ZTp9g4SP6t73rH2s8zgmtqdXyT?service=CL-Schema",
    "tag": "some_tag",
    "value": {
      "primary": "...",
      "revocation": "..."
    }
}
```

Don't store Schema DIDDoc in the State.

CredDef URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue`

CredDef Entity URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-CredDef`

`CRED_DEF` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
  "controller": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu", // CredDef Issuer DID
  "service":[
    {
      "id": "cheqd-cred-def",
      "type": "CL-CredDef",
      "serviceEndpoint": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?service=CL-CredDef"
    }
  ]
}
```

`CRED_DEF` state format:

`"credDef:<id>" -> {CredDefEntity, txHash, txTimestamp}`

**Note**: `CRED_DEF` **cannot** be updated.

### Rationale and Alternatives

#### Schema options not used

##### Option 2

Schema URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue#<schema_entity_id>`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
  "schema":[
    {
      "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue#schema1",
      "type": "CL-Schema",
      "controller": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
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

Schema URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
  "schema": {
              "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
              "type": "CL-Schema",
              "controller": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
              "value": {
                "version": "1.0",
                "name": "Degree",
                "attrNames": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
              },
            },
}
```

##### Option 4

Schema URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue#<schema_entity_id>`

`SCHEMA` DID Document transaction format:

```jsonc
{
  "id": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue",
  "schema":[
              {
                "id": "cheqd-schema",
                "type": "CL-Schema",
                "schemaRef": "did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue?resource=true"
              }
          ]
}
```

`SCHEMA` State format:

- `"schema:<id>" -> {SchemaEntity, txHash, txTimestamp}`

`id` example: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue`

#### Cred Def options not used

##### Option 2

Store inside Issuer DID Document

CredDef URL: `did:cheqd:mainnet-1:N22KY2Dyvmuu2PyyqSFKue#<cred_def_entity_id>`

`CRED_DEF` DID Document transaction format:

```jsonc
{
  "@context": [
    "https://www.w3.org/ns/did/v1",
    "https://w3id.org/security/suites/jws-2020/v1",
    "https://w3id.org/security/suites/ed25519-2020/v1"
  ],
  "verificationMethod": [
    {
      "id": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu#creddef-1", // Cred Def ID
      "type": "CamLysCredDefJwk2021", // TODO: define and register this key type
      "controller": "did:cheqd:mainnet-1:IK22KY2Dyvmuu2PyyqSFKu"
      "publicKeyJwk":
      {
        [TODO: Define structure for CL JWK]
      }
    }
  ],
  "assertionMethod": [
    "#creddef-1"
  ]

}
```

`CRED_DEF` state format:

###### Positive

- Credential Definition is a set of Issuer keys. So storing them in Issuer's DIDDoc reasonable.

###### Negative

- Credential Definition name means that it contains more than just a key and `value` field
  provides this flexibility.
- Adding all Cred Defs to Issuer's DIDDoc makes it too large. For every DIDDoc or Cred Def request
  a client will receive the whole list of Issuer's Cred Defs.
- Impossible to put a few controllers for Cred Def.
- In theory, we need to make Credential Definitions mutable.

## Consequences

### Backward Compatibility

### Positive

### Negative

### Neutral

## References

- [Hyperledger Indy](https://wiki.hyperledger.org/display/indy) official project background on Hyperledger Foundation wiki
  - [`indy-node`](https://github.com/hyperledger/indy-node) GitHub repository: Server-side blockchain node for Indy ([documentation](https://hyperledger-indy.readthedocs.io/projects/node/en/latest/index.html))
  - [`indy-plenum`](https://github.com/hyperledger/indy-plenum) GitHub repository: Plenum Byzantine Fault Tolerant consensus protocol; used by `indy-node` ([documentation](https://hyperledger-indy.readthedocs.io/projects/plenum/en/latest/index.html))
  - [Indy DID method](https://hyperledger.github.io/indy-did-method/) (`did:indy`)
  - [Indy identity-domain transactions](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md)
- [Hyperledger Aries](https://wiki.hyperledger.org/display/ARIES/Hyperledger+Aries) official project background on Hyperledger Foundation wiki
  - [`aries`](https://github.com/hyperledger/aries) GitHub repository: Provides links to implementations in various programming languages
  - [`aries-rfcs`](https://github.com/hyperledger/aries-rfcs) GitHub repository: Contains Requests for Comment (RFCs) that define the Aries protocol behaviour
- [W3C Decentralized Identifiers (DIDs)](https://www.w3.org/TR/did-core/) specification
  - [DID Core Specification Test Suite](https://w3c.github.io/did-test-suite/)
- [Cosmos blockchain framework](https://cosmos.network/) official project website
  - [`cosmos-sdk`](https://github.com/cosmos/cosmos-sdk) GitHub repository ([documentation](https://docs.cosmos.network/))
- [Sovrin Foundation](https://sovrin.org/)
  - [Sovrin Networks](https://sovrin.org/overview/)
  - [`libsovtoken`](https://github.com/sovrin-foundation/libsovtoken): Sovrin Network token library
  - [Sovrin Ledger token plugin](https://github.com/sovrin-foundation/token-plugin)
