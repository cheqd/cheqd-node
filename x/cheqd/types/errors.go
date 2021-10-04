package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cheqd module sentinel errors
var (
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 1100, "invalid signature detected")
)
