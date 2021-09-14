# ADR 008: Estimating transaction fees

## Status

| Category | Status |
| :--- | :--- |
| **ADR Stage** | DRAFT |
| **Implementation Status** | Not Implemented |

## Summary

The goal of this ADR is to define how fees will be defined and processed on the cheqd network.

A linked ADR associated with this is [ADR 006: **Community tax**](adr-006-community-tax.md), which describes additional steps that take place based on the fees for a particular transaction on the network.

## Context

To do understand how fees would be implemented on the cheqd network, it is important to understand what capabilities [Tendermint](https://docs.tendermint.com/master/) and [Cosmos SDK](https://docs.cosmos.network/v0.43/basics/gas-fees.html) provide to set and implement transaction fees.

From [Cosmos SDK's **Introduction to `Gas` and `Fees`**](https://docs.cosmos.network/v0.43/basics/gas-fees.html#introduction-to-gas-and-fees):

> In the Cosmos SDK, **`gas` is a special unit that is used to track the consumption of resources during execution**. `gas` is typically consumed whenever read and writes are made to the store, but it can also be consumed if expensive computation needs to be done. It serves two main purposes:

> * **Make sure blocks are not consuming too many resources and will be finalized**. This is implemented by default in the SDK via the [block gas meter](https://docs.cosmos.network/v0.43/basics/gas-fees.html#block-gas-meter).
> * **Prevent spam and abuse from end-user**. To this end, `gas` consumed during [`message`](https://docs.cosmos.network/v0.43/building-modules/messages-and-queries.html#messages) execution is typically priced, resulting in a `fee` \(`fees = gas * gas-prices`\). `fees` generally have to be paid by the sender of the `message`. Note that the SDK does not enforce `gas` pricing by default, as there may be other ways to prevent spam \(e.g. bandwidth schemes\). Still, most applications will implement `fee` mechanisms to prevent spam. This is done via the [`AnteHandler`](https://docs.cosmos.network/v0.43/basics/gas-fees.html#antehandler).

### Investigations related to `gas` and `fees` parameters

#### `gas-wanted`

`gas-wanted` is the amount of gas that a transaction is allowed to use and is set by the transaction _sender_.

While there is an estimation functionality in Cosmos CLI to predict what the gas needed for a transaction should be, it doesn't take into account reads and writes to state. As a consequence, it can be inaccurate and result in a scenario where `gas-wanted` is underestimated.

The `gas-wanted` estimation functionality can also be accessed via RPC calls by querying the `app/simulate` endpoint.

#### **`gas-used`**

`gas-used` is the amount of gas used by a transaction and is always less than  `gas-wanted`.

If a transaction requires more gas than specified in `gas-wanted`, execution is interrupted. The fee is NOT returned and becomes fully exhausted.

#### **`gas-prices`**

`gas-prices` is the price of gas suggested by transaction sender. It is set by the transaction sender and can be set to any value.

Validators can define a parameter called `min-gas-prices`, which acts as a filter on the transactions in each validators' mempool. Transactions with lower `gas-prices` are not considered as candidates for inclusion into the next blocks by the validator. This value is configured by the node/validator operator in the `app.toml` config file.

If `gas-prices` specified in a transaction is lower than the `min-gas-prices` of _all_ validators, no validators consider that transaction in their mempool. The transaction eventually times out and is not processed.

If `gas_prices` specified in a transaction is lower than `min-gas-prices` of _some_ validators, the transaction has a chance at being committed within the transaction timeout window. Whether this actually happens depends on the sum of voting powers of validators that are ready to commit it.

There is currently no way to request `min-gas-prices`.

#### **`fee`**

`fee` is the amount of tokens a validator takes for a transaction processing and is calculated as `fee` = `gas` \* `gas-prices`

`gas` above must be at least equal to or more than `gas-wanted`. 

### Options considered for estimating transaction fees

#### Gas estimation

**Option 1: Client-side calculation**

1. A gas estimation request could be built into the client libraries \(such as [VDR Tools](https://gitlab.com/evernym/verity/vdr-tools)\) used in applications using cheqd network
2. The client library would use a precalculated adjustment coefficient to increase the chances of the transaction being committed.

**Option 2: Use estimated values**

1. Estimated gas for most common transactions could be pre-calculated based on best-guess.
2. This method can be used as a workaround in the absence of a formal gas estimation request in Cosmos.
3. This approach can be used because gas needed to process simple transactions will remain constant over time.

#### Fee price estimations

**Option 1: Use fixed values**

1. Set a fixed recommended gas price on the cheqd network
2. Embed the fixed value into applications
3. This scenario is more likely to fail over time because validator operators can change gas prices.

**Option 2: Dynamically calculate values**

1. Gas prices can be dynamically determined based on recent transactions.
2. Possible ways of estimation:
   1. Moving average
   2. Weighed moving average
   3. Exponential moving average
3. Possible ways of implementation:
   1. Client side
   2. Provided as a service
   3. On chain (the most trusted and robust way)
3. To prevent constant growth of the mean value it's proposed to use multiplier `Y < 1` to get star fee price.
4. If the dynamically calculated fees are insufficient to put a transaction though, use exponential growth to calculate higher values and retry.

### Proposed transaction sending flow in VDR Tools

1. Build and sign a transaction
2. Send gas estimation request
   1. Multiply response to the precalculated coefficient `X > 1` to increase probability of the transaction being committed
3. Set initial gas price
   1. Request moving average from the network
   2. Multiply response by `Y < 1` to prevent constant growth of the mean value
4. Try sending the transaction to the network
   1. Exponentially increase the `gas` limit, in case of failure due to gas being lower than gas-wanted.
   2. Exponentially increase the `gas-prices`, in case transaction time-out \(as this indicates `min-gas-prices` was not met\)

### Proposed API in VDR Tools

#### Build simulate query

Used to build gas estimation request. 

```
indy_cheqd_ledger_tx_build_query_simulate(tx: byte[]) -> String
```

Parameters:
* `tx` - transaction bytes from `build_tx` method

Result:
* request body

`indy_cheqd_pool_abci_query` should be used to send the request and receive a response.

#### Parse simulate query response

Used to parse gas estimation response.

```
indy_cheqd_ledger_tx_parse_query_simulate_resp(resp: String) -> String
```

Parameters:
* `resp` - response string from `indy_cheqd_pool_abci_query` call

Result is json encoded  structure:

```
SimulateResponse {
   GasInfo {
      GasWanted uint64,
      GasUsed uint64
   },
   Result {
      ...
   }
}
```

Full type description can be found [here](https://github.com/cosmos/cosmos-sdk/blob/master/proto/cosmos/tx/v1beta1/service.proto#L112).

`GasUsed` is the estimated gas value.

#### Get gas multiplier

It's recommended to use this precomputed multiplier to increase chances of a transaction being committed.

```
indy_cheqd_ledger_tx_get_gas_multiplier() -> double
```

Params:
* None

Result:
* Multiplier value

#### Build average gas price query

```
indy_cheqd_ledger_tx_build_query_prices() -> String
```

Parameters:
* None

Result:
* Request body

`indy_cheqd_pool_abci_query` should be used to send the request and receive a response.

#### Parse average gas price response

Used to parse gas estimation response.

```
indy_cheqd_ledger_tx_parse_query_prices_resp(resp: String) -> String
```
Params:
* None

Return string is json encoded structure:

```
PricesResponse {
   LastAverage uint64
   SubsetSize uint64
}
```

`LastAverage` - average gas price for the recent `N` committed transactions where `N` is `SubsetSize`.

#### Get gas average price multiplier

To minimize `gas price` it's recommended to apply a multiplier for `last average`.

```
indy_cheqd_ledger_tx_get_gas_price_multiplier() -> double
```

Parameters:
* None

Result:
* Multiplier value

## Decision



## Consequences

There is no solution currently available that simplifies transaction fee estimation beyond the mechanisms already by Cosmos which are not client side.

### Backward Compatibility

* This proposal is compatible with all recent versions of Cosmos.

### Positive

* 
### Negative

* Client-side library implementation for gas/fee estimation can be complex to achieve

### Neutral

* 
## References

* [ADR 006: Community Tax](adr-006-community-tax.md)
* [Tendermint Core documentation](https://docs.tendermint.com/master/)
* [Cosmos SDK documentation](https://docs.cosmos.network/)
* [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) \(client side library for cheqd\)
* Open issues on Cosmos SDK Github on fee estimation improvements:
  * [\#2150 Refunding unused but allocated gas](https://github.com/cosmos/cosmos-sdk/issues/2150)
  * [\#4938 Tx Gas Estimation Improvement](https://github.com/cosmos/cosmos-sdk/issues/4938)
  * [\#6555 More robust gas pricing](https://github.com/cosmos/cosmos-sdk/issues/6555)
  * [\#9569 Refund Unused Fee](https://github.com/cosmos/cosmos-sdk/issues/9569)

