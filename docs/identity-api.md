# Identity API

## Base write flow

* _Step 1._ Build a request. Example: `build_create_did_request(id, verkey, alias)`
* _Step_ 2_._ Sign request using DID key. Example:  `indy_crypto_sign(did, verkey)`
* _Step 3._ Built a transaction with the request from the previous step. Example: `build_tx(pool_alias, pub_key, builded_request, account_number, account_sequence, max_gas, max_coin_amount, denom, timeout_height, memo)`
* _Step 4._ Sign a transaction from the previous step. `cheqd_keys_sign(wallet_handle, key_alias, tx)`. 
* _Step 5._ Broadcast a signed transaction from the previous step. `broadcast_tx_commit(pool_alias, signed)`.

  Response format:

  ```text
    Response {
     check_tx: TxResult {
        code: 0,
        data: None,
        log: "",
        info: "",
        gas_wanted: 0,
        gas_used: 0,
        events: [
        ],
        codespace: ""
     },
     deliver_tx: TxResult {
        code: 0,
        data: Some(Data([...])),
        log: "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1pvnjjy3vz0ga6hexv32gdxydzxth7f86mekcpg\"},{\"key\":\"sender\",\"value\":\"cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"amount\",\"value\":\"100cheq\"}]}]}]",
        info: "",
        gas_wanted: 0,
        gas_used: 0,
        events: [...], 
        codespace: ""
     },
     hash: "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E",
     height: 353
  }
  ```

  `hash` - transaction hash

  `height` - ledger height

##   DID

### Create DID

_VDR tools: ****_build\_create\_did\_request\(id, verkey, alias\)

_Builds request in the follow format:_

