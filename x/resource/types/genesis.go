package types

import (
	"fmt"
)

const (
	DefaultCreateResourceImageFee   = 5e9
	DefaultCreateResourceJsonFee    = 2e9
	DefaultCreateResourceDefaultFee = 1e9
	DefaultBurnFactorRepresentation = 0.500000000000000000 // 0.5 or 50%
	_Precision                      = 1
	_PrecisionFactor                = 1e1 // CONTRACT: 1e(`_Precision`) <-- `sdk.Dec(1 <= `gs.BurnFactor` < `_PrecisionFactor`, `_Precision`).
)

// Bump `_Precision` if more decimals are needed, along with the exponent.
// e.g. `DefaultBurnFactor = 0.510000000000000000` --> `_Precision = 2` and `_PrecisionFactor = 1e2`, etc.

// DefaultGenesis returns the default `resource` genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ResourceList: []*Resource{},
		FeeParams:    DefaultFeeParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	resourceIdMap := make(map[string]bool)

	for _, resource := range gs.ResourceList {
		collectionResourceId := resource.Header.CollectionId + ":" + resource.Header.Id

		if _, ok := resourceIdMap[collectionResourceId]; ok {
			return fmt.Errorf("duplicated id for resource within the same collection: %s", collectionResourceId)
		}

		resourceIdMap[collectionResourceId] = true
	}

	if err := gs.FeeParams.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
