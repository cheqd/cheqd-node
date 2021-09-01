# ADR 001: Payment mechanism for issuing credentials

## Status

PROPOSED

## Summary

The Aries protocol describes the payment mechanism for issuing credentials. It is necessary to establish which public API cheqd-node can provide for implementation of payments for issuing through cheqd coins transfer for Connect.me.

## Decision

According to Aries protocol we can use payment decorators `~payment_request` and `~payment_receipt` credential issuing.

### Step 1: Credential Offer

A message sent by the Issuer to the potential Holder, describing the credential they intend to offer and possibly the price they expect to be paid. [Aries Credential Offer](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md#offer-credential)

```text
    "@type": "https://didcomm.org/issue-credential/1.0/offer-credential",
    "@id": "<uuid-of-offer-message>",
    "comment": "some comment",
    "credential_preview": <json-ld object>,
    "offers~attach": [
                        {
                            "@id": "libindy-cred-offer-0",
                            "mime-type": "application/json",
                            "data": {
                                        "base64": "<bytes for base64>"
                                    }
                        }
                    ]
    "~payment_request": { ... }
}
```

And use a payment decorator to add information about an issuing price and address for sending payment transaction. [Aries payment\_request decorator](https://github.com/hyperledger/aries-rfcs/blob/main/features/0075-payment-decorators/README.md#payment_request)

```text
   "~payment_request": {
        "methodData": [
          {
            "supportedMethods": "cheqd",
            "data": {
              "payeeId": "cosmos1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd"
            },
          }
        ],
        "details": {
          "id": "0a2bc4a6-1f45-4ff0-a046-703c71ab845d",
          "displayItems": [
            {
              "label": "commercial driver's license",
              "amount": { "currency": "cheq", "value": "10" },
            }
          ],
          "total": {
            "label": "Total due",
            "amount": { "currency": "cheq", "value": "10" }
          }
        }
      }
```

* `details.id` field contains an invoice number that unambiguously identifies a credential for which payment is requested. When paying, this value should be placed in `memo` field for Cheqd payment transaction.
* `payeeId` field contains a Cheqd account address in the cosmos format.

### Step 2: Payment transaction

This operation has 5 steps:

* _Step 2.1._ Build a request for transferring coins. Example: `cheqd_ledger::bank::build_msg_send(from_account, to_account, amount_for_transfer, denom)`.
  * `from_account` the potential Holder Cheqd account address
  * `to_account` the same with `payeeId` from a payment request
  * `amount_for_transfer` the same with `details.total.amount.value` from a payment request
  * `denom` the same with `details.total.amount.currency` from a payment request
* _Step 2.2._ Built a transaction with the request from the previous step. Example: `cheqd_ledger::auth::build_tx(pool_alias, pub_key, builded_request, account_number, account_sequence, max_gas, max_coin_amount, denom, timeout_height, memo)`. 
  * `memo` tha same with `details.id` from a payment request
* _Step 2.3._ Sign a transaction from the previous step. `cheqd_keys::sign(wallet_handle, key_alias, tx)`. 
* _Step 2.4._ Broadcast a signed transaction from the previous step. `cheqd_pool::broadcast_tx_commit(pool_alias, signed)`.

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

[Read more about Cheqd payment transaction](https://gitlab.com/evernym/verity/vdr-tools/-/tree/cheqd/docs/design/014-bank-transactions)

### Step 3: Credential Request

This is a message sent by the potential Holder to the Issuer, to request the issuance of a credential. After sending a payment transaction.

```text
{
    "@type": "https://didcomm.org/issue_credential/1.0/request_credential",
    "@id": "94af9be9-5248-4a65-ad14-3e7a6c3489b6",
    "~thread": { "thid": "5bc1989d-f5c1-4eb1-89dd-21fd47093d96" },
    "cred_def_id": "KTwaKJkvyjKKf55uc6U8ZB:3:CL:59:tag1",
    "~payment_receipt": {
      "request_id": "0a2bc4a6-1f45-4ff0-a046-703c71ab845d",
      "selected_method": "cheqd",
      "transaction_id": "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E",
      "payeeId": "0xD15239C7e7dDd46575DaD9134a1bae81068AB2A4",
      "amount": { "currency": "cheq", "value": "10.0" }
    }
}
```

`request_id` the same with `details.id` from payment\_request and with `memo` from a payment transaction

### Step 4: Check payment\_receipt

Issuer receives Credential Request + `payment_receipt` with payment `transaction_id`. It allows Issuer

* get the payment transaction by hash from Cheqd Ledger using `get_tx_by_hash(hash)` method. `hash` parameter is `transaction_id` from previous steps.
* check that `memo` field from received transaction contains `request_id`.

### Step 5: Credential issuing

Credential issuing according Aries protocol.

## References

* [Cheqd payment transaction](https://gitlab.com/evernym/verity/vdr-tools/-/tree/cheqd/docs/design/014-bank-transactions)
* [Aries Credential Offer](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md#offer-credential)
* [Aries payment\_request decorator](https://github.com/hyperledger/aries-rfcs/blob/main/features/0075-payment-decorators/README.md#payment_request)

