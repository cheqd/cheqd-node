package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cheqd module sentinel errors
var (
	ErrBadRequest           = sdkerrors.Register(ModuleName, 1000, "bad request")
	ErrInvalidSignature     = sdkerrors.Register(ModuleName, 1100, "invalid signature detected")
	ErrDidDocExists         = sdkerrors.Register(ModuleName, 1200, "DID Doc exists")
	ErrDidDocNotFound       = sdkerrors.Register(ModuleName, 1201, "DID Doc not found")
	ErrUnexpectedDidVersion = sdkerrors.Register(ModuleName, 1202, "unexpected DID version")
	ErrInvalidCredDefValue  = sdkerrors.Register(ModuleName, 1300, "invalid cred def value")
)
