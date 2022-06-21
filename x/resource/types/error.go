package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/resource module sentinel errors
// TODO: Rework these errors to be used in the module
var (
	ErrBadRequest                 = sdkerrors.Register(ModuleName, 2000, "bad request")
	ErrInvalidSignature           = sdkerrors.Register(ModuleName, 2100, "invalid signature detected")
	ErrSignatureNotFound          = sdkerrors.Register(ModuleName, 2101, "signature is required but not found")
	ErrResourceExists             = sdkerrors.Register(ModuleName, 2200, "Resoure exists")
	ErrVerificationMethodNotFound = sdkerrors.Register(ModuleName, 2202, "verification method not found")
	ErrUnexpectedDidVersion       = sdkerrors.Register(ModuleName, 2203, "unexpected DID version")
	ErrBasicValidation            = sdkerrors.Register(ModuleName, 2205, "basic validation failed")
	ErrNamespaceValidation        = sdkerrors.Register(ModuleName, 2206, "DID namespace validation failed")
	ErrUnpackStateValue           = sdkerrors.Register(ModuleName, 2300, "invalid did state value")
	ErrInternal                   = sdkerrors.Register(ModuleName, 2500, "internal error")
)
