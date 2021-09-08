# ADR 002: Fee processing

## Status

PROPOSED \| Not Implemented

## Summary

Goals of this ARD is to define how fees will be processed on cheqd network. To do it we first need to understand what capabilities Tendermint and Cosmos SDK provide.

## Context

Gas and fees in Cosmos SDK and Tendermint:

* Gas:
  * What is it?
    * It's a unit used to track the consumption of resources during transaction execution.
  * What is it used for?
    * Block computational complexity limitation;
    * Spam prevention.
* Gas wanted:
  * What is it?
    * The amount of gas that transaction is allowed to use.
    * Set by transaction sender
  * How can it be predicted?
    * There is an estimation functionality in Cosmos CLI.
    * It doesn't take into account reads and writes to state so it's inaccurate and can underestimate `gas_wanted`.
  * Is it accessible via \`RPC\` call?
    * Yes, we can query `app/simulate` endpoint
* Gas used:
  * What is it?
    * The amount of gas used by a transaction.
    * Always less than `gas_wanted`.
  * What happens if a transaction requires more gas than specified in `gas_wanted`?
    * Execution interrupts
    * The fee isn't returned \(fully charged\)
* Gas prices:
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
* Fee:
  * What is it?
    * The amount of tokens validator takes for a transaction processing.
  * How is it calculated?
    * `fee` = `gas_wanted` \* `gas_price`
  * Is the extra fee that is unused returned?
    * No, all fee suggested by user is charged.

## Decision

Proposals:

* Gas estimation:
  * Request gas esctimation from node
  * Multiply response by `safety coefficient`
* Fee prices estimations:
  * Option 1: Set recommended gas price for the network, embed it into the applications
  * Option 2: Dynamically determine the gas price based on recent transactions
    * Can be either implemented on client size or provided as a service
  * Then use exponential growth if the fee isn't enough

Transaction sending algorithm for VDR tools library consumer:

1. Build and sign a transaction
2. Send gas estimation request
3. Set initial gas price to a value: 1. Either proposed by the community; 2. Or retrieved from recent transactions, median for example: 1. Can be retrieved on a client-side; 2. Or received from public service.
4. Try to send the transaction
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

* Github issues with suggestions to improve fee mechanism:
  * [https://github.com/cosmos/cosmos-sdk/issues/6555](https://github.com/cosmos/cosmos-sdk/issues/6555)
  * [https://github.com/cosmos/cosmos-sdk/issues/2150](https://github.com/cosmos/cosmos-sdk/issues/2150)
  * [https://github.com/cosmos/cosmos-sdk/issues/4938](https://github.com/cosmos/cosmos-sdk/issues/4938)

