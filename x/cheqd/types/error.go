package types

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
	ErrBasicValidation            = sdkerrors.Register(ModuleName, 1205, "basic validation failed")
	ErrNamespaceValidation = sdkerrors.Register(ModuleName, 1206, "DID namespace validation failed")
	ErrUnpackStateValue    = sdkerrors.Register(ModuleName, 1300, "invalid did state value")
	ErrSetToState          = sdkerrors.Register(ModuleName, 1304, "cannot set to state")
	ErrInternal                   = sdkerrors.Register(ModuleName, 1500, "internal error")
	ErrNotImplemented             = sdkerrors.Register(ModuleName, 1501, "not implemented")
	ErrValidatorInitialisation    = sdkerrors.Register(ModuleName, 1502, "can't init validator")
	// Static validation errors
	ErrStaticDIDBadMethod              = sdkerrors.Register(ModuleName, 1600, "DID method is not cheqd")
	ErrStaticDIDNamespaceNotAllowed    = sdkerrors.Register(ModuleName, 1601, "Namespace is not allowed for this network")
	ErrStaticDIDNamespaceNotValid      = sdkerrors.Register(ModuleName, 1602, "Namespace is not valid")
	ErrStaticDIDBadUniqueIDLen         = sdkerrors.Register(ModuleName, 1603, "Length of unique ID should be 16 or 32 symbols")
	ErrStaticDIDNotBase58ID            = sdkerrors.Register(ModuleName, 1604, "Not base58 symbols for unique ID string")
	ErrStaticDIDURLPathAbemptyNotValid = sdkerrors.Register(ModuleName, 1605, "There are not allowed symbols in Path.")
	ErrStaticDIDURLQueryNotValid       = sdkerrors.Register(ModuleName, 1606, "Query part in DIDUrl is not valid.")
	ErrStaticDIDURLFragmentNotValid    = sdkerrors.Register(ModuleName, 1607, "Fragment part in DIDUrl is not valid.")
)
