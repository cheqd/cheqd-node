# Account and key management in VDR Tools SDK

## Overview

This page describes how [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) works with cheqd accounts in identity wallets built using the VDR Tools SDK.

## Identity wallet key methods for cheqd accounts

These methods below are used to manage [cheqd accounts/wallets](../cheqd-cli/cheqd-cli-accounts.md) within identity wallets built using VDR Tools SDK. (For example, to pay for or receive payment for transactions that need to be written to ledger.)

For compatibility purposes, VDR Tools SDK method names use the `indy_` prefix. This may be updated in the future as work is done on the upstream project to refactor method names to be ledger-agnostic.

### indy_cheqd_keys_add_random

Create a new identity wallet key with specified `alias`, without specifying any other additional information, such as mnemonics.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.
* `alias` (string): Memorable key name/alias that makes it easier to reference for a user.

#### Example output

```jsonc
{
    "alias": "some_alias",
    "account_id": "cheqd1gudhsalrhsurucr5gkvga5973359etv6q0xvwz",
    "pub_key": "xSMzGopEnnTPCwjQwryDcrG9MGw3sVyb4ecYVaJrfkoA"
}
```

### indy_cheqd_keys_add_from_mnemonic

Similar to the technique used in [cheqd CLI to recover a key](../cheqd-cli/cheqd-cli-key-management.md), it allows identity wallets that use VDR Tools SDK to recover a key from mnemonic.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.
* `alias` (string): Memorable key name/alias that makes it easier to reference for a user.
* `mnemonic` (string): 12 or 24 word seed phrase to recover the keys associated with a wallet.

#### Example output

```jsonc
{
	"alias": "some_alias_2",
    "account_id": "cheqd10hcwm576uprz53wj2p8vv2dg0u8zu3n6l0wsxr",
    "pub_key": "fPn3LGakGrbHJTEk5fs7hAfa65DfpefgWawmwfCwjTHF"
}
```

### indy_cheqd_keys_get_info

Display information about an existing key based on its `alias` in an identity wallet.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.
* `alias` (string): Memorable key name/alias that makes it easier to reference for a user.

#### Example output

```jsonc
{
    "alias": "some_alias",
    "account_id": "cheqd17t7fmt3vpkkxa04hql0gyx8dufumq05vr9ztp6",
    "pub_key": "juXSChod4MmAAniezu44pDtTMgTUizUi84RaWqnhf43j"
}
```

### indy_cheqd_keys_get_list_keys

Return a list of all keys stored in the local identity wallet.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.

#### Example output

```js
 [
 Object({"account_id": String("cheqd1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0"), "alias": String("alice"), "pub_key": String("fTsZShn9KkgYKyDmbP5bLhVucNuPRdo4N6zGjAfzSSgv")}), 
 Object({"account_id": String("cheqd1c4n6j030trrljqsphmw9tcrcpgdf33hd3jd0jn"), "alias": String("some_alias_2"), "pub_key": String("g2vxGLkuYg84s3UcsKSxSttNgCKoQgRBQXizzSqbHdRJ")}), 
 Object({"account_id": String("cheqd1kcwadpmfreuvdrvkz7v79ydclfnn4ukdhp57c2"), "alias": String("some_alias_1"), "pub_key": String("27taHHmKLcxZEHQPmiNHuXTXQYF7u2CxGrzKQMxBZALsQ")})
 ]
 ```

### indy_cheqd_keys_sign

Sign a transaction using key in identity wallet with specified `alias`. The cheqd account associated with the key must have sufficient balance to fund the transaction.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.
* `alias` (string): Memorable key name/alias that makes it easier to reference for a user.
* `tx_raw`: Raw bytecode representation of of a correctly formattted cheqd/Cosmos transaction.
* `tx_len`: Length of string for bytes in transaction.

#### Example output

```text
[10, 146, 1, 10, 134, 1, 10, 37, 47, 99, 104, 101, 113, 100, 105, 100, 46, 99, 104, 101, 113, 100, 110, 111, 100, 101, 46, 99, 104, 101, 113, 100, 46, 77, 115, 103, 67, 114, 101, 97, 116, 101, 78, 121, 109, 18, 93, 10, 45, 99, 111, 115, 109, 111, 115, 49, 120, 51, 51, 120, 107, 106, 100, 51, 103, 113, 108, 102, 104, 122, 53, 108, 57, 104, 54, 48, 109, 53, 51, 112, 114, 50, 109, 100, 100, 52, 121, 51, 110, 99, 56, 54, 104, 48, 18, 10, 116, 101, 115, 116, 45, 97, 108, 105, 97, 115, 26, 11, 116, 101, 115, 116, 45, 118, 101, 114, 107, 101, 121, 34, 8, 116, 101, 115, 116, 45, 100, 105, 100, 42, 9, 116, 101, 115, 116, 45, 114, 111, 108, 101, 18, 4, 109, 101, 109, 111, 24, 192, 2, 18, 97, 10, 78, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 59, 126, 95, 52, 102, 213, 99, 251, 102, 62, 148, 101, 72, 226, 188, 243, 222, 31, 35, 148, 19, 127, 79, 75, 79, 37, 160, 132, 193, 33, 148, 7, 18, 4, 10, 2, 8, 1, 18, 15, 10, 9, 10, 4, 99, 104, 101, 113, 18, 1, 48, 16, 224, 167, 18, 26, 64, 130, 229, 164, 76, 214, 244, 157, 39, 135, 11, 118, 223, 29, 196, 41, 92, 247, 126, 129, 194, 18, 154, 136, 165, 153, 76, 202, 85, 187, 195, 40, 69, 10, 206, 165, 238, 223, 245, 35, 140, 92, 123, 246, 110, 23, 39, 32, 215, 239, 230, 196, 146, 168, 5, 147, 9, 67, 113, 242, 163, 0, 223, 233, 73]
```
