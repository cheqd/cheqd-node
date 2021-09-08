# ADR 007: Fees on cheqd network

## Status

| Category | Status |
| :--- | :--- |
| **ADR Stage** | PROPOSED |
| **Implementation Status** | Not Implemented |

## Summary

The goal of this ADR is to define how fees will be defined and processed on the cheqd network.

A linked ADR associated with this is [ADR 006: **Community tax**](adr-006-community-tax.md), ****which describes additional steps that take place based on the fees for a particular transaction on the network.

## Context

To do understand how fees would be implemented on the cheqd network, it is important to understand what capabilities [Tendermint](https://docs.tendermint.com/master/) and [Cosmos SDK](https://docs.cosmos.network/v0.43/basics/gas-fees.html) provide to set and implement transaction fees.

From [Cosmos SDK's **Introduction to `Gas` and `Fees`**](https://docs.cosmos.network/v0.43/basics/gas-fees.html#introduction-to-gas-and-fees):

> In the Cosmos SDK, **`gas` is a special unit that is used to track the consumption of resources during execution**. `gas` is typically consumed whenever read and writes are made to the store, but it can also be consumed if expensive computation needs to be done. It serves two main purposes:

> * **Make sure blocks are not consuming too many resources and will be finalized**. This is implemented by default in the SDK via the [block gas meter](https://docs.cosmos.network/v0.43/basics/gas-fees.html#block-gas-meter).
> * **Prevent spam and abuse from end-user**. To this end, `gas` consumed during [`message`](https://docs.cosmos.network/v0.43/building-modules/messages-and-queries.html#messages) execution is typically priced, resulting in a `fee` \(`fees = gas * gas-prices`\). `fees` generally have to be paid by the sender of the `message`. Note that the SDK does not enforce `gas` pricing by default, as there may be other ways to prevent spam \(e.g. bandwidth schemes\). Still, most applications will implement `fee` mechanisms to prevent spam. This is done via the [`AnteHandler`](https://docs.cosmos.network/v0.43/basics/gas-fees.html#antehandler).

### Investigations related to `gas` and `fees` parameters

#### `gas_wanted`

`gas_wanted` is the amount of gas that a transaction is allowed to use and is set by the transaction _sender_.

While there is an estimation functionality in Cosmos CLI to predict what the gas needed for a transaction should be, it doesn't take into account reads and writes to state. As a consequence, it can be inaccurate and result in a scenario where `gas_wanted` is underestimated.

The `gas_wanted` estimation functionality can also be 

* Is it accessible via \`RPC\` call?
  * Yes, we can query `app/simulate` endpoint

Gas used:

* What is it?
  * The amount of gas used by a transaction.
  * Always less than `gas_wanted`.
* What happens if a transaction requires more gas than specified in `gas_wanted`?
  * Execution interrupts
  * The fee isn't returned \(fully charged\)

Gas prices:

* What is transaction's `gas_prices`?
  * It's prises of gas suggested by transaction sender
  * Set by transaction sender
  * Can be any value
* What is validator's `min_gas_prices`?
  * It's a filter used to include transactions in the validator's mempool.
  * Transactions with less `gas_prices` are not considered as candidates for inclusion into the next blocks.
  * Specific for each validator.
  * Set by validator operator in `app.toml`.
* What happens if `gas_price` is less then `min_gas_prices` of all validators?
  * No validators include this transaction in mempool, it times out.
* What happens if `gas_price` is less then `min_gas_prices` of some validators?
  * The transaction has chance to be committed within timeout.
  * It depends on the sum of voting powers of validators that are ready to commit it.
* Can `gas_prices` be requested?
  * There is no way to request gas prices.

Fee:

* What is it?
  * The amount of tokens validator takes for a transaction processing.
* How is it calculated?
  * `fee` = `gas_wanted` \* `gas_price`
* Is the extra fee that is unused returned?
  * No, all fee suggested by user is charged.

## Decision

Proposals:

* Gas estimation:
  * Option 1: Implement gas estimation request in VDR tools and find out adjustment coefficient
  * Option 2: Estimate gas for most common transactions in advance
    * Can be used as a workaround while estimation request isn't implemented
* Fee prices estimations:
  * Option 1: Set recommended gas price for the network, embed it into the applications
  * Option 2: Dynamically determine the gas price based on recent transactions
    * Can be either implemented on client size or provided as a service
  * Then use exponential growth if the fee isn't enough

Transaction sending algorithm for VDR tools library consumer:

1. Build and sign a transaction
2. Send gas estimation request
3. Set initial gas price to a value:
   1. Either proposed by the community;
   2. Or retrieved from recent transactions, median for example:
      1. Can be retrieved on a client-side;
      2. Or received from public service.

* Try to send the transaction
  * Exponentially increase the gas limit in case of gas limit failure
  * Exponentially increase gas price in case of time out

## Consequences

* We use standard Cosmos mechanisms for fee estimation
* Client-side becomes pretty complex

### Backward Compatibility

This proposal is compatible with all recent versions of Cosmos

### Positive

### Negative

* Client-side logic complication

### Neutral

## References

* ADR 006: Community Tax
* [Tendermint Core documentation](https://docs.tendermint.com/master/)
* [Cosmos SDK documentation](https://docs.cosmos.network/)
* Github issues with suggestions to improve fee mechanism:
  * [https://github.com/cosmos/cosmos-sdk/issues/6555](https://github.com/cosmos/cosmos-sdk/issues/6555)
  * [https://github.com/cosmos/cosmos-sdk/issues/2150](https://github.com/cosmos/cosmos-sdk/issues/2150)
  * [https://github.com/cosmos/cosmos-sdk/issues/4938](https://github.com/cosmos/cosmos-sdk/issues/4938)

