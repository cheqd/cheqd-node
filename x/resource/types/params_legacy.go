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
func (p *FeeParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{

		paramtypes.NewParamSetPair(KeyResourceImageFee, p.Image, validateImage),
		paramtypes.NewParamSetPair(KeyCreateResourceJSONFee, p.Json, validateJSON),
		paramtypes.NewParamSetPair(KeyCreateResourceKeyFee, p.Default, validateDefault),
		paramtypes.NewParamSetPair(KeyBurnFactor, p.BurnFactor, validateBurnFactor),
	}
}
