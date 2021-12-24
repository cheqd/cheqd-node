# Registered cheqd errors

## Overview

| Name | Code   | Description  |
|---|---|---|
| ErrBadRequest  |  1000 | The request the client made is incorrect or corrupt |
| ErrBadRequestIsRequired  |  1001 | The request does not contain required property |
| ErrBadRequestIsNotDid  | 1002  | The request contains invalid `id` property |
| ErrBadRequestInvalidVerMethod  | 1003 | The request contains invalid verification method  |
| ErrBadRequestInvalidService  | 1004  | The request contains invalid service |
| ErrBadRequestIsNotDidFragment  |  1005 | The request contains invalid verification method id |
| ErrInvalidSignature  | 1100  | Invalid signature detected |
| ErrDidDocExists  | 1200  | An attempt to create a DID Doc that exists in the ledger detected |
| ErrDidDocNotFound  | 1201  | The DID Doc not found in the ledger |
| ErrVerificationMethodNotFound  | 1202  | The DID Doc does not contain the requested verification method  |
| ErrUnexpectedDidVersion  | 1203  | Replay protected failed. An attempt to update DID Doc with wrong version detected |
| ErrInvalidPublicKey  | 1204  | Unable to decode public key |
| ErrInvalidDidStateValue  | 1300  | Unable to unmarshall stored document |
| ErrSetToState  |  1304 | Unable to set value into the ledger |
| ErrNotImplemented  |  1501 | The method is not implemented |
