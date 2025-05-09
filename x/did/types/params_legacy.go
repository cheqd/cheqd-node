package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var ParamStoreKey = []byte("feeparams")

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&FeeParams{})
}

func (tfp *FeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKey, &FeeParams{}, validateFeeParams),
	}
}
