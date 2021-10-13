package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cheqd module sentinel errors
var (
	ErrInvalidSignature     = sdkerrors.Register(ModuleName, 1100, "invalid signature detected")
	ErrDidDocExists         = sdkerrors.Register(ModuleName, 1200, "did doc exists")
	ErrUnexpectedDidVersion = sdkerrors.Register(ModuleName, 1201, "unexpected did version")
	ErrInvalidCredDefValue  = sdkerrors.Register(ModuleName, 1300, "invalid cred def value")
)
