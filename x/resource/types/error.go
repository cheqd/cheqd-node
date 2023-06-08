package types

import sdkerrors "cosmossdk.io/errors"

// DONTCOVER

// x/resource module sentinel errors
var (
	ErrBadRequest      = sdkerrors.Register(ModuleName, 2000, "bad request")
	ErrResourceExists  = sdkerrors.Register(ModuleName, 2200, "Resource exists")
	ErrBasicValidation = sdkerrors.Register(ModuleName, 2205, "basic validation failed")
	ErrInternal        = sdkerrors.Register(ModuleName, 2500, "internal error")
)
