package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyResourceImageFee      = []byte("KeyResourceImageFee")
	KeyCreateResourceJSONFee = []byte("KeyCreateResourceJSONFee")
	KeyCreateResourceKeyFee  = []byte("KeyCreateResourceKeyFee")
	KeyBurnFactor            = []byte("KeyBurnFactor")
)

// NewParams creates a new FeeParams object with specified parameters
func (tfp *FeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyResourceImageFee, tfp.Image, validateImage),
		paramtypes.NewParamSetPair(KeyCreateResourceJSONFee, tfp.Json, validateJSON),
		paramtypes.NewParamSetPair(KeyCreateResourceKeyFee, tfp.Default, validateDefault),
		paramtypes.NewParamSetPair(KeyBurnFactor, tfp.BurnFactor, validateBurnFactor),
	}
}
