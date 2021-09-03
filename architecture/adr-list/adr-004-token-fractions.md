---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 004: Token fractions

## Status

ACCEPTED \| Not Implemented

## Summary

The aim of this ADR is to define the smallest token fraction on cheqd network.

## Context

Cosmos SDK doesn't provide native support for token fractions. The minimal amount you can operate in transactions is 1token. To address this issue networks assume that they use **N** digits after the decimal point and multiply all values by **10^\(-N\)** in UI.

### Examples

How many digits after the decimal point popular networks use:

* Cosmos - **6**
* IRIS - **6**
* Fetch.ai - **18**
* Binance - **8**

## Decision

It was decided to go with 10^-9 as the smallest fraction and call it **ncheq** \(nano cheq\).

## Consequences

### Backward Compatibility

There is no backward compatibility. To adjust the number of digits after the decimal point network should be restarted.

### Positive

* The value chosen is more precise than in the Cosmos chain so we have a reserve.

### Negative

* This decision is hard to change in the future.

### Neutral

## References

* [Cosmos ADR proposal to add coin metadata](https://docs.cosmos.network/master/architecture/adr-024-coin-metadata.html)

