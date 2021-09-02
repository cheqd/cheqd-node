# ADR 001: Fee processing

## Status

DRAFT | Not Implemented

## Summary

The goals of this ARD is to define how fees will be processed on cheqd network. To do it we firstly need to understand what capabilities Tendermint and Cosmos SDK provide.

## Context

Gas and fees in Cosmos SDK and Tendermint:

- Gas:
  - What is it?
    - It's a unit used to track the consumption of resources during transaction execution.
  - What is it used for?
    - Block computational complexity limitation;
    - Spam prevention.
- Gas wanted:
  - What is it?
    - The amount of gas that tx is allowed to use.
    - Set by transaction sender
  - How it can be predicted?
    - __TODO:__ simulations
- Gas used:
  - What is it?
    - The amount of gas used by a transaction.
    - Always less than `gas_wanted`.
  - What happens if a transaction requires more gas then specified in `gas_wanted`?
    - Execution interrupts
    - Fee isn't returned (fully charged)
- Gas prices:
  - What is transaction's `gas_prices`?
    - It's prises of gas suggested by transaction sender
    - Set by transaction sender
    - Can be any value
  - What is validator's `min_gas_prices`?
    - It's a filter used to include transactions in validator's mempool.
    - Transactions with less `gas_prices` are not considered as candidates for inclusion into next blocks.
    - Specific for each validator.
    - Set by validator operator in `app.toml`.
  - What happens if `gas_price` is less then `min_gas_prices` of all validators?
    - No validators include this transaction in mempool, it times out.
  - What happens if `gas_price` is less then `min_gas_prices` of some validators?
    - The transaction has chances to be committed within timeout.
    - It it depends on sum of voting powers of validators that are ready to commit it.
  - Can `gas_prices` be requested?
    - __TODO:__ most likely no
- Fee:
  - What is it?
    - The amount of tokens validator takes for a transaction processing.
  - How is it calculated?
    - `fee` = `gas_wanted` * `gas_price`
  - Is extra fee that is unused returned?
    - No, all fee suggested by user is charged.


How other notworks handle fees? Here are some thoughts:

- https://github.com/cosmos/cosmos-sdk/issues/6555
- https://github.com/cosmos/cosmos-sdk/issues/2150
- https://github.com/cosmos/cosmos-sdk/issues/4938

> This section describes the forces at play, such as business, technological, social, and project local. These forces are probably in tension, and should be called out as such. The language in this section is value-neutral. It is simply describing facts. It should clearly explain the problem and motivation that the proposal aims to resolve.

## Decision

What is the change that we're proposing and/or doing?

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## References

- {reference link}