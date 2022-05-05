# ADR 008: cheqd DIDDoc resources: Schemas and Credential Definitions

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar, Brent Zundel, Ankur Banerjee |
| **ADR Stage** | DRAFT |
| **Implementation Status** | Draft |
| **Start Date** | 2021-09-23 |

## Summary

This Architecture Decision Record (ADR) will define how Verifiable Credential schemas can be represented through the use of a DID URL, which when dereferenced, fetches the credential schemas a resource. 

In this ADR we look only at CL (Camenisch-Lysyanskaya) schemas and credential definition transactions that are needed for verifying issued Verifiable Credentials.

The identity entities and transactions for the cheqd network are designed to extend usage scenarios and functionality currently supported by [Hyperledger Indy](https://github.com/hyperledger/indy-node) into a model which is [W3C DID Core compliant](https://www.w3.org/TR/did-core/). 


## Context

Hyperledger Indy is a verifiable data registry (VDR) built for DIDs with a strong focus on privacy-preserving techniques. It is one of the most widely-adopted SSI blockchain ledgers. Most notably, Indy is used by the [Sovrin Network](https://sovrin.org/overview/).

Our aim is to support the functionality enabled by [identity-domain transactions in by Hyperledger Indy](https://github.com/hyperledger/indy-node/blob/master/docs/source/transactions.md) into `cheqd-node`. This will partly enable the goal of allowing use cases of existing SSI networks on Hyperledger Indy to be supported by the cheqd network.

### Identity-domain transaction types in Hyperledger Indy

The following identity-domain transactions from Indy were considered:

1. `NYM`: Equivalent to "DIDs" on other networks
2. `ATTRIB`: Payload for DID Document generation
3. `SCHEMA`: Schema used by a credential
4. `CRED_DEF`: Credential definition by an issuer for a particular schema
5. `REVOC_REG_DEF`: Credential revocation registry definition
6. `REVOC_REG_ENTRY`: Credential revocation registry entry

Revocation registries for credentials are not covered under the scope of this ADR. This topic is discussed separately in [ADR 007: **Revocation registry**](adr-007-revocation-registry.md) as there is ongoing research by the cheqd project on how to improve the privacy and scalability of credential revocations.

### Resolving DID vs Dereferencing DID

Before diving into the specific architecture of cheqd's schemas, it is imporntant to explain the difference between resolving a DID and dereferencing a DID URL. 

When you resolve a DID, a DID Document is returned. For example, resolving: "did:cheqd:example1234" would return the full DID Document associated with the specific DID "did:cheqd:example1234". 

Dereferecing a DID URL is slightly different. When you dereference a DID URL, you are parsing the URL for specific actions, such as to take a certain path, point to a specific fragment, or query a specific resource. 

For example, "did:cheqd:example1234?service=CLSchema" can be dereferenced. In this case it will query the service of type "CLSchema" within the DID Document and will return the resource specified at the Service Endpoint within the DID Document. It is this type of architecture that cheqd uses in this ADR to fetch schemas. 

## Decision

### Schema

The most important concept to understand for this architecture is that each Schema will be written to the cheqd network as its own DID, with its own DID Document, which is able to be resolved or dereferenced.  

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
This creates an on-ledger artifact, that when dereferenced will return the appropriate fields of the schema. 

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


## Credential Definitions

A Credential Definition is an on-ledger artifact used in Hyperledger Indy. Its purpose is:
- To references the DID of the issuer and the public keys of the issuer
- To references a Verifiable Credential schema and for each attribute in the schema, list a public key so that verifier can verify that each claim in the proof has not been tampered with.
- To references the Revocation Registry to evaluate the holders "proofs of non-revocation"
- To list the specific attributes which may be included within a presented claim by the holder

A single credential definition can incorporate attribute data types from multiple schemas. This encourages the reuse of the same data types and schema definitions across multiple credentials, encouraging interoperability.

The benefits of Credential Defintiions are that they add an extra layer of compostibility to Verifiable Credential claims and presentations, enabling a holder to selectively disclose specific attributes from multiple different Credentials in one claim, and prove that they have not been tamperered with. This is possible because each attribute is linked atomically to a public key. This means that multiple attributes, from multiple Credentials, can be bundled into a claim that is still verifiable and trustworthy. 

However, with advancements in the W3C DID Core specification, such as the ability to have Verification Keys usable **only** for specific purposes, the compostability of Credential Definitions has become much less important. 

For example, A DID Document could specify that did:cheqd:example:1234#key1 can **only** be used as an assertion method for signing claims in Verifiable Credentials. Next, when a holder presents such claims to a verifier, the verifier will be able to resolve the issuer's DID included in the Presentation and dereference the signing key used as a Primary Resource. This will prove to the verifier that this claim is definitively from the issuer. In order to check that the claims have not been tampered with, the verifier is able to also verify the schema and syntax of the claims presented, since the schema will also be referenced within the Verifiable Credential. Here, the schema in combination with Verification Method Relationships can augment large amounts of the purpose of the Credential Definition.

#### Consequences

##### Backward Compatibility

- Through augmenting the necessity for Credential Definitions using schemas and Verification Method Relationships, Hyperledger Indy Credentials can be closely replicated within a format which is W3C DID Core compliant. This will enable companies who currently use Indy Credentials on Indy to switch to using Indy Credentials on cheqd without needing to rip and replace the existing processes. 

##### Positive

- Credential Definition is a set of Issuer keys. So storing them in Issuer's DIDDoc is reasonable.
- Blending Indy Credentials into a more standardised format will increase semantic interoperability.
- cheqd will be able to support JSON/JSON-LD and a closely-tied and compatible format of AnonCreds.
- Schemas on ledger is far more secure than using a web-based service like [schema.org](https://schema.org/). 

##### Negative

- Since schemas do not have a public key per attribute, but reference a set of attributes, cheqd-based Credentials with CL-schemas on ledger will not enable the same compostability in terms of predicate proofs as Indy-based Credentials. 
- Indy-based revocation is currently operative and functional, whilst cheqd's privacy-preserving revocation is still a work in progress, and in the meantime, companies may need to rely on RevocationList2020 or StatusList 2021.

##### Neutral

### Rationale and Alternatives

#### Schema options not used

##### Option 1

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

##### Option 2

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

##### Option 3

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
