# ADR 008: Importing/exporting cheqd account mnemonics in VDR Tools SDK

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Andrew Nikitn, Ankur Banerjee |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Implemented |
| **Start Date** | 2021-09-23 |

## Summary

This ADR describes how cheqd/Cosmos account keys can be imported/exported into identity wallet applications built on Evernym VDR Tools SDK.

## Context

Client SDK applications such as [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) need to work with cheqd accounts in identity wallets to be able to interact with the cheqd network ledger.

For example, an identity wallet application or backend application would need to pay network transaction fees for writing [cheqd DIDs to the ledger](adr-002-cheqd-did-method.md). This may also need to be extended in the future to support [peer-to-peer payments for credential exchange](adr-001-payment-mechanism-for-issuing-credentials.md).

### Assumptions / Considerations

Cosmos SDK uses [known algorithms for deriving private keys from mnemonics](https://docs.cosmos.network/master/basics/accounts.html#keyring). This can be replicated using standard crypto libraries to carry out the same steps as in Cosmos SDK:

```text
rounds of iteration :    2048
length              :    64
algorithm           :    sha512
salt                :    "mnemonic" + passphrase
```

The mnemonic above is assumed to be a pre-existing one [cheqd/Cosmos CLI](../../docs/cheqd-cli/cheqd-cli-accounts.md). The "passphrase" above is user-defined, and defaults to blank if not defined.

## Decision

Mnemonic import/export can be achieved using pre-existing [BIP39](https://github.com/bitcoin/bips/tree/master/bip-0039) packages and [Cosmos SDK's Rust library `cosmrs`](https://github.com/cosmos/cosmos-rust).

Using these pre-existing libraries, cheqd accounts can be recovered using the standard [BIP44 `HDPath`](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) for Cosmos SDK chains described below:

```text
"m/44'/118'/0'/0/0"
```

## Consequences

Functionality will be added to VDR Tools SDK to import/export cheqd accounts using mnemonics paired with the `--recover` flag as done with Cosmos wallets.

### Backwards Compatibility

Not applicable, since this is an entirely new feature in VDR Tools SDK for integration with the new blockchain framework.

### Positive

* Adding/recovering cheqd accounts in VDR Tools SDK will follow a similar, familiar process that users have for Cosmos wallets.

### Negative

N/A

### Neutral

N/A

## References

* [Cosmos SDK account generation and keyrings](https://docs.cosmos.network/master/basics/accounts.html)
