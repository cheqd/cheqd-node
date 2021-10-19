# ADR 005: Genesis parameters

## Status

| Category | Status |
| :--- | :--- |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Not Implemented |

## Summary

The aim of this document is to define the genesis parameters that will be used in cheqd network testnet and mainnet.

> Cosmos v0.42.5 parameters are described.

## Context

Genesis consists of Tendermint consensus engine parameters and Cosmos app-specific parameters.

### Consensus parameters

Tendermint requires [genesis parameters](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis) to be defined for basic consensus conditions on any Cosmos network.

#### Proposed values

* **`block`**
  * `max_bytes` = `200000` \(~200KB\)
    * Cosmos Hub: `200000` \(~200KB\)
  * `max_gas` = `2000000` \(~20 txs\)
    * Cosmos Hub: `2000000` \(~20 txs\)
* **`evidence`**
  * `max_age_num_blocks` = `1576800`
    * Maximum age of evidence, in blocks. The basic formula for calculating this is: `MaxAgeDuration / {average block time}`.
  * `max_age_duration` = `7884000`
    * Maximum age of evidence, in time. It should correspond with an app's "unbonding period".
  * `max_bytes` = `1048576`
    * This sets the maximum size of total evidence in bytes that can be committed in a single block and should fall comfortably under `max_bytes` for a block.
* **`validator`**
  * `pub_key_types` = `[ "ed25519" ]`

### Application parameters

Cosmos application is divided [into a list of modules](https://docs.cosmos.network/v0.44/modules/). Each module has parameters that help to adjust the module's behaviour. Here are proposed values for default modules:

* **`auth`**
  * `max_memo_characters` = `512`
    * Maximum number of characters in the memo field
  * `tx_sig_limit` = `7`
    * Max number of signatures
  * `tx_size_cost_per_byte` = `10`
    * Gas cost of transaction byte
  * `sig_verify_cost_ed25519` = `590`
    * Cost of `ed25519` signature verification
  * `sig_verify_cost_secp256k1` = `1000`
    * Cost of `secp256k1` signature verification
* **`bank`**
  * `send_enabled` = `[]`
    * Enables send for specific denominations
  * `default_send_enabled` = `true`
    * The default send enabled value allows send transfers for all coin denominations
* **`crisis`**
  * `constant_fee` = `{ "denom": "ncheq", "amount": "10000000000000" }` \(10,000 `cheq`\)
    * The fee is used to verify the [invariant\(s\)](https://docs.cosmos.network/v0.44/building-modules/invariants.html) in the `crisis` module.
* **`distribution`**
  * `community_tax` = `0.02`
    * The percent of rewards that goes to the community fund pool
  * `base_proposer_reward` = `0.01`
    * Base reward that proposer gets
  * `bonus_proposer_reward` = `0.04`
    * Bonus reward that proposer gets. This depends on the number of pre-commits included to the block
  * `withdraw_addr_enabled` = `true`
    * Whether withdrawal address can be changed or not. By default, it's the delegator's address.
* **`evidence`**
  * No parameters
* **`genutil`**
  * Used to manage initial transactions such as genesis validators creation
* **`gov`**
  * `deposit_params`
    * min\_deposit = `[{ "denom": "ncheq", "amount": "8000000000000" }]` \(8,000 `cheq`\)
      * The minimum deposit for a proposal to enter the voting period.
    * `max_deposit_period` = `1210000s` \(2 weeks\)
      * The maximum period for Atom holders to deposit on a proposal. Initial value: 2 months.
  * `voting_params`
    * voting\_period = `1210000s` \(2 weeks\)
  * `tally_params`
    * `quorum` = `0.334`
      * Minimum percentage of total stake needed to vote for a result to be considered valid. 
    * `threshold` = `0.5`
      * Minimum percentage of total stake needed to vote for a result to be considered valid.
    * `veto_threshold` = `0.334`
      * The minimum value of veto votes to total votes ratio for proposal to be vetoed. Default value: 1/3.
* **`mint`**
  * `mint_denom` = `ncheq`
  * `inflation_rate_change` = `0.02`
    * Maximum inflation rate change per year
    * In Cosmos Hub they use `1.0`
    * Formula: `inflationRateChangePerYear = (1 - BondedRatio/ GoalBonded) * MaxInflationRateChange`
  * `inflation_max` = `0.04`
    * Inflation aims to this value if `bonded_ratio` &lt; `bonded_goal`
    * Cosmos Hub: `0.20`
  * `inflation_min` = `0.01`
    * Inflation aims to this value if `bonded_ratio` &lt; `bonded_goal`
    * Cosmos Hub: `0.07`
  * `goal_bonded` = `0.60`
    * Cosmos Hub: `0.67`
  * `blocks_per_year` = `6311520` \(~5s\)
    * Cosmos Hub: `4360000`
* **`slashing`**
  * `signed_blocks_window` = `120960` \(1 week\)
    * Cosmos Hub: `10000` \(~20h\)
  * `min_signed_per_window`= `0.50`
    * This percentage of blocks must be signed within the window
  * `downtime_jail_duration` = `600s`
    * The minimal time validator have to stay in jail
  * `slash_fraction_double_sign` = `0.05`
    * Slash for double sign
  * `slash_fraction_downtime` = `0.01`
    * Slash for downtime
* **`staking`**
  * `unbonding_time` = `1814400s` \(3 weeks\)
    * A delegator must wait this time before tokens become unbonded
  * `max_validators` = `125`
    * The maximum number of validators in the network
  * `max_entries` = `7`
    * Max amount of unbound/redelegation operations in progress per account
  * `historical_entries` = `10000`
    * Amount of unbound/redelegate entries to store
  * `bond_denom` = `ncheq`
    * Denomination used in staking
* _\[ibc\]_ **`ibc`**
  * `max_expected_time_per_block` = `20000000000`
    * Maximum expected time per block (in nanoseconds), used to enforce block delay. This parameter should reflect the largest amount of time that the chain might reasonably take to produce the next block under normal operating conditions. A safe choice is 3-5x the expected time per block.
  * `allowed_clients` = `[ "07-tendermint" ]` (allow connections only from other chains)
    * Defines the list of allowed client state types. We allow connections only from other chains.
* _\[ibc\]_ **`capability`**
  * No params
* _\[ibc\]_ **`transfer`**
  * `send_enabled` = `true`
    * Enables or disables all cross-chain token transfers from this chain
  * `receive_enabled` = `true`
    * Enables or disables all cross-chain token transfers to this chain

## Decision

The parameters above were agreed to be used for the cheqd network testnet, with a view towards testing them for cheqd mainnet.

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

