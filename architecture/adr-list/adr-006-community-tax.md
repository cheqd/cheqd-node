---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 006: Community tax

## Status

PROPOSAL \| Not Implemented

## Summary

The aim of this ADR is to define how community tax will work on cheqd network.

## Context

### Community tax collection

In Cosmos **v0.42**:

* `communityTax` tax value is set in genesis and can be changed via proposals + voting
* The tax is applied to fees collected in each block
* Tokens charged accumulate in community pool

```text
communityFunding = feesCollectedDec * communityTax
feePool.CommunityFund += communityFunding
```

Here is more information about fees distribution [\[ref\]](https://docs.cosmos.network/v0.42/modules/distribution/03_end_block.html).

### Community tax distribution

To spend tokens from community pool:

* `community-pool-spend` proposal should be submitted
  * Recipient address and amount of tokens should be specified
* If proposal is approved recipient will recieve tokens and will be able to spend them

## Decision

Use the same community tax as in Cosmos mainnet: `2%`.

## Consequences

### Backward Compatibility

### Positive

* Community will receive tokens that can be spent on network support.

### Negative

### Neutral

## References