```text
CreateDidRequest 
{
    "data": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "Alice did"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `id` \(base58-encoded string\) Target DID as base58-encoded string for 16 or 32 byte DID value.
* `verkey` \(base58-encoded string, possibly starting with "~"; optional\) Target verification key. It can start with "~", which means that it's abbreviated verkey and should be 16 bytes long when decoded, otherwise it's a full verkey which should be 32 bytes long when decoded.
* `alias` \(string; optional\).

_Returns:_

```text
CreateDidResponse {
    "key": "did:GEzcdDLhCpGCYRHW82kjHd" 
}  
```

* `key`\(string\): a key is used to store this DID in a state

#### Validation

* `CreateDidRequest` must be signed by  DID from `id` field. It means that this DID must be an owner of this DID transaction.

### Update DID

_VDR tools: ****_build\_update\_did\_request\(id, verkey, alias\)

_Builds request in the follow format:_

```text
UpdateDidRequest 
{
    "data": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "Alice did"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `id` \(base58-encoded string\) Target DID as base58-encoded string for 16 or 32 byte DID value.
* `verkey` \(base58-encoded string, possibly starting with "~"; optional\) Target verification key. It can start with "~", which means that it's abbreviated verkey and should be 16 bytes long when decoded, otherwise it's a full verkey which should be 32 bytes long when decoded.
* `alias` \(string; optional\).

_Returns:_

```text
UpdateDidResponse {
    "key": "did:GEzcdDLhCpGCYRHW82kjHd" 
}  
```

* `key`\(string\): a key is used to store this DID in a state

#### Validation

* A transaction with `id` from `UpdateDidRequest`must already be in a ledger created by `CreateDidRequest`
* `UpdateDidRequest` must be signed by  DID from `id` field. It means that this DID must be an owner of this DID transaction.

### Get DID

_VDR tools: ****_build\_query\_get\_did\(id\)

* `id` \(base58-encoded string\) Target DID as base58-encoded string for 16 or 32 byte DID value.

_Builds request in the follow format:_

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`_-_ path for RPC Endpoint for Cheqd pool; 
* `data` - query with an entity key from a state. String `did:<id>` encoded to bytes;
* `height` - a height of ledger \(size\). `None` for auto calculation;
* `prove` - boolean value. `True` - for getting state proof in a pool response. 

_Returns:_

```text
QueryGetDidResponse{
        "did": {
               "id": "GEzcdDLhCpGCYRHW82kjHd",
               "verkey": "~HmUWn928bnFT6Ephf65YXv",
               "alias": "Alice did"
             },
}  
```

## ATTRIB

### Create ATTRIB

_VDR tools: ****_build\_create\_attrib\_request\(did, raw\)

_Builds request in the follow format:_

```text
CreateAttribRequest 
{
    "data": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `did` \(base58-encoded string\):

  Target DID as base58-encoded string for 16 or 32 byte DID value.

* `raw` \(json; mutually exclusive with `hash` and `enc`\): Raw data is represented as json, where the key is attribute name and value is attribute value.

_Returns:_

```text
CreateAttribResponse {
    "key": "attrib:GEzcdDLhCpGCYRHW82kjHd" 
} 
```

* `key`\(string\): a key is used to store these attributes in a state

#### Validation

* A DID transaction with `id` from `UpdateAttribRequest`must already be in a ledger created by `CreateDidRequest`
* `CreateAttribRequest` must be signed by  DID from `did` field. It means that this DID must be an owner of this ATTRIB transaction.

### Update ATTRIB

_VDR tools: ****_build\_update\_attrib\_request\(id, raw\)

_Builds request in the follow format:_

```text
UpdateAttribRequest 
{
    "data": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `did` \(base58-encoded string\):

  Target DID as base58-encoded string for 16 or 32 byte DID value.

* `raw` \(json; mutually exclusive with `hash` and `enc`\): Raw data is represented as json, where the key is attribute name and value is attribute value.

_Returns:_

```text
UpdateAttribResponse {
        "key": "attrib:GEzcdDLhCpGCYRHW82kjHd" 
} 
```

* `key`\(string\): a key is used to store these attributes in a state

#### Validation

* A DID transaction with `id` from `UpdateAttribRequest`must already be in a ledger created by `CreateDidRequest`
* `UpdateAttribRequest` must be signed by  DID from `did` field. It means that this DID must be an owner of this ATTRIB transaction.

### Get ATTRIB

_VDR tools: ****_build\_query\_get\_attrib\(did\)

* `did` \(base58-encoded string\) Target DID as base58-encoded string for 16 or 32 byte DID value.

_Builds request in the follow format:_

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`_-_ path for RPC Endpoint for Cheqd pool; 
* `data` - query with an entity key from a state. String `attrib:<did>` encoded to bytes;
* `height` - a height of ledger \(size\). `None` for auto calculation;
* `prove` - boolean value. `True` - for getting state proof in a pool response. 

_Returns:_

```text
QueryGetAttribResponse{
        "attrib": {
               "did": "GEzcdDLhCpGCYRHW82kjHd",
               "raw": "{'name': 'Alice'}"
             },
}  
```

## SCHEMA

### Create Schema

_VDR tools: ****_build\_create\_schema\_request\(version, name, attr\_names\)

_Builds request in the follow format:_

```text
CreateSchemaRequest 
{
    "data": {
            "version": "1.0",
            "name": "Degree",
            "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `attr_names`\(array\): array of attribute name strings \(125 attributes maximum\)
* `name`\(string\): Schema's name string
* `version`\(string\): Schema's version string

_Returns:_

```text
CreateSchemaResponse {
        "key": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0" 
} 
```

* `key`\(string\): a key is used to store this schema in a state

#### Validation

* A SCHEMA transaction with did from `owner` field must already be in a ledger created by `CreateDidRequest`
* `CreateSchemaRequest` must be signed by  DID from `owner` field. 

### Get Schema

_VDR tools: ****_build\_query\_get\_schema\(name, version, owner\)

* `name`\(string\): Schema's name string
* `version`\(string\): Schema's version string
* `owner` \(string\): Schema's owner did

_Builds request in the follow format:_

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`_-_ path for RPC Endpoint for Cheqd pool; 
* `data` - query with an entity key from a state. String `schema:<owner>:<name>:<version>` encoded to bytes;
* `height` - a height of ledger \(size\). `None` for auto calculation;
* `prove` - boolean value. `True` - for getting state proof in a pool response. 

_Returns:_

```text
QueryGetSchemaResponse{
        "attrib": {
                "version": "1.0",
                "name": "Degree",
                "attr_names": ["undergrad", "last_name", "first_name", "birth_date", "postgrad", "expiry_date"]
             },
}  
```

## CRED\_DEF

### Create Credential Definition

_VDR tools: ****_build\_create\_cred\_def\_request\(cred\_def, schema\_id, signature\_type, tag\)

_Builds request in the follow format:_

```text
CreateCredDefRequest 
{
    "data": {
                "signature_type": "CL",
                "schema_id": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0",
                "tag": "some_tag",    
                "cred_def": {
                    "primary": ....,
                    "revocation": ....
            },
    "owner": "GEzcdDLhCpGCYRHW82kjHd",
    "signature": "49W5WP5jr7x1fZhtpAhHFbuUDqUYZ3AKht88gUjrz8TEJZr5MZUPjskpfBFdboLPZXKjbGjutoVascfKiMD5W7Ba",
    "metadata": {}
}
```

* `cred_def` \(dict\): Dictionary with Cred Definition's data:
  * `primary` \(dict\): primary credential public key
  * `revocation` \(dict\): revocation credential public key
* `schema_id` \(string\): Schema_'_s key from a state
* `signature_type` \(string\): Type of the credential definition \(that is credential signature\). `CL` \(Camenisch-Lysyanskaya\) is the only supported type now.
* `tag` \(string, optional\): A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

_Returns:_

```text
CreateCredDefResponse {
        "key": "cred_def:GEzcdDLhCpGCYRHW82kjHd:schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0:some_tag:CL" 
} 
```

* `key`\(string\): a key is used to store this Credential Definition in a state

#### Validation

* A CRED\_DEF transaction with did from `owner` field must already be in a ledger created by `CreateDidRequest`
* `CreateCredDefRequest` must be signed by  DID from `owner` field. 

### Get Credential Definition

_VDR tools: ****_build\_query\_get\_cred\_def\(name, version, owner\)

* `schema_id`\(string\): Schema's key from a state
* `signature_type`\(string\): Type of the credential definition \(that is credential signature\). CL \(Camenisch-Lysyanskaya\) is the only supported type now.
* `owner` \(string\): Credential Definition's owner did
* `tag` \(string, optional\): A unique tag to have multiple public keys for the same Schema and type issued by the same DID. A default tag `tag` will be used if not specified.

_Builds request in the follow format:_

```text
Request 
{
    "path": "/store/cheqd/key",
    "data": <bytes>,
    "height": 642,
    "prove": true
}
```

* `path`_-_ path for RPC Endpoint for Cheqd pool; 
* `data` - query with an entity key from a state. String `cred_def:<owner>:<schema_id>:<tag>:<signature_type>` encoded to bytes;
* `height` - a height of ledger \(size\). `None` for auto calculation;
* `prove` - boolean value. `True` - for getting state proof in a pool response. 

_Returns:_

```text
QueryGetCredDefResponse{
        "cred_def": {
                "signature_type": "CL",
                "schema_id": "schema:GEzcdDLhCpGCYRHW82kjHd:Degree:1.0",
                "tag": "some_tag",    
                "cred_def": {
                    "primary": ....,
                    "revocation": ....
         },
}  
```

## 

