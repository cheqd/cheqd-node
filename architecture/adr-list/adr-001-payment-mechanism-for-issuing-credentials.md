# ADR 001: Payment mechanism for issuing credentials

## Status

| Category | Status |
| :--- | :--- |
| **Authors** | Ankur Banerjee |
| **ADR Stage** | ACCEPTED |
| **Implementation Status** | Not Implemented |
| **Start Date** | 2021-09-01 |

## Summary

The [Hyperledger Aries](https://github.com/hyperledger/aries) protocol describes a payment mechanism that can used to pay for the issuance of credentials.

It is necessary to establish which public APIs from Hyperledger Aries can be implemented in `cheqd-node` to provide an implementation of payments using CHEQ tokens using a well-understood SSI protocol.

## Decision

Hyperledger Aries protocol has the concept of payment "decorators" `~payment_request` and `~payment_receipt` in requests, that can be used to pay using tokens for credential issuance.

### Step 1: Credential Offer

A message is sent by the Issuer to the potential Holder, describing the credential they intend to offer and optionally, the price the issuer would be expected to be paid for said credential. This is based on the [Hyperledger Aries credential offer RFC](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md#offer-credential).

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

A payment request can then be defined using the [Hyperledger Aries Payment Decorator](https://github.com/hyperledger/aries-rfcs/blob/main/features/0075-payment-decorators/README.md#payment_request) to add information about an issuing price and address where payment should be sent.

```json
   "~payment_request": {
        "methodData": [
          {
            "supportedMethods": "cheqd",
            "data": {
              "payeeId": "cheqd1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd"
            },
          }
        ],
        "details": {
          "id": "0a2bc4a6-1f45-4ff0-a046-703c71ab845d",
          "displayItems": [
            {
              "label": "commercial driver's license",
              "amount": { "currency": "ncheq", "value": "1000" },
            }
          ],
          "total": {
            "label": "Total due",
            "amount": { "currency": "ncheq", "value": "1000" }
          }
        }
      }
```

* **`details.id`** field contains an invoice number that unambiguously identifies a credential for which payment is requested. When paying, this value should be placed in `memo` field for the cheqd payment transaction.
* **`payeeId`** field contains a cheqd account address in the correct format for cheqd network.

### Step 2: Payment transaction flow

The payment flow can be broken down into five steps:

1. **Build a request for transferring tokens**. Example: `cheqd_ledger::bank::build_msg_send(from_account, to_account, amount_for_transfer, denom)`
   * **`from_account`**: The prospective credential holder's cheqd account address
   * **`to_account`**: Same as `payeeId` from the Payment Request
   * **`amount_for_transfer`**: Price of credential issuance defined as `details.total.amount.value` from the Payment Request
   * **`denom`**: Defined in `details.total.amount.currency` from the Payment Request
2. **Build a transaction with the request from the previous step** Example: `cheqd_ledger::auth::build_tx(pool_alias, pub_key, builded_request, account_number, account_sequence, max_gas, max_coin_amount, denom, timeout_height, memo)`
   * `memo`: This should be the same as `details.id` from the Payment Request
3. **Sign the transaction** Example:`cheqd_keys::sign(wallet_handle, key_alias, tx)`.
4. **Broadcast the signed transaction** Example: `cheqd_pool::broadcast_tx_commit(pool_alias, signed)`.

#### Response format

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
      log: "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cheqd1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cheqd1pvnjjy3vz0ga6hexv32gdxydzxth7f86mekcpg\"},{\"key\":\"sender\",\"value\":\"cheqd1fknpjldck6n3v2wu86arpz8xjnfc60f99ylcjd\"},{\"key\":\"amount\",\"value\":\"1000ncheq\"}]}]}]",
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

Key fields in the response above are:

* `hash`: Transaction hash
* `height`: Ledger height

### Step 3: Credential Request

This is a message sent by the potential Holder to the Issuer, to request the issuance of a credential after tokens are transferred to the nominated account using a Payment Transaction.

```json
{
    "@type": "https://didcomm.org/issue_credential/1.0/request_credential",
    "@id": "94af9be9-5248-4a65-ad14-3e7a6c3489b6",
    "~thread": { "this": "5bc1989d-f5c1-4eb1-89dd-21fd47093d96" },
    "cred_def_id": "KTwaKJkvyjKKf55uc6U8ZB:3:CL:59:tag1",
    "~payment_receipt": {
      "request_id": "0a2bc4a6-1f45-4ff0-a046-703c71ab845d",
      "selected_method": "cheqd",
      "transaction_id": "1B3B00849B4D50E8FCCF50193E35FD6CA5FD4686ED6AD8F847AC8C5E466CFD3E",
      "payeeId": "0xD15239C7e7dDd46575DaD9134a1bae81068AB2A4",
      "amount": { "currency": "ncheq", "value": "1000.0" }
    }
}
```

**`request_id`** should be the same as `details.id` from Payment Request and `memo` from Payment Transaction.

### Step 4: Check payment\_receipt

Issuer receives Credential Request + `payment_receipt` with payment `transaction_id`. It allows the Issuer to:

* Get the payment transaction by hash from cheqd network ledger using `get_tx_by_hash(hash)` method, where `hash` is `transaction_id` from previous steps.
* Check that `memo` field from received transaction contains the correct `request_id`.

### Step 5: Credential issuing

If steps 1-4 are successful, the Issuer is able to confirm that the requested payment has been made using CHEQ tokens. The credential issuing process can then proceed using standard Hyperledger Aries protocol procedures.

### Overview of steps 1-5

REPLACE WITH PNG

#### UML version

Editable version available on [swimlanes.io](https://swimlanes.io/u/6_9Qx9GOe?rev=2) or as text for compatible UML diagram generators below:

```text
Issuer -> Holder: Credential Offer (+ payment_request)
Holder -> Ledger: payment transaction (with payment_request id in memo)
Ledger -> Holder: payment transaction response (with transaction_hash)
Holder -> Issuer: Credential Request (+ payment_receipt)
Issuer -> Ledger: Get payment transaction by hash
Ledger -> Issuer: Payment transaction
Issuer -> Issuer: Check `memo` field from received transaction
Issuer -> Holder: Credential
Holder -> Issuer: Accept
```

## Consequences

### Backward Compatibility

* Credential issuance outside of the payment flow is compatible with and carried out using existing Hyperledger Aries protocol procedures. This should provide a level of compatibility with existing apps/SDKs that implement Aries protocol.
* Defining the transaction in CHEQ tokens is specific to the cheqd network.

### Positive

* By defining the payment mechanism using Hyperledger Aries protocols, this allows the possibility in the future to support payments on multiple networks.
* Existing SSI app developers should already be familiar with Hyperledger Aries (if building on Hyperledger Indy) and provides a transition path to add new functionality.

### Negative

* Hyperledger Aries may not be a familiar protocol for other Cosmos projects.
* Using the Payment Decorator in practice means there could be interoperability challenges at in implementations that impact credential issuance and exchange.

### Neutral

* N/A

## References

* [Hyperledger Aries RFC 0036: Issue Credential Protocol 1.0](https://github.com/hyperledger/aries-rfcs/blob/main/features/0036-issue-credential/README.md)
* [Hyperledger Aries RFC 0075: Payment Decorators](https://github.com/hyperledger/aries-rfcs/blob/main/features/0075-payment-decorators/README.md)
* [Evernym VDR Tools cheqd network payments ADR](https://gitlab.com/evernym/verity/vdr-tools/-/tree/main/docs/design/014-bank-transactions)

