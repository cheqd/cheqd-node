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

func (p *FeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCreateDid, &p.CreateDid, validateCreateDid),
		paramtypes.NewParamSetPair(KeyUpdateDid, &p.UpdateDid, validateUpdateDid),
		paramtypes.NewParamSetPair(KeyDeactivateDid, &p.DeactivateDid, validateDeactivateDid),
		paramtypes.NewParamSetPair(KeyBurnFactor, &p.BurnFactor, validateBurnFactor),
	}
}
