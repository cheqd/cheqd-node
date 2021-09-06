---
description: This is the suggested template to be used for ADRs on the cheqd-node project.
---

# ADR 005: Genesis parameters

## Status

PROPOSED \| Not Implemented

## Summary

The aim of this document is to define genesis parameters that will be used in testnets and the mainnet.

**NB: v0.42.5**

## Context

Genesis consists of Tendermint consensus engine parameters and Cosmos app-specific parameters.

### Consensus parameters

Proposed values:

* block
  * max\_bytes = `22020096` \(~22MB\)
    * Cosmos hub: `200000` \(~200KB\)
  * max\_gas = `-1` \(no gas limit\)
    * Cosmos hub: `2000000`
  * \(?\) time\_iota\_ms = `1000`
    * Cosmos hub: `1000`
* \(?\) evidence
  * max\_age\_num\_blocks = `100000`
  * max\_age\_duration = `172800000000000`
  * max\_bytes = `1048576`
* validator
  * pub\_key\_types = `[ "ed25519" ]`

### Application parameters

Cosmos application is divided into modules. Each module has parameters that help to adjust the module's behavior. Here are proposed values for default modules:

* auth
  * max\_memo\_characters = `512`
    * Max number of characters in the memo field
  * tx\_sig\_limit = `7`
    * Max number of signatures
  * tx\_size\_cost\_per\_byte = `10`
    * Gas cost of transaction byte
  * sig\_verify\_cost\_ed25519 = `590`
    * Cost of `ed25519` signature verification
  * sig\_verify\_cost\_secp256k1 = `1000`
    * Cost of `secp256k1` signature verification
* \(?\) bank
  * "send\_enabled": \[\],

    "default\_send\_enabled": true
* \(?\) capability
  * ?
* \(?\) cheqd
  * None
* \(?\) crisis
  * ?
* distribution
  * community\_tax = `0.02`
    * The percent of rewards that goes to community fund pool
  * "base\_proposer\_reward": "0.010000000000000000",
  * "bonus\_proposer\_reward": "0.040000000000000000",
  * \(?\) withdraw\_addr\_enabled = true
* \(?\) evidence
  * ?
* \(?\) genutil
  * None
* \(?\) gov
  * "deposit\_params": {

      "min\_deposit": \[

        {

          "denom": "stake",

          "amount": "10000000"

        }

      \],

      "max\_deposit\_period": "172800s"

    },

    "voting\_params": {

      "voting\_period": "172800s"

    },

    "tally\_params": {

      "quorum": "0.334000000000000000",

      "threshold": "0.500000000000000000",

      "veto\_threshold": "0.334000000000000000"

    }
* \(?\) ibc
  * ?
* mint
  * mint\_denom = `cheq`
  * inflation\_rate\_change = `0.13`
    * Max inflation rate change per year
    * In Cosmos hub they use `1.0`
    * Formula: `inflationRateChangePerYear = (1 - BondedRatio/ GoalBonded) * MaxInflationRateChange`
  * inflation\_max = `0.20`
    * Inflation aims to this value if `bonded_ratio` &lt; `bonded_goal`
    * Cosmos hub: `0.20`
  * inflation\_min = `0.07`
    * Inflation aims to this value if `bonded_ratio` &lt; `bonded_goal`
    * Cosmos hub: `0.07`
  * goal\_bonded = `0.67`
    * Cosmos hub: `0.67`
  * blocks\_per\_year = `6311520`
    * Cosmos hub: `4360000`
* slashing
  * signed\_blocks\_window = `120960` \(1 week\)
    * Cosmos hub: `10000` \(~20h\)
  * min\_signed\_per\_window = `0.50`
    * This percentage of blocks must be signed within the window
  * downtime\_jail\_duration = `600s`
    * The minimal time validator have to stay in jail
  * slash\_fraction\_double\_sign = `0.05`
    * Slash for double sign
  * slash\_fraction\_downtime = `0.01`
    * Slash for downtime
* staking
  * unbonding\_time = `1814400s`
    * A delegator must wait this time before tokens become unbonded
  * max\_validators = `125`
    * The maximum number of validators in the network
  * max\_entries = `7`
    * Max amount of unbound/redelegation operations in progress per account
  * historical\_entries = `10000`
    * Amount of unbound/redelegate entries to store
  * bond\_denom = `stake`
    * Denomination used in staking
* \(?\) transfer
  * "send\_enabled": true,

    "receive\_enabled": true

### Parameter adjustment

All parameters can be changed via change proposals + voting.

## Decision

What is the change that we're proposing and/or doing?

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Backward Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## References

* {reference link}

