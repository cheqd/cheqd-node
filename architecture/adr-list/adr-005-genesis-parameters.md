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

Genesis consists of Tendermint consensus engine parameters and Cosmos app specific parameters.

### Consensus parameters

Proposed values:

* block
  * max\_bytes = `22020096`
  * max\_gas = `-1`
    * No gas limit
  * time\_iota\_ms = `1000`
* evidence
  * max\_age\_num\_blocks = `100000`
  * max\_age\_duration = `172800000000000`
  * max\_bytes = `1048576`
* validator
  * pub\_key\_types = `[ "ed25519" ]`

### Application parameters

Cosmos application is divided into modules. Each module have parameters that help to adjust module's behavior. Here are proposed values for default modules:

* auth
  * max\_memo\_characters = 256
  * tx\_sig\_limit = 7
  * tx\_size\_cost\_per\_byte = 10
  * sig\_verify\_cost\_ed25519 = 590
  * sig\_verify\_cost\_secp256k1 = 1000
* bank
  * "send\_enabled": \[\],

    "default\_send\_enabled": true
* capability
  * ?
* cheqd
  * None
* crisis
  * ?
* distribution
  * "community\_tax": "0.020000000000000000",

    "base\_proposer\_reward": "0.010000000000000000",

    "bonus\_proposer\_reward": "0.040000000000000000",

    "withdraw\_addr\_enabled": true
* evidence
  * ?
* genutil
  * None
* gov
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
* ibc
  * ?
* mint
  * "mint\_denom": "stake",

    "inflation\_rate\_change": "0.130000000000000000",

    "inflation\_max": "0.200000000000000000",

    "inflation\_min": "0.070000000000000000",

    "goal\_bonded": "0.670000000000000000",

    "blocks\_per\_year": "6311520"
* slashing
  * "signed\_blocks\_window": "100",

    "min\_signed\_per\_window": "0.500000000000000000",

    "downtime\_jail\_duration": "600s",

    "slash\_fraction\_double\_sign": "0.050000000000000000",

    "slash\_fraction\_downtime": "0.010000000000000000"
* staking
  * "unbonding\_time": "1814400s",

    "max\_validators": 100,

    "max\_entries": 7,

    "historical\_entries": 10000,

    "bond\_denom": "stake"
* transfer
  * "send\_enabled": true,

    "receive\_enabled": true

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

* {reference link}

