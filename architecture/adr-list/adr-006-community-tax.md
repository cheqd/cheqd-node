---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 006: Community tax

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Alexandr Kolesov |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Not Implemented |
| **Start Date** | 2021-09-08 |

## Summary

The aim of this ADR is to define how ["community tax" as described in the Cosmos blockchain framework](https://docs.cosmos.network/v0.44/modules/distribution/07_params.html#parameters) will work on cheqd network.

## Context

### What is "community tax"?

`communityTax` is a value set in genesis for each Cosmos network and defined as a percentage that is applied to the fees collected in each block.

Tokens collected through this process accumulate in the **community pool**. The percentage charged as `communityTax` can be changed by [making proposals on the network and voting for acceptance](https://docs.cosmos.network/v0.44/modules/gov/) by the network.

### Community tax collection

From [Cosmos SDK documentation, `distribution` module](https://docs.cosmos.network/master/modules/distribution/03_begin_block.html#reward-to-the-community-pool):

> The community pool gets `community_tax * fees`, plus any remaining dust after validators get their rewards that are always rounded down to the nearest integer value.

```text
communityFunding = feesCollectedDec * communityTax
feePool.CommunityFund += communityFunding
```

### Community tax distribution

To spend tokens from the **community pool**:

1. `community-pool-spend` proposal can be submitted on the network.
   1. Recipient address and amount of tokens should be specified.
   2. The purpose for which the requested community pools tokens will be spent should be described.
2. If proposal is approved using the voting process, the recipient address specified will receive the requested tokens.
3. The expectation on the recipient is that they spend the tokens for the purpose specified in their proposal.

More information about fee distribution is available in the [**End Block** section of Cosmos's `distribution` module](https://docs.cosmos.network/master/modules/distribution/03_begin_block.html) documentation.

## Decision

* cheqd's network will keep the `communityTax` parameter enabled, i.e., non-zero.
* The value of `communityTax`, based on a review of similar Cosmos networks will be set to `2%`.

## Consequences

### Backward Compatibility

* The behavior of `communityTax` is the across Cosmos SDK **v0.42** and **v0.43**.

### Positive

* The cheqd network will have a pool of tokens that can be used to spend on initiatives valued by the community.

### Negative

* N/A

### Neutral

* cheqd's Governance Framework should provide guidance on how to submit proposals and recommended areas of investment in community efforts.

## References

* [Cosmos SDK `distribution` module parameters](https://docs.cosmos.network/v0.44/modules/distribution/07_params.html#parameters)
* [Cosmos SDK `governance` module](https://docs.cosmos.network/v0.44/modules/gov/)

