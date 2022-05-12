# ADR 005: Genesis parameters

## Status

| Category                  | Status           |
| ------------------------- | ---------------- |
| **Authors**               | Alexandr Kolesov |
| **ADR Stage**             | ACCEPTED         |
| **Implementation Status** | Implemented      |
| **Start Date**            | 2021-09-15       |

## Summary

The aim of this document is to define the genesis parameters that will be used in cheqd network testnet and mainnet.

> Cosmos v0.44.3 parameters are described.

## Context

Genesis consists of Tendermint consensus engine parameters and Cosmos app-specific parameters.

### Consensus parameters

Tendermint requires [genesis parameters](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis) to be defined for basic consensus conditions on any Cosmos network.

#### Block parameters

| Parameter      | Description                                                                                                                             | Mainnet            | Testnet            |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------- | ------------------ | ------------------ |
| `max_bytes`    | This sets the maximum size of total bytes that can be committed in a single block. This should be larger than `max_bytes for evidence.` | 200,000 (\~200 KB) | 200,000 (\~200 KB) |
| `max_gas`      | This sets the maximum gas that can be used in any single block.                                                                         | 200,000 (\~200 KB) | 200,000 (\~200 KB) |
| `time_iota_ms` | Unused. This has been deprecated and will be removed in a future version of Cosmos SDK.                                                 | 1,000 (1 second)   | 1,000 (1 second)   |

#### Evidence parameters

| Parameter            | Description                                                                                                                                              | Mainnet                                                                                          | Testnet                                                                        |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `max_age_num_blocks` | Maximum age of evidence, in blocks. The basic formula for calculating this is: `MaxAgeDuration / {average block time}`.                                  | **12,100**                                                                                       | 25,920                                                                         |
| `max_age_duration`   | Maximum age of evidence, in time. It should correspond with an app's "unbonding period".                                                                 | <p><strong>1,209,600,000,000,000</strong> </p><p></p><p>(expressed in nanoseconds, ~2 weeks)</p> | <p>259,200,000,000,000 </p><p></p><p>(expressed in nanoseconds, ~72 hours)</p> |
| `max_bytes`          | This sets the maximum size of total evidence in bytes that can be committed in a single block and should fall comfortably under `max_bytes` for a block. | **50,000** (\~ 50 KB)                                                                            | 5,000 (\~ 5 KB)                                                                |

#### Validator

| Parameter       | Description                                                   | Mainnet | Testnet |
| --------------- | ------------------------------------------------------------- | ------- | ------- |
| `pub_key_types` | Types of public keys supported for validators on the network. | Ed25519 | Ed25519 |

### Application parameters

