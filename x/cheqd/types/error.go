package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cheqd module sentinel errors
var (
	ErrBadRequest                 = sdkerrors.Register(ModuleName, 1000, "bad request")
	ErrInvalidSignature           = sdkerrors.Register(ModuleName, 1100, "invalid signature detected")
	ErrSignatureNotFound          = sdkerrors.Register(ModuleName, 1101, "signature is required but not found")
	ErrDidDocExists               = sdkerrors.Register(ModuleName, 1200, "DID Doc exists")
	ErrDidDocNotFound             = sdkerrors.Register(ModuleName, 1201, "DID Doc not found")
	ErrVerificationMethodNotFound = sdkerrors.Register(ModuleName, 1202, "verification method not found")
	ErrUnexpectedDidVersion       = sdkerrors.Register(ModuleName, 1203, "unexpected DID version")
	ErrBasicValidation            = sdkerrors.Register(ModuleName, 1205, "basic validation failed")
	ErrNamespaceValidation        = sdkerrors.Register(ModuleName, 1206, "DID namespace validation failed")
	ErrDIDDocDeactivated          = sdkerrors.Register(ModuleName, 1207, "DID Doc already deactivated")
	ErrUnpackStateValue           = sdkerrors.Register(ModuleName, 1300, "invalid did state value")
	ErrInternal                   = sdkerrors.Register(ModuleName, 1500, "internal error")
)
