# ADR 002: cheqd DID method, identity entities, and transactions

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar, Brent Zundel, Ankur Banerjee |
| **ADR Stage** | PROPOSED |
| **Implementation Status** | Implementation in progress |
| **Start Date** | 2021-09-23 |

## TODO

- Describe the resolution process for DID
- Describe the deactivation process
- Add security considerations
- Add privacy considerations

## Summary

This ADR defines the cheqd DID method and describes the identity entities, queries, and transaction types for the cheqd network: a purpose-built self-sovereign identity (SSI) network based on the [Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk).

[Decentralized identifiers](https://www.w3.org/TR/did-core) (DIDs) are a type of identifier that enables verifiable, decentralized digital identity. A DID refers to any subject (for example, a person, organization, thing, data model, abstract entity, and so on) as determined by the controller of the DID.

## Context

### Rationale for baselining against Hyperledger Indy

Hyperledger Indy is a verifiable data registry (VDR) built for DIDs with a strong focus on privacy-preserving techniques. It is one of the most widely-adopted SSI blockchain ledgers. Most notably, Indy is used by the [Sovrin Network](https://sovrin.org/overview/).

The Sovrin Foundation initiated a project called [`libsovtoken`](https://github.com/sovrin-foundation/libsovtoken) in 2018 to create a native token for Hyperledger Indy. `libsovtoken` was intended to be a payment handler library that could work with `libindy` and be merged upstream. This native token would allow interactions on Hyperledger Indy networks (such as Sovrin) to be paid for using tokens.

Due to challenges the project ran into, the `libsovtoken` codebase saw its [last official release in August 2019](https://github.com/sovrin-foundation/libsovtoken/releases/tag/v1.0.1).

### Rationale for using the Cosmos blockchain framework for cheqd

The cheqd network aims to support similar use cases for SSI as seen on Hyperledger Indy networks, with a similar focus on privacy-resspecting techniques.

Since the core of Hyperledger Indy's architecture was designed before the [W3C DID specification](https://www.w3.org/TR/did-core/) started to be defined, the [Indy DID Method](https://hyperledger.github.io/indy-did-method/) (`did:indy`) has aspects that are not fully-compliant with latest specifications.

However, the [rationale for why the cheqd team chose the Cosmos blockchain framework instead of Hyperledger Indy](https://blog.cheqd.io/why-cheqd-has-joined-the-cosmos-4db8845722c5) were primarily down to the following reasons:

1. **Hyperledger Indy is a permissioned ledger**: Indy networks are permissioned networks where the ability to have write capability is restricted to a limited number of nodes. Governance of such a permissioned network is therefore also not decentralised.
2. **Limitations of Hyperledger Indy's consensus mechanism**: Linked to the permissioned nature of Indy are the drawbacks of its [Plenum Byzantine Fault Tolerant (BFT) consensus](https://github.com/hyperledger/indy-plenum) mechanism, which effectively limits the number of nodes with write capability to approximately 25 nodes. This limit is due to limited transactions per second (TPS) for an Indy network with a large number of nodes, rather than a hard cap implemented in the consensus protocol.
3. **Wider ecosystem for token functionality outside of Hyperledger Indy**: Due to its origins as an identity-specific ledger, Indy does not have a fully-featured token implementation with sophisticated capabilities. Moreover, this also impacts end-user options for ecosystem services such as token wallets, cryptocurrency exchanges, custodianship services etc that would be necessary to make a viable, enterprise-ready SSI ledger with token functionality.

By selecting the Cosmos blockchain framework, the maintainers of the cheqd project aim to address the limitations of Hyperledger Indy outlined above. However, with an eye towards interoperability, the cheqd project aims to use [Hyperledger Aries](https://wiki.hyperledger.org/display/ARIES/Hyperledger+Aries) for ledger-related peer-to-peer interactions.

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

Schema and Credential definition have been moved to another ADR.

## Decision

Identity entities and transactions for the cheqd network may differ in name from those in Hyperledger Indy, but aim enable equivalent support for privacy-respecting SSI use cases.

The differences stem primarily from aiming to achieve better compliance with the W3C [DID Core](https://www.w3.org/TR/did-core/) specification and architectural differences between Hyperledger Indy and [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) (used to build `cheqd-node`).

With better compliance against the DID Core specification, the goal of the **cheqd DID method** is to maximise interoperability with compatible third-party software librarires, tools and projects in the SSI ecosystem.

### DID Method Name

The `method-name` for the **cheqd DID Method** will be identified by the string `cheqd`.

A DID that uses the cheqd DID method MUST begin with the prefix `did:cheqd`. This prefix string MUST be in lowercase. The remainder of the DID, after the prefix, is as follows:

#### Method Specific Identifier

The cheqd DID `method-specific-id` is made up of a component:

**`namespace`**: A string that identifies the cheqd network `chain_id` (e.g., "mainnet", "testnet") where the DID reference is stored. Different cheqd networks may be differentiated based on whether they are production vs non-production, governance frameworks in use, participants involved in running nodes, etc.

A `did:cheqd` DID MUST be self-certifying by having the namespace component of the DID (last element)
derived from the initial public key of the DID, as follows:
For an Ed25519 key: Convert into Base58char the first 16 bytes of the 256 bit public key.

#### cheqd DID method syntax

The cheqd DID method ABNF to conform with [DID syntax guidelines](https://www.w3.org/TR/did-core/#did-syntax) is as follows:

```abnf
cheqd-did         = "did:cheqd:" namespace
namespace         = 1*namespace-char "-" 1*namespace-ver 
namespace-char    = ALPHA / DIGIT
namespace-ver     = ALPHA / DIGIT
```

#### Examples of `did:cheqd` identifiers

A DID written to the cheqd "mainnet" ledger `namespace`:

```abnf
did:cheqd:mainnet:7Tqg6BwSSWapxgUDm9KKgg
```

A DID written to the cheqd "testnet" ledger `namespace`:

```abnf
did:cheqd:testnet:6cgbu8ZPoWTnR5Rv5JcSMB
```

### DID Documents (DIDDocs)

A DID Document ("DIDDoc") associated with a cheqd DID is a set of data describing a DID subject. The [representation of a DIDDoc when requested for production](https://www.w3.org/TR/did-core/#representations) from a DID on cheqd networks MUST meet the DID Core specifications.

#### Elements needed for DIDDoc representation

1. **`id`**: Target DID as base58-encoded string for 16 or 32 byte DID value with cheqd DID Method prefix `did:cheqd:<namespace>:`.
2. **`controller`** (optional): A list of fully qualified DID strings or one string. Contains one or more DIDs who can update this DIDdoc. All DIDs must exist.
3. **`verificationMethod`** (optional): A list of Verification Methods
4. **`authentication`** (optional): A list of strings with key aliases or IDs
5. **`assertionMethod`** (optional): A list of strings with key aliases or IDs
6. **`capabilityInvocation`** (optional): A list of strings with key aliases or IDs
7. **`capabilityDelegation`** (optional): A list of strings with key aliases or IDs
8. **`keyAgreement`** (optional): A list of strings with key aliases or IDs
9. **`service`** (optional): A set of Service Endpoint maps
10. **`alsoKnownAs`** (optional): A list of strings. A DID subject can have multiple identifiers for different purposes, or at different times. The assertion that two or more DIDs refer to the same DID subject can be made using the `alsoKnownAs` property.
11. **`@context`** (optional): A list of strings with links or JSONs for
describing specifications that this DID Document is following to.

##### Example of DIDDoc representation

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

#### State format for DIDDocs on ledger

`"diddoc:<id>" -> {DIDDoc, DidDocumentMetadata, txHash, txTimestamp }`

`didDocumentMetadata` is created by the node after transaction ordering and before adding it to a State.

#### DIDDoc metadata

1. **`created`** (string): Formatted as an XML Datetime normalized to UTC
00:00:00 and without sub-second decimal precision. For example:
2020-12-20T19:17:47Z.
2. **`updated`** (string): The value of the property MUST follow the same
formatting rules as the created property. The `updated` field is null if an Update operation has never been performed on the DID document. If an updated property exists, it can be the same value as the created property when the difference between the two timestamps is less than one second.
3. **`deactivated`** (strings): If DID has been deactivated, DID document
metadata MUST include this property with the boolean value true. By default this is set to `false`.
4. **`versionId`** (strings): Contains transaction hash of the current DIDDoc version.

##### Example of DIDDoc metadata

```jsonc
{
  "created": "2020-12-20T19:17:47Z",
  "updated": "2020-12-20T19:19:47Z",
  "deactivated": false,
  "versionId": "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E",
}
```

#### Verification method

Verification methods are used to define how to authenticate / authorise interactions with a DID subject or delegates. Verification method is an OPTIONAL property.

1. **`id`** (string): A string with format `<DIDDoc-id>#<key-alias>`
2. **`controller`**: A string with fully qualified DID. DID must exist.
3. **`type`** (string)
4. **`publicKeyJwk`** (`map[string,string]`, optional): A map representing a JSON Web Key that conforms to [RFC7517](https://tools.ietf.org/html/rfc7517). See definition of `publicKeyJwk` for additional constraints.
5. **`publicKeyMultibase`** (optional): A base58-encoded string that conforms to a [MULTIBASE](https://datatracker.ietf.org/doc/html/draft-multiformats-multibase-03)
encoded public key.

**Note**: Verification method cannot contain both `publicKeyJwk` and
`publicKeyMultibase` but must contain at least one of them.

##### Example of Verification method in a DIDDoc

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

#### Service

Services can be defined in a DIDDoc to express means of communicating with the DID subject or associated entities.

1. **`id`** (string): The value of the `id` property for a Service MUST be a URI conforming to [RFC3986](https://www.rfc-editor.org/rfc/rfc3986). A conforming producer MUST NOT produce multiple service entries with the same ID. A conforming consumer MUST produce an error if it detects multiple service entries with the same ID. It has a follow formats: `<DIDDoc-id>#<service-alias>` or `#<service-alias>`.
2. **`type`** (string): The service type and its associated properties SHOULD be registered in the [DID Specification Registries](https://www.w3.org/TR/did-spec-registries/)
3. **`serviceEndpoint`** (strings): A string that conforms to the rules of [RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs, a map, or a set composed of a one or more strings that conform to the rules of
[RFC3986](https://www.rfc-editor.org/rfc/rfc3986) for URIs and/or maps.

##### Example of Service in a DIDDoc

```jsonc
{
  "id":"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#linked-domain",
  "type": "LinkedDomains",
  "serviceEndpoint": "https://bar.example.com"
}
```

### `DID` transactions

#### Create `DID`

If there is no DID entry on the ledger with the specified DID (`DID.id`), it is considered as a creation request for a new DID.

**Note**: The field `signatures`(from `WriteRequest`) must contain signatures from all new controllers.

**Client request**: 
`WriteRequest(CreateDidRequest(id, controller, verificationMethod, authentication, assertionMethod, capabilityInvocation, capabilityDelegation, keyAgreement, service, alsoKnownAs, context), signatures)`

- **`signatures`**: `CreateDidRequest` should be signed by all `controller` private keys. This field contains a dict there key's URI from
  `DIDDoc.authentication`, and the signature is a value. The `signatures` must contains signatures from all controllers.
- **`id`**: fully qualified did
- **`controller, verificationMethod, authentication, assertionMethod, capabilityInvocation, capabilityDelegation, keyAgreement, service, alsoKnownAs, context`**: optional parameters in accordance with DIDDoc properties

Example of a client request:

```jsonc
WriteRequest{
              "data": 
                      CreateDidRequest{   
                                        "context": [
                                            "https://www.w3.org/ns/did/v1",
                                            "https://w3id.org/security/suites/ed25519-2020/v1"
                                        ],
                                        "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
                                        "controller": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue"],
                                        "verificationMethod": [
                                          {
                                            "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#capabilityInvocationKey",
                                            "type": "Ed25519VerificationKey2020", // external (property value)
                                            "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
                                            "publicKeyMultibase": "z4BWwfeqdp1obQptLLMvPNgBw48p7og1ie6Hf9p5nTpNN"
                                          }
                                        ],
                                        "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#authKey1"],
              },
              "signatures": {
                  "verification method URI": "<signature>"
                  // Multiple verification methods and corresponding signatures can be added here
              },
}
```

#### Update `DID`

If there is a DID entry on the ledger with the specified DID (`DID.id`), then this considered a request for updating an existing DID.

For updating `versionId` from `UpdateDIDRequest` should be filled by a
transaction hash of the previous DIDDoc version. Can be received by GetDid query. It is needed for a replay protection.

**Note**: The field `signatures`(from `WriteRequest`) must contain signatures from all old controllers and all new controllers.

**Client request**:
`WriteRequest(UpdateDidRequest(id, controller, verificationMethod, authentication, assertionMethod, capabilityInvocation, capabilityDelegation, keyAgreement, service, alsoKnownAs, context, versionId), signatures)`

- **`signatures`**: `CreateDidRequest` should be signed by all `controller` private keys. This field contains a dict there key's URI from
  `DIDDoc.authentication`, and the signature is a value. The `signatures` must contains signatures from all controllers.
- **`id`**: fully qualified did
- **`versionId`**: transaction hash of the previous DIDDoc version.
- **`controller, verificationMethod, authentication, assertionMethod, capabilityInvocation, capabilityDelegation, keyAgreement, service, alsoKnownAs, context`**: optional parameters in accordance with DIDDoc properties

Example of a client request:

```jsonc
WriteRequest{
              "data": 
                     UpdateDidRequest{   
                                      "context": [
                                          "https://www.w3.org/ns/did/v1",
                                          "https://w3id.org/security/suites/ed25519-2020/v1"
                                      ],
                                      "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
                                      "controller": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue"],
                                      "verificationMethod": [
                                        {
                                          "id": "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#capabilityInvocationKey",
                                          "type": "Ed25519VerificationKey2020", // external (property value)
                                          "controller": "did:cheqd:mainnet:N22N22KY2Dyvmuu2PyyqSFKue",
                                          "publicKeyMultibase": "z4BWwfeqdp1obQptLLMvPNgBw48p7og1ie6Hf9p5nTpNN"
                                        }
                                      ],
                                      "authentication": ["did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#authKey1"],
                                      "versionId": "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E"
              },
              "signatures": {
                  "verification method URI": "<signature>"
                  // Multiple verification methods and corresponding signatures can be added here
              },
}
```

#### Resolve DID

Using GetDid query DIDDoc can be received from the ledger.
`QueryGetDidResponse(id)`

Response format from Tendermint RPC interface is `QueryGetDidResponse` protobuf:
```cjson


```

### Security Considerations

For adding a new DIDDoc or update an old one should be signed by controller signatures. 

Changing of a DIDDoc fragment without a controller (any field except VerificationMethods) MUST be signed by DIDDoc's controller(s).

Changing of a DIDDoc fragment with a controller MUST be signed by DIDDoc's controller(s) **and** fragment's controller.

Changing the controller requires a list of signatures as before for changing any field.

### Privacy Considerations

TODO: add privacy considerations

### Changes from Indy entities and transactions

#### Rename `NYM` transactions to `DID` transactions

[**NYM** is the term used by Hyperledger Indy](https://hyperledger-indy.readthedocs.io/projects/node/en/latest/transactions.html#nym) for DIDs. 

cheqd uses the term `DID` instead of `NYM` in transactions, which should
make it easier to understand the context of a transaction easier by bringing it closer to W3C DID terminology used by the rest of the SSI ecosystem.

#### Remove `role` field from `DID` transactions

Hyperledger Indy is a public-permissioned distributed ledger and therefore use the `role` field to distinguish transactions from different types of nodes. As cheqd networks are public-permissionless, the `role` scope has been removed.

#### Dropping `ATTRIB` transactions

`ATTRIB` was originally used in Hyperledger Indy to add document content similar to DID Documents (DIDDocs). The cheqd DID method replaces this by implementing DIDDocs for most transaction types.

## Consequences

### Backward Compatibility

- `cheqd-node` [release v0.1.17](https://github.com/cheqd/cheqd-node/releases/tag/v0.1.17) and 
  earlier had a transaction type called `NYM` which would allow writing/reading a unique identifier 
  on ledger. However, this `NYM` state was not fully defined as a DID method and did not contain DID Documents 
  that resolved when the DID identifier was read. This `NYM` transaction type is deprecated and the data written 
  to cheqd testnet with legacy states will not be retained.

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
