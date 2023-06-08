package types

import errorsmod "cosmossdk.io/errors"

// DONTCOVER

// x/resource module sentinel errors
var (
	ErrBadRequest      = errorsmod.Register(ModuleName, 2000, "bad request")
	ErrResourceExists  = errorsmod.Register(ModuleName, 2200, "Resource exists")
	ErrBasicValidation = errorsmod.Register(ModuleName, 2205, "basic validation failed")
	ErrInternal        = errorsmod.Register(ModuleName, 2500, "internal error")
)
