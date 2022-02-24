# Ledger connections in VDR Tools SDK

## Overview

This page describes how [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) connects to the cheqd network ledger "pool".

It is worth noting here that the terminology of "pool" connection is specifically a legacy term originally used in [Hyperledger Indy](https://github.com/hyperledger/indy-node), which as a permissioned blockchain assumes there is a finite pools of servers. While this paradigm is no longer true in the public, permissionless world of the cheqd network, the identity APIs in VDR Tools SDK and similar Hyperledger Aries-based frameworks is retained for explanations.

## Ledger pool connection methods

Establishing a ledger "pool" connection in VDR Tools SDK broadly has the following steps:

1. Generate keys or restore them from mnemonic as described in [key management using VDR Tools SDK](vdr-tools-sdk-key-management.md).
2. Add ledger "pool" configuration as described using the methods below, i.e., `indy_cheqd_pool_add`. This only adds the configuration, without actually establishing the connection. The connection is established once the first transaction is sent.

### indy_cheqd_pool_add

Add a new cheqd network ledger `PoolConfig` configuration.

#### Input parameters

* `alias` (string): Friendly-name for pool connection
* `rpc_address` (string): Tendermint RPC endpoint (e.g., `http://localhost:26657`) for a cheqd network node(s) to send/receive transactions to.
* `chain_id` (string): cheqd network identifier, e.g., `cheqd-mainnet-1`

#### Example output

```jsonc
{
    "alias": "cheqd_pool",
    "rpc_address": "rpc_address",
    "chain_id": "chain_id"
}
```

#### **indy_cheqd_pool_get_config**
This method is needed for getting config information about connecting to the pool by alias.

#### Input parameters
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

#### Input parameters
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

#### Input parameters
* `pool_alias` - human-readable string, like `test_pool`,
* `req_json`- String of ABCI query in json format,

Expected response should be in json format like:
```
{
   "nym":
   {
      "creator":"cheqd1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0",
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

