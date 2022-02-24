# VDR Tools SDK ledger connection API

## Overview

This page describes the API for how [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) works with identity wallet keys and how it connects to the ledger "pool".

It is worth noting here that the terminology of "pool" connection is specifically a legacy term originally used in [Hyperledger Indy](https://github.com/hyperledger/indy-node), which as a permissioned blockchain assumes there is a finite pools of servers. While this paradigm is no longer true in the public, permissionless world of the cheqd network, the identity APIs in VDR Tools SDK and similar Hyperledger Aries-based frameworks is retained for explanations.

## Identity wallet key methods

For compatibility purposes, VDR Tools SDK method names use the `indy_` prefix. This may be updated in the future as work is done on the upstream project to refactor method names to be ledger-agnostic.

### indy_cheqd_keys_add_random

This method implements the logic of identity wallet key creation just using `alias`, without specifying any other additional information, like mnemonic.

#### Input parameters

* `wallet_handle` (integer): Linked to previously created and opened wallet.
* `alias` (string): Human-readable representation of alias to user,

#### Example output

```jsonc
{
    "alias": "some_alias",
    "account_id":"cheqd1gudhsalrhsurucr5gkvga5973359etv6q0xvwz",
    "pub_key":"xSMzGopEnnTPCwjQwryDcrG9MGw3sVyb4ecYVaJrfkoA"
}
```

### indy_cheqd_keys_add_from_mnemonic

This method realised logic for recovering keys from `mnemonic` string. Mnemonic string - it's a human-readable combination of words, in general can include 12 or 24 words.

Input parameters:

* `wallet_handle` - integer, which is connected to previously created and opened wallet.
* `alias` - human-readable representation of alias to user,
* `mnemonic` - string of 12 or 24 words.

As result, the next structure is expected:
```
{
    "alias":"some_alias_2",
    "account_id":"cosmos10hcwm576uprz53wj2p8vv2dg0u8zu3n6l0wsxr",
    "pub_key":"fPn3LGakGrbHJTEk5fs7hAfa65DfpefgWawmwfCwjTHF"
}
```

#### **indy_cheqd_keys_get_info**
This method is needed for getting information about already generated and stored keys by using only an `alias`.

Input parameters:

* `wallet_handle` - integer, which is connected to previously created and opened wallet.
* `alias` - human-readable representation of alias to user,

As result, the list of next structures is expected:
```
{
    "alias":"some_alias",
    "account_id":"cosmos17t7fmt3vpkkxa04hql0gyx8dufumq05vr9ztp6",
    "pub_key":"juXSChod4MmAAniezu44pDtTMgTUizUi84RaWqnhf43j"
}
```

#### **indy_cheqd_keys_get_list_keys**
The method returns a list of all keys which are placed locally.

Input parameters:

* `wallet_handle` - integer, which is connected to previously created and opened wallet.

As result, the next structure is expected:
```
 [
 Object({"account_id": String("cosmos1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0"), "alias": String("alice"), "pub_key": String("fTsZShn9KkgYKyDmbP5bLhVucNuPRdo4N6zGjAfzSSgv")}), 
 Object({"account_id": String("cosmos1c4n6j030trrljqsphmw9tcrcpgdf33hd3jd0jn"), "alias": String("some_alias_2"), "pub_key": String("g2vxGLkuYg84s3UcsKSxSttNgCKoQgRBQXizzSqbHdRJ")}), 
 Object({"account_id": String("cosmos1kcwadpmfreuvdrvkz7v79ydclfnn4ukdhp57c2"), "alias": String("some_alias_1"), "pub_key": String("27taHHmKLcxZEHQPmiNHuXTXQYF7u2CxGrzKQMxBZALsQ")})
 ]
 ```

#### **indy_cheqd_keys_sign**
This method can sign a transaction by using a key which can be found by `alias`

Input parameters:

* `wallet_handle` - integer, which is connected to previously created and opened wallet.
* `alias` - human-readable representation of alias to user, 
* `tx_raw` - byte representation of transaction,
* `tx_len` - length of string with bytes of transaction,

As result, raw byte's string is expected.
```
[10, 146, 1, 10, 134, 1, 10, 37, 47, 99, 104, 101, 113, 100, 105, 100, 46, 99, 104, 101, 113, 100, 110, 111, 100, 101, 46, 99, 104, 101, 113, 100, 46, 77, 115, 103, 67, 114, 101, 97, 116, 101, 78, 121, 109, 18, 93, 10, 45, 99, 111, 115, 109, 111, 115, 49, 120, 51, 51, 120, 107, 106, 100, 51, 103, 113, 108, 102, 104, 122, 53, 108, 57, 104, 54, 48, 109, 53, 51, 112, 114, 50, 109, 100, 100, 52, 121, 51, 110, 99, 56, 54, 104, 48, 18, 10, 116, 101, 115, 116, 45, 97, 108, 105, 97, 115, 26, 11, 116, 101, 115, 116, 45, 118, 101, 114, 107, 101, 121, 34, 8, 116, 101, 115, 116, 45, 100, 105, 100, 42, 9, 116, 101, 115, 116, 45, 114, 111, 108, 101, 18, 4, 109, 101, 109, 111, 24, 192, 2, 18, 97, 10, 78, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 59, 126, 95, 52, 102, 213, 99, 251, 102, 62, 148, 101, 72, 226, 188, 243, 222, 31, 35, 148, 19, 127, 79, 75, 79, 37, 160, 132, 193, 33, 148, 7, 18, 4, 10, 2, 8, 1, 18, 15, 10, 9, 10, 4, 99, 104, 101, 113, 18, 1, 48, 16, 224, 167, 18, 26, 64, 130, 229, 164, 76, 214, 244, 157, 39, 135, 11, 118, 223, 29, 196, 41, 92, 247, 126, 129, 194, 18, 154, 136, 165, 153, 76, 202, 85, 187, 195, 40, 69, 10, 206, 165, 238, 223, 245, 35, 140, 92, 123, 246, 110, 23, 39, 32, 215, 239, 230, 196, 146, 168, 5, 147, 9, 67, 113, 242, 163, 0, 223, 233, 73]
```

## Pool

### List of methods
- [VDR Tools SDK ledger connection API](#vdr-tools-sdk-ledger-connection-api)
	- [Overview](#overview)
	- [Identity wallet key methods](#identity-wallet-key-methods)
		- [indy_cheqd_keys_add_random](#indy_cheqd_keys_add_random)
			- [Input parameters](#input-parameters)
			- [Example output](#example-output)
		- [indy_cheqd_keys_add_from_mnemonic](#indy_cheqd_keys_add_from_mnemonic)
			- [**indy_cheqd_keys_get_info**](#indy_cheqd_keys_get_info)
			- [**indy_cheqd_keys_get_list_keys**](#indy_cheqd_keys_get_list_keys)
			- [**indy_cheqd_keys_sign**](#indy_cheqd_keys_sign)
	- [Pool](#pool)
		- [List of methods](#list-of-methods)
			- [**indy_cheqd_pool_add**](#indy_cheqd_pool_add)
			- [**indy_cheqd_pool_get_config**](#indy_cheqd_pool_get_config)
			- [**indy_cheqd_pool_get_all_config**](#indy_cheqd_pool_get_all_config)
			- [**indy_cheqd_pool_broadcast_tx_commit**](#indy_cheqd_pool_broadcast_tx_commit)
			- [**indy_cheqd_pool_abci_query**](#indy_cheqd_pool_abci_query)
			- [**indy_cheqd_pool_abci_info**](#indy_cheqd_pool_abci_info)
	- [Base connection workflow:](#base-connection-workflow)

#### **indy_cheqd_pool_add**
This method is needed for adding information about pool which will be used to connect to.

Input parameters:
* `alias` - is a human-readable string,
* `rpc_address` - address for connecting to the node, like `http://1.2.3.4:26657`, port `26657` is default value.
* `chain_id` - identifier of the network.
As result structure like PoolConfig is expected in response:
```
{
    "alias":"test_pool",
    "rpc_address":"rpc_address",
    "chain_id":"chain_id"
}
```
#### **indy_cheqd_pool_get_config**
This method is needed for getting config information about connecting to the pool by alias.

Input parameters:
* `alias` - human-readable string, represents pool alias, like `test_pool`.

Expected result is structure PoolConfig, like:
```
{
    "alias":"test_pool",
    "rpc_address":"rpc_address",
    "chain_id":"chain_id"
}
```
#### **indy_cheqd_pool_get_all_config**
The same as [indy_cheqd_pool_get_config](#indy_cheqd_pool_get_config) but returns the list of structures.
Response should be like:
```
[Object({
	"alias": String("test_pool_1"),
	"chain_id": String("chain_id"),
	"rpc_address": String("rpc_address")
}), 
Object({
	"alias": String("test_pool_2"),
	"chain_id": String("chain_id"),
	"rpc_address": String("rpc_address")
})]
```
#### **indy_cheqd_pool_broadcast_tx_commit**
This method allows to send a txn to all the nodes.

Input parameters:
* `pool_alias` - human-readable string, like `test_pool`,
* `signed_tx_raw` - string of bytes which includes raw signed transaction,
* `signed_tx_len` - length of signed txn string,

Request signed txn as input in raw format.
Expected response should be in json format like:
```
{
	"check_tx": {
		"code": 0,
		"data": "",
		"log": "[]",
		"info": "",
		"gas_wanted": "300000",
		"gas_used": "38591",
		"events": [],
		"codespace": ""
	},
	"deliver_tx": {
		"code": 0,
		"data": "Cg8KCUNyZWF0ZU55bRICCAU=",
		"log": [{
			"events ": [{
				"type ": "message",
				"attributes": [{
					"key": "action",
					"value": "CreateNym"
				}]
			}]
		}],
		"info": "",
		"gas_wanted": "300000",
		"gas_used": "46474",
		"events": [{
			"type": "message",
			"attributes": [{
				"key": "YWN0aW9u",
				"value": "Q3JlYXRlTnlt"
			}]
		}],
		"codespace": ""
	},
	"hash": "364441EDC5266A0B6AF5A67D4F05AC5D1FE95BFEDFBEBBE195723BEDBA877CAE",
	"height": "121"
}
```
#### **indy_cheqd_pool_abci_query**
Needs to send abci_query to the pool.
With pool `alias` it requires request in json format as input parameter.

Input parameters:
* `pool_alias` - human-readable string, like `test_pool`,
* `req_json`- String of ABCI query in json format,

Expected response should be in json format like:
```
{
   "nym":
   {
      "creator":"cosmos1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0",
      "id":4,
      "alias":"test-alias",
      "verkey":"test-verkey",
      "did":"test-did",
      "role":"test-role"
   }
}
```
#### **indy_cheqd_pool_abci_info**
Get general pool information. 
Requires only `pool alias` as input parameter.
Returns the response in json format, like:
```
"{
  "response":
     {
         "data": "cheqd-node",
         "version":"1.0",
         "app_version":"1",
         "last_block_height":"119",
         "last_block_app_hash":[120,105,48,70,72,98,101,55,55,97,112,84,54,98,65,116,76,71,88,76,43,65,90,107,114,75,73,78,104,88,83,119,102,118,115,105,111,54,105,67,53,106,99,61]
     }
 }
```

## Base connection workflow:
* Generate keys or restore them from mnemonic string. Useful methods here [`indy_cheqd_keys_add_random`](#indy_cheqd_keys_add_random) or [`indy_cheqd_keys_add_from_mnemonic`](#indy_cheqd_keys_add_from_mnemonic). 
* Add configuration about pool, by calling [`indy_cheqd_pool_add`](#indy_cheqd_pool_add)

The real connection will be created only when request or txn will be sent.
