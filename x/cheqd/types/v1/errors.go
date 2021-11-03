package v1

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cheqd module sentinel errors
var (
	ErrBadRequest                 = sdkerrors.Register(ModuleName, 1000, "bad request")
	ErrBadRequestIsRequired       = sdkerrors.Register(ModuleName, 1001, "is required")
	ErrBadRequestIsNotDid         = sdkerrors.Register(ModuleName, 1002, "is not DID")
	ErrBadRequestInvalidVerMethod = sdkerrors.Register(ModuleName, 1003, "invalid verification method")
	ErrBadRequestInvalidService   = sdkerrors.Register(ModuleName, 1004, "invalid service")
	ErrBadRequestIsNotDidFragment = sdkerrors.Register(ModuleName, 1005, "is not DID fragment")
	ErrInvalidSignature           = sdkerrors.Register(ModuleName, 1100, "invalid signature detected")
	ErrDidDocExists               = sdkerrors.Register(ModuleName, 1200, "DID Doc exists")
	ErrDidDocNotFound             = sdkerrors.Register(ModuleName, 1201, "DID Doc not found")
	ErrVerificationMethodNotFound = sdkerrors.Register(ModuleName, 1202, "verification method not found")
	ErrUnexpectedDidVersion       = sdkerrors.Register(ModuleName, 1203, "unexpected DID version")
	ErrInvalidPublicKey           = sdkerrors.Register(ModuleName, 1204, "invalid public key")
	ErrInvalidDidStateValue       = sdkerrors.Register(ModuleName, 1300, "invalid did state value")
	ErrSetToState                 = sdkerrors.Register(ModuleName, 1304, "cannot set to state")
	ErrNotImplemented             = sdkerrors.Register(ModuleName, 1501, "not implemented")
)
