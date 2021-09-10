# Keys and pool API

## Overview

This page describes API about how to work with keys and how to connect to the pool. 
All the libraries methods/calls are defined on the [`cheqd-sdk`](https://github.com/cheqd/cheqd-sdk) \(forked from [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools)\).

## Keys
### List of methods
* [`indy_cheqd_keys_add_random`](#indy_cheqd_keys_add_random)
* [`indy_cheqd_keys_add_from_mnemonic`](#indy_cheqd_keys_add_from_mnemonic)
* [`indy_cheqd_keys_get_info`](#indy_cheqd_keys_get_info)
* [`indy_cheqd_keys_get_list_keys`](#indy_cheqd_keys_get_list_keys)
* [`indy_cheqd_keys_sign`](#indy_cheqd_keys_sign)

#### **indy_cheqd_keys_add_random**
This method implements the logic of creation a key just using `alias`, without specifying any other additional information, like mnemonic.
As result, the next structure is expected:
```
alias: String,
// Cosmos address
account_id: String,
// Base58-encoded SEC1-encoded secp256k1 ECDSA key
pub_key: String,
```

#### **indy_cheqd_keys_add_from_mnemonic**
This method realised logic for recovering keys from `mnemonic` string. Mnemonic string - it's a human-readable combination of words, in general can include 12 or 24 words.
As result, the next structure is expected:
```
alias: String,
// Cosmos address
account_id: String,
// Base58-encoded SEC1-encoded secp256k1 ECDSA key
pub_key: String,
```

#### **indy_cheqd_keys_get_info**
This method is needed for getting information about already generated and stored keys by using only an `alias`.
As result, the list of next structures is expected:
```
alias: String,
// Cosmos address
account_id: String,
// Base58-encoded SEC1-encoded secp256k1 ECDSA key
pub_key: String,
```

#### **indy_cheqd_keys_get_list_keys**
The method returns a list of all keys which are placed locally.
As result, the next structure is expected:
```
alias: String,
// Cosmos address
account_id: String,
// Base58-encoded SEC1-encoded secp256k1 ECDSA key
pub_key: String,
```

#### **indy_cheqd_keys_sign**
This method can sign a transaction by using a key which can be found by `alias`
As result, raw byte's string is expected.

## Pool

### list of methods
* [`indy_cheqd_pool_add`](#indy_cheqd_pool_add)
* [`indy_cheqd_pool_get_config`](#indy_cheqd_pool_get_config)
* [`indy_cheqd_pool_get_all_config`](#indy_cheqd_pool_get_all_config)
* [`indy_cheqd_pool_broadcast_tx_commit`](#indy_cheqd_pool_broadcast_tx_commit)
* [`indy_cheqd_pool_abci_query`](#indy_cheqd_pool_abci_query)
* [`indy_cheqd_pool_abci_info`](#indy_cheqd_pool_abci_info)

#### **indy_cheqd_pool_add**
This method is needed for adding information about pool which will be used to connect to.
Input parameters:
* `alias` - is a human-readable string,
* `rpc_address` - address for connecting to the node, like `http://1.2.3.4:26657`, port `26657` is default value.
* `chain_id` - identifier of the network.
As result structure like PoolConfig is expected:
```
alias: String,
rpc_address: String,
chain_id: String,
}
```
#### **indy_cheqd_pool_get_config**
This method is needed for getting config information about connecting to the pool by alias.
Expected result is structure PoolConfig, like:
```
alias: String,
rpc_address: String,
chain_id: String,
}
```
#### **indy_cheqd_pool_get_all_config**
The same as [indy_cheqd_pool_get_config](#indy_cheqd_pool_get_config) but returns the list of structures.
#### **indy_cheqd_pool_broadcast_tx_commit**
This method allows to send a txn to all the nodes.
Request signed txn as input in raw format.
#### **indy_cheqd_pool_abci_query**
Needs to send abci_query to the pool.
With pool `alias` it requires request in json format as input parameter.
#### **indy_cheqd_pool_abci_info**
Get general pool information. 
Reuires only `pool alias`
Returns the response in json format, like:
```
"response": {
      "data": "cheqdnode",
      "version": "1.0",
      "app_version": "1",
      "last_block_height": "791638",
      "last_block_app_hash": "dNurd7DVM06TCDJtZ4rd6RfWraKhxjptbARRbKPAF30="
    }
```

## Base connection workflow:
* Generate keys or restore them from mnemonic string. Useful methods here [`indy_cheqd_keys_add_random`](#indy_cheqd_keys_add_random) or [`indy_cheqd_keys_add_from_mnemonic`](#indy_cheqd_keys_add_from_mnemonic). 
* Add configuration about pool, by calling [`indy_cheqd_pool_add`](#indy_cheqd_pool_add)
THe real connection will be created only when request or txn will be sent.
