package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyCreateDid     = []byte("KeyCreateDid")
	KeyUpdateDid     = []byte("UpdateDid")
	KeyDeactivateDid = []byte("DeactivateDid")
	KeyBurnFactor    = []byte("BurnFactor")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&FeeParams{})
}

func (tfp *FeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCreateDid, &tfp.CreateDid, validateCreateDid),
		paramtypes.NewParamSetPair(KeyUpdateDid, &tfp.UpdateDid, validateUpdateDid),
		paramtypes.NewParamSetPair(KeyDeactivateDid, &tfp.DeactivateDid, validateDeactivateDid),
		paramtypes.NewParamSetPair(KeyBurnFactor, &tfp.BurnFactor, validateBurnFactor),
	}
}