Cosmos application is divided [into a list of modules](https://docs.cosmos.network/v0.44/modules/). Each module has parameters that help to adjust the module's behaviour. Here are proposed values for default modules:

#### `Auth`

| Parameter                   | Description                                    | Mainnet | Testnet |
| --------------------------- | ---------------------------------------------- | ------- | ------- |
| `max_memo_characters`       | Maximum number of characters in the memo field | 512     | 512     |
| `tx_sig_limit`              | Max number of signatures                       | 7       | 7       |
| `tx_size_cost_per_byte`     | Gas cost of transaction byte                   | 10      | 10      |
| `sig_verify_cost_ed25519`   | Cost of `ed25519` signature verification       | 590     | 590     |
| `sig_verify_cost_secp256k1` | Cost of `secp256k1` signature verification     | 1,000   | 1,000   |

#### Bank

| Parameter              | Description                                                                     | Mainnet | Testnet |
| ---------------------- | ------------------------------------------------------------------------------- | ------- | ------- |
| `default_send_enabled` | The default send enabled value allows send transfers for all coin denominations | True    | True    |

#### Crisis

| Parameter      | Description                                                                                                                             | Mainnet                                                              | Testnet                                                              |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `constant_fee` | The fee is used to verify the [invariant(s)](https://docs.cosmos.network/v0.44/building-modules/invariants.html) in the `crisis` module | <p>10,000,000,000,000 <em>ncheq</em> </p><p></p><p>(10,000 CHEQ)</p> | <p>10,000,000,000,000 <em>ncheq</em> </p><p></p><p>(10,000 CHEQ)</p> |

#### Distribution

| Parameter               | Description                                                                                                          | Mainnet   | Testnet   |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------- | --------- | --------- |
| `community_tax`         | The percent of rewards that goes to the community fund pool                                                          | 0.02 (2%) | 0.02 (2%) |
| `base_proposer_reward`  | Base reward that proposer of a block receives                                                                        | 0.01 (1%) | 0.01 (1%) |
| `bonus_proposer_reward` | Bonus reward that proposer gets for proposing block. This depends on the number of pre-commits included to the block | 0.04 (4%) | 0.04 (4%) |
| `withdraw_addr_enabled` | Whether withdrawal address can be changed or not. By default, it's the delegator's address.                          | True      | True      |

#### `Gov`

| Parameter            | Description                                                                                         | Mainnet                                                             | Testnet                                                             |
| -------------------- | --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------- |
| **`deposit_params`** |                                                                                                     | ``                                                                  | ``                                                                  |
| `min_deposit`        | The minimum deposit for a proposal to enter the voting period.                                      | \[{ "denom": "ncheq", "amount": "8,000,000,000,000" }] (8,000 cheq) | \[{ "denom": "ncheq", "amount": "8,000,000,000,000" }] (8,000 cheq) |
| `max_deposit_period` | The maximum period for Atom holders to deposit on a proposal. Initial value: 2 months.              | 604,800s (1 week)                                                   | **172,800s (48 hours)**                                             |
| **`voting_params`**  |                                                                                                     | ``                                                                  | ``                                                                  |
| `voting_period`      | The defined period for an on-ledger vote from start to finish.                                      | 604,800s (1 week)                                                   | **172,800s (48 hours)**                                             |
| **`tally_params`**   |                                                                                                     |                                                                     |                                                                     |
| `quorum`             | Minimum percentage of total stake needed to vote for a result to be considered valid.               | 0.334 (33.4%)                                                       | 0.334 (33.4%)                                                       |
| `threshold`          | Minimum proportion of Yes votes for proposal to pass.                                               | 0.5 (50%)                                                           | 0.5 (50%)                                                           |
| `veto_threshold`     | The minimum value of veto votes to total votes ratio for proposal to be vetoed. Default value: 1/3. | 0.334 (33.4%)                                                       | 0.334 (33.4%)                                                       |

#### `Mint`

| Parameter               | Description                                                                                                                                                                                                                  | Mainnet                                | Testnet                                |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------- | -------------------------------------- |
| `mint_denom`            | Name of the cheq smalles denomination                                                                                                                                                                                        | ncheq                                  | ncheq                                  |
| `inflation_rate_change` | <p>Maximum inflation rate change per year. </p><p></p><p>In Cosmos Hub they use <code>1.0</code></p><p></p><p>Formula: <code>inflationRateChangePerYear = (1 - BondedRatio / GoalBonded) * MaxInflationRateChange</code></p> | 0.045 (4.5%)                           | 0.045 (4.5%)                           |
| `inflation_max`         | Inflation aims to this value if `bonded_ratio` < `bonded_goal`                                                                                                                                                               | 0.04 (4%)                              | 0.04 (4%)                              |
| `inflation_min`         | Inflation aims to this value if `bonded_ratio` > `bonded_goal`                                                                                                                                                               | 0.01 (1%)                              | 0.01 (1%)                              |
| `goal_bonded`           | Percentage of bonded tokens at which inflation rate will neither increase nor decrease                                                                                                                                       | 0.60 (60%)                             | 0.60 (60%)                             |
| `blocks_per_year`       | Number of blocks generated per year                                                                                                                                                                                          | 3,155,760 (1 block every \~10 seconds) | 3,155,760 (1 block every \~10 seconds) |

#### `Slashing`

| Parameter                    | Description                                                         | Mainnet                                                              | Testnet                                                                  |
| ---------------------------- | ------------------------------------------------------------------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `signed_blocks_window`       | Number of blocks a validator can miss signing before it is slashed. | 25,920 (expressed in blocks, equates to 259,200 seconds or \~3 days) | **17,280 (expressed in blocks, equates to 172,800 seconds or \~2 days)** |
| `min_signed_per_window`      | This percentage of blocks must be signed within the window.         | 0.50 (50%)                                                           | 0.50 (50%)                                                               |
| `downtime_jail_duration`     | The minimal time validator have to stay in jail                     | 600s (\~10 minutes)                                                  | 600s (\~10 minutes)                                                      |
| `slash_fraction_double_sign` | Slashed amount as a percentage for a double sign infraction         | 0.05 (5%)                                                            | 0.05 (5%)                                                                |
| `slash_fraction_downtime`    | Slashed amount as a percentage for downtime                         | 0.01 (1%)                                                            | 0.01 (1%)                                                                |

#### `Staking`

| Parameter            | Description                                                                    | Mainnet                | Testnet                 |
| -------------------- | ------------------------------------------------------------------------------ | ---------------------- | ----------------------- |
| `unbonding_time`     | A delegator must wait this time after unbonding before tokens become available | 1,210,000s (\~2 weeks) | **259,200s (\~3 days)** |
| `max_validators`     | The maximum number of validators in the network                                | 125                    | 125                     |
| `max_entries`        | Max amount of unbound/redelegation operations in progress per account          | 7                      | 7                       |
| `historical_entries` | Amount of unbound/redelegate entries to store                                  | 10,000                 | 10,000                  |
| `bond_denom`         | Denomination used in staking                                                   | ncheq                  | ncheq                   |

#### `ibc`

| Parameter                     | Description                                                                                                                                                                                                                                                                                                                                    | Mainnet                                                  | Testnet                                                  |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| `max_expected_time_per_block` | Maximum expected time per block, used to enforce block delay. This parameter should reflect the largest amount of time that the chain might reasonably take to produce the next block under normal operating conditions. A safe choice is 3-5x the expected time per block.                                                                    | 30,000,000,000 (expressed in nanoseconds, \~ 30 seconds) | 30,000,000,000 (expressed in nanoseconds, \~ 30 seconds) |
| `allowed_clients`             | Defines the list of allowed client state types. We allow connections from other chains using the [Tendermint client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint), and with light clients using the [Solo Machine client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/06-solomachine). | \[ "06-solomachine", "07-tendermint" ]                   | \[ "06-solomachine", "07-tendermint" ]                   |

#### `ibc-transfer`

| Parameter         | Description                                                         | Mainnet | Testnet |
| ----------------- | ------------------------------------------------------------------- | ------- | ------- |
| `send_enabled`    | Enables or disables all cross-chain token transfers from this chain | true    | true    |
| `receive_enabled` | Enables or disables all cross-chain token transfers to this chain   | true    | true    |

## Decision

The parameters above were agreed separate the cheqd mainnet and testnet parameters. We have **bolded** the testnet parameters that differ from mainnet.

## Consequences

### Backward Compatibility

* The token denomination has been changed to make the smallest denomination 10^-9 `cheq` instead of 1 `cheq`. This is a breaking change from the previous version of the cheqd testnet that will potentially require new tokens to be transferred and issued to testnet node operators.

### Positive

* Inflation allows fees to be collected from block rewards in addition to transaction fees.
* In production/mainnet, parameters can only be changed via a majority vote without veto defeat according to the cheqd network governance principles. This allows for more democratic governance frameworks to be created for a self-sovereign identity network.

### Negative

* Existing node operators will need to re-establish staking with new staking denomination and staking parameters.

### Neutral

* Voting time, unbonding period, and deposit period have all been reduced to 2 weeks to balance the speed at which decisions can be reached vs giving enough time to validators to participate.

## References

* [List of Cosmos modules](https://docs.cosmos.network/v0.44/modules/)
* [Tendermint genesis parameters](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis)
