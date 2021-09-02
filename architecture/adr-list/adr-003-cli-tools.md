---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 003: CLI tools

## Status

ACCEPTED \| Implemented

## Summary

Due to the nature of the verim product, there are two CLI tools for operating the network:

* Cosmos CLI \(which comes with Cosmos, but doesn’t yet have any identity related modules\)
* VDR CLI \(previously known as Indy CLI\): traditionally, this is the tool that node operators / stewards would have used 

Now there are two different CLI tools with different feature sets available to operate with verim networks. Need to decide how CLI tools should look like for the end user.

## Context

### Assumptions / Considerations

1. Anything that is done by verim or Evernym should NOT increase the likelihood of introducing potential bugs or security vulnerabilities in Cosmos CLI. This is a must-have requirement.
   1. Cosmos CLI is battle-tested and reviewed by the wider Cosmos open source community.
   2. Cosmos features/functions are also what directly relate to the most financially-sensitive part of the verim product stack. E.g., losing control of a DID is bad, but recoverable. Losing control of tokens or stake can have massive financial impact or loss.
2. Would node operators want a single CLI to manage everything?
   1. This might be the case with node operators from an SSI / digital identity background, or node operators coming from Hyperledger Indy
   2. A “single CLI” could be a single tool as far as the user sees, but actually consist of multiple modules beneath it in how it’s implemented
3. Or, would node operators be okay with having two separate CLIs: one for Cosmos-related functions, and one for verim/identity specific functions.
   1. Unlike Sovrin/Hyperledger Indy, many of the node operators on the verim network will be existing native Cosmos validators. For this group, having to learn a new “single” CLI tool is a worse user experience than what they have now.
   2. It might be acceptable for node operators to have two separate CLIs for security reasons too, i.e., for a separation of concerns in terms of functionality.

### Current architecture overview

![Current Cosmos CLI and VDR tools CLI architecture](https://lh3.googleusercontent.com/cMdfEe19vqDVaRJ0kP97KGCUHauEpnh2TV1OhmvGqOFqqIkhWXGkdKxONDLjW2rnU83k9yelFWK_jhsqQoF57tNf8ChrPeIZsiLys3LKVT_QKG9Gk7Mir4ChbCeiUKs2V7l7jE8d=s0)

### Features overview

| Cosmos only | Need to decide whether on Cosmos CLI or VDR CLI | VDR CLI only |
| :--- | :--- | :--- |
| Cosmos Txs + Queries | Signing service + Key storage | verim Txs + Queries |
| Multisig | \(Tx + Queries\) sending + Proof validation | DIDs + VCs \(+ DID storage\) |
| Network bootstrapping commangs |  |  |

### Options

#### 1. Keep two separate CLIs for different purposes

Pros:

* Simple to do, no changes needed - explain the difference between them in documentation
* Might be acceptable enough to dev audience

Cons:

* Two key storages for cosmos accounts

#### 3.2. Just create aliases in one direction

Pros:

* Single tool

Options:

* Option a: map Cosmos CLI commands as aliases to VDR CLI
  * Better ux for those who is familiar with Indy
* Option b: map VDR CLI commands as aliases to Cosmos CLI
  * Better ux for those who is familiar with Cosmos

#### 3.3. Actually move modules around between the two libraries

Cons:

* Potentially big effort

Considerations:

* VDR has DID / VC modules that are highly unlikely to be moved to Cosmos CLI
* If we try to move verim/VDR modules to Cosmos CLI, it is very likely that this will be our own standalone fork and won’t actually be integrated into the main Cosmos CLI release as this is functionality/project specific code. This effectively means we end up with a fork of Cosmos CLI, deviating from mainline branch, that we have to patch and maintain.
* If we try to make Cosmos functions easier to use in VDR, we have the possibility that a bug or vulnerability in VDR CLI somehow results in an action that makes tokens/stake vulnerable.
* Similarly, if we try to put in VDR functions in Cosmos CLI, we’re introducing an element of deviation from the mainline Cosmos CLI release. We may therefore accidentally end up introducing a security vulnerability in our modified version of Cosmos CLI \(with verim/VDR functionality\) that results in tokens/stake being lost.

## Decision

After the discussion on Jul 1 we agreed on:

* Keep 2 separate CLI tools for different purposes:
  * VC/DID stuff management
  * Cosmos network management
* Get feedback after net launch and decide how to improve UX
  * Possible direction - use single keyring back-end for CLIs

## Consequences

### Backwards Compatibility

* None

### Positive

* Simple to do, no changes needed - explain the difference between them in documentation
* Might be acceptable enough to dev audience

### Negative

* Two key storages for cosmos accounts

### Neutral

* None

## References

* {reference link}

