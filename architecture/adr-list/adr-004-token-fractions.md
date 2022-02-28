---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 004: Token fractions

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Alexandr Kolesov |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Not Implemented |
| **Start Date** | 2021-09-08 |

## Summary

The aim of this ADR is to define the smallest fraction for CHEQ tokens.

## Context

Cosmos SDK doesn't provide native support for token fractions. The lowest denomination out-of-the-box that can be used in transactions is `1token`.

To address this issue, similar Cosmos networks assume that they use **N** digits after the decimal point and multiply all values by **10^(-N)** in UI.

### Examples of lowest token denominations in Cosmos

Popular Cosmos networks were compared to check how many digits after the decimal point are used by them:

* Cosmos: **6**
* IRIS: **6**
* Fetch.ai: **18**
* Binance: **8**

## Decision

Fractions of CHEQ tokens will be referred by their [SI/metric prefix](https://en.wikipedia.org/wiki/Metric_prefix#List_of_SI_prefixes), based on the power of 10 of CHEQ tokens being referred to in context. This notation system is common across other Cosmos networks as well.

It was decided to go with **10^-9** as the smallest fraction, with the whole number token being 1 CHEQ. Based on the SI prefix system, the lowest denomination would therefore be called "**nanocheq**".

## Consequences

### Backward Compatibility

* There is no backward compatibility. To adjust the number of digits after the decimal point (lowest token denomination), the network should be restarted.

### Positive

* The power of 10 chosen for the lowest denomination of CHEQ tokens is more precise than for Cosmos ATOMs, which allows transactions to be defined in smaller units.

### Negative

* This decision is hard to change in the future, as changes to denominations require significant disruption when a network is already up and running.

### Neutral

* N/A

## References

* [Cosmos ADR proposal to add coin metadata](https://docs.cosmos.network/master/architecture/adr-024-coin-metadata.html)

