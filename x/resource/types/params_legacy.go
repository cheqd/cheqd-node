package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var ParamStoreKeyFeeParams = []byte("feeparams")

// ParamKeyTable returns the key declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&LegacyFeeParams{})
}

// NewParams creates a new FeeParams object with specified parameters
func (tfp *LegacyFeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyFeeParams, &LegacyFeeParams{}, validateFeeParams),
	}
}
