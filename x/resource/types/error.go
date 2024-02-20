package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/resource module sentinel errors
var (
	ErrBadRequest       = errorsmod.Register(ModuleName, 2000, "bad request")
	ErrResourceExists   = errorsmod.Register(ModuleName, 2200, "Resource exists")
	ErrBasicValidation  = errorsmod.Register(ModuleName, 2205, "basic validation failed")
	ErrInternal         = errorsmod.Register(ModuleName, 2500, "internal error")
	ErrInvalidVersion   = errorsmod.Register(ModuleName, 2505, "invalid ibc version")
	ErrUnexpectedAck    = errorsmod.Register(ModuleName, 2510, "resource module never sends packets")
	ErrUnexpectedPacket = errorsmod.Register(ModuleName, 2515, "IBC packet is incorrect")
	ErrResourceNotAvail = errorsmod.Register(ModuleName, 2525, "resource not available")
)
