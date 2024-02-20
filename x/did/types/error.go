package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cheqd module sentinel errors
var (
	ErrBadRequest                 = errorsmod.Register(ModuleName, 1000, "bad request")
	ErrInvalidSignature           = errorsmod.Register(ModuleName, 1100, "invalid signature detected")
	ErrSignatureNotFound          = errorsmod.Register(ModuleName, 1101, "signature is required but not found")
	ErrDidDocExists               = errorsmod.Register(ModuleName, 1200, "DID Doc exists")
	ErrDidDocNotFound             = errorsmod.Register(ModuleName, 1201, "DID Doc not found")
	ErrVerificationMethodNotFound = errorsmod.Register(ModuleName, 1202, "verification method not found")
	ErrUnexpectedDidVersion       = errorsmod.Register(ModuleName, 1203, "unexpected DID version")
	ErrBasicValidation            = errorsmod.Register(ModuleName, 1205, "basic validation failed")
	ErrNamespaceValidation        = errorsmod.Register(ModuleName, 1206, "DID namespace validation failed")
	ErrDIDDocDeactivated          = errorsmod.Register(ModuleName, 1207, "DID Doc already deactivated")
	ErrUnpackStateValue           = errorsmod.Register(ModuleName, 1300, "invalid did state value")
	ErrInternal                   = errorsmod.Register(ModuleName, 1500, "internal error")
)
