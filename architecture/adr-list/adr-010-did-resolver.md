
# ADR 010: DID Resolver

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Renata Toktar|
| **ADR Stage** | DRAFT |
| **Implementation Status** | Not Implemented |
| **Start Date** | 22 February 2022 |

## Summary

This document defines the architecture of two DID resolvers: a DID resolver for the [cheqd DID method](adr-002-cheqd-did-method.md), and a DID driver for the [Universal Resolver](https://github.com/decentralized-identity/universal-resolver) project.

## Context

First, let's look at the discrepancies between the DID Doc model stored in the ledger and the format from the specification. Data is written to the ledger based on protobufs, which is due to the use of the Cosmos framework. However, according to the specification, DID Doc must be provided in JSON format, which provides more features that are not achievable in protobuffs.
Therefore, a new DID resolver is needed.

Inconsistencies between DIDDoc from the ledger and specification that should be corrected:

1. Rename "context" to "@context" - will be done on the ledger side
2. Change snake_case to camelCase for field names - will be done on Cheqd resolver side using the marshaller setting
3. Remove empty lists  - will be done on Cheqd resolver side using the marshaller setting
4. Convert a list of pairs to a map for jwk_pubkey and other cases
5. In metadata make `did` property to be of `DID` type instead of `Any`
6. DID URL Dereferencing: handle links to provide DID fragments and convert them to the desired format - Cheqd DID resolver functionality.

### Design principles

 helping to define its architecture in detail:

- Parallel executing of requests
- Synchronous replying for client requests (?)
- Marshal/unmarshal JSON - object - protobuff
- Programming language: Go (?)

## Architecture of DID Resolver(s)

### DID Resolution from the cheqd network ledger

Host the resolver separately from the ledger as an additional resolution service. Interaction with other applications and resolvers will implement the following schema:

![Cheqd did resolver](assets/adr-010-did-resolver/universal-resolver-sequence-diagram.png)
[Cheqd did sequence diagram. Schema 1.](assets/adr-010-did-resolver/universal-resolver-sequence-diagram.puml)

#### Pros

- Updating the resolver software does not need updating the application on the node side
- Separation of the system into microservices, moving away from a monolithic structure

#### Cons

For using the resolver separately from the ledger as an additional resolution service. 
All options for application interaction will be described in more detail below in [Possible flows for DID resolution](#possible-flows-for-did-resolution).

- Longer chain of trust. As a result, more resources required by the client to maintain the security of the system (`node + resolver` instead of `node`)

### Possible flows for DID resolution

To level out downsides of this approach a client can choose one of suitable flows.

#### 1.  


Cheqd DIDDoc resolving module at this stage will implement simple functionality that can be a lightweight architecture of threads without synchronization. Just several classes without the use of complex design patterns.

![cheqd did resolver class diagram](assets/adr-010-did-resolver/resolver-class-diagram.png)
[Cheqd did resolver class diagram](assets/adr-010-did-resolver/resolver-class-diagram.puml)

There are two ways to use Cheqd DIDDoc resolving module. As a library (go module) and as a standalone web service.

- In the first case, the module can be imported simply by adding the necessary import:

```golang
import (
     "github.com/cheqd/cheqd-did-resolver/src"
)
```

- In the second one, by launching the application with the command.

```bash
go run cheqd_resolver_web_service.go
```

### Universal resolver driver

A universal resolver driver is not required, since the resolver web application will be able to implement all the needs of Universal resolver. In case of unexpected issues with its integration, the cheqd-did-resolver module can always be used to import as a library.

Put the resolver inside Cheqd-node as a new module or as a new handler (keeper) inside the node application.

#### Pros

The presentation of the data takes place next to the base where the data is stored. This

- speeds up the process due to because of unnecessary data transferring between services
- does not allow compromising the resolver, only the entire node, which is a more difficult task
- fault tolerance and availability of the blockchain network is higher than a single web server

#### Cons

- Unable to update resolver without updating node. However, expanding the functionality without breaking changes is also possible with minor releases, which allows update the node without upgrade transaction.

## Decision

Add a new application/library for did-resolution

## Consequences

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## References

- [W3C Decentralized Identifiers (DIDs)](https://www.w3.org/TR/did-core/) specification
  - [DID Core Specification Test Suite](https://w3c.github.io/did-test-suite/)

## Unresolved questions

- Should the web service find another node for the request if it is not possible to connect to the node? So will the web service have a pool of trusted nodes for requesting?
- Synchronous replying for client requests?
