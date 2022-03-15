# ADR 003: Command Line Interface (CLI) tools

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Alexandr Kolesov |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Implemented |
| **Start Date** | 2021-09-10 |

## Summary

Due to the nature of the cheqd project merging concepts from the [Cosmos blockchain framework](https://github.com/cosmos/cosmos-sdk) and self-sovereign identity (SSI), there are two potential options for creating Command Line Interface (CLI) tools for developers to use:

1. **Cosmos-based CLI:** Most likely route for Cosmos projects for their node application. Most existing Cosmos node validators will be familiar with this method of managing their node.
2. **VDR CLI**: Traditionally, a lot of SSI networks have used [Hyperledger Indy](https://github.com/hyperledger/indy-node) and therefore the Indy CLI tool for managing and interacting with the ledger. This has now been renamed to [Verifiable Data Registry (VDR) Tools CLI](https://gitlab.com/evernym/verity/vdr-tools) and is the tool that most existing SSI node operators ("stewards") would be familiar with.

Ideally, the `cheqd-node` project would provide a consistent set of CLI tools rather than two separate tools with varying feature sets between them.

This ADR will focus on the CLI tool architecture choice for `cheqd-node`.

## Context

### Assumptions / Considerations

#### Likelihood of introducing bugs or security vulnerabilities

1. Any CLI tool architecture chosen should not increase the likelihood of introducing bugs, security vulnerabilities, or design pattern deviations from upstream Cosmos SDK.
2. Actions that are carried out on ledger through a CLI tool in `cheqd-node` now include token functionality as well as identity functionality. E.g., if a DID gets compromised, there could be mechanisms to recover or signal that fact to issuers/verifiers. If tokens or the staking balance of node operators get compromised, this may potentially have more severe consequences for them.

#### User / Node Operator preferences

1. Would node operators want a single CLI to manage everything?
   1. This might be the case with node operators from an SSI / digital identity background, or node operators familiar with Hyperledger Indy CLI / VDR Tools CLI.
   2. A “single CLI” could be a single tool as far as the user sees, but actually consist of multiple modules beneath it in how it’s implemented.
2. Would node operators be okay with having two separate CLIs?
   1. One for Cosmos-ledger functions, and one for identity-specific functions.
   2. Unlike existing Hyperledger Indy networks, it is anticipated that some of the node operators on the cheqd network will have experience running Cosmos validator nodes. For this group, having to learn a new “single” CLI tool could cause a steeper learning curve and a worse user experience than what they have now.
   3. Node operators may want one/separate CLIs for security and operational reasons, i.e., for a separation of concerns in terms of functionality.

### Options considered

#### 1. Keep both Cosmos CLI and VDR Tools CLI, but use them for different purposes

**Pros:**

* Simple to do, no changes needed in code developed.
* Differences in functionality between the two CLIs can be explained in documentation.
* Node operators with good technical skills will understand the difference.
* Cosmos CLI design patterns would be consistent the wider Cosmos open source ecosystem.
* No steep learning curve for potential node operators who only want to run a node, without implementing SSI functionality in apps.

**Cons:**

* Key storage for Cosmos accounts may need to be done in two different keystores.
* Potentially confusing for node operators who use both CLIs to know which one to use for what purpose.
* Potentially a steeper learning curve for existing SSI node operators.

#### 2. Implement overlapping functionality in both CLI tools

**Pros:**

* Both Cosmos CLI and VDR Tools CLI would have native support for identity as well as token transactions.
* Node operators/developers could pick their preferred CLI tool.

**Cons:**

* Significant development effort required to implement very similar functionality two separate times, with little difference to the end user in actions that can be executed.
* VDR Tools CLI has DID / VC modules that would take significant effort to recreate in Cosmos CLI
* Cosmos CLI has token related functionality that would take significant development effort to replicate in VDR Tools CLI, and opens up the possibility that errors in implementation could introduce security vulnerabilities.

#### 3. Create aliases for commands in one of the CLI tools in the other CLI tool

_Commands in the Cosmos CLI could be made available as aliases in the VDR Tools CLI, or vice versa._

**Pros:**

* Single CLI tool to learn and understand for node operators.
* Development effort is simplified, as overlapping functionality is not implemented in two separate tools.

**Cons:**

* Less development effort required than Option 2, but greater than Option 1.
* Opens up the possibility that there's deviation in feature coverage between the two CLIs if aliases are not created to make 1:1 feature parity in both tools.

## Decision

Based on the options considerations above and an analysis of development required, the decision was taken to maintain two separate CLI tools:

1. **`cheqd-node` Cosmos CLI**: Any Cosmos-specific features, such as network & node management, token functionality required by node operators, etc.
2. **VDR Tools CLI**: Any identity-specific features required by issuers, verifiers, holders on SSI networks.

### CLI tools feature matrix

| Only available in Cosmos CLI | In both Cosmos CLI and VDR CLI | Only available in VDR CLI |
| :--- | :--- | :--- |
| Cosmos transactions + Queries | Signing service + Key storage | Identity transactions + Queries |
| MultiSig | (Transaction + Query) sending + Proof validation | DIDs + VCs (+ DID storage) |
| Network bootstrapping commands |  |  |

### CLI components overview


> [Editable versions of the diagrams](https://github.com/cheqd/cheqd-node/tree/e5f850355609f35a9a62c557ebf4adc73e766a44/architecture/adr-list/assets/adr003-cli-components-editable.excalidraw) (in Excalidraw format)

## Consequences

### Positive

* Faster time-to-market on the CLI tools, while freeing up time to build out user-facing functionality.

### Negative

* Cosmos account keys may need to be replicated in two separate key storages. A potential resolution for this in the future is to integrate the ability to use a single keyring for both CLI tools.

### Neutral

* Seek feedback from cheqd's open source community and node operators during testnet phase on whether the documentation and user experience is easy to understand and appropriate tools are available.

## References

