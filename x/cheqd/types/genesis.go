package types

import (
	"fmt"
)

const DefaultDidNamespace = "testnet"
const DefaultCreateDidTxFee = 5e9
const DefaultUpdateDidTxFee = 2e9
const DefaultDeactivateDidTxFee = 1e9
const DefaultBurnFactorRepresentation = 0.500000000000000000 // 0.5 or 50%
const _Precision = 1
const _PrecisionFactor = 1e1 // CONTRACT: 1e(`_Precision`) <-- `sdk.Dec(1 <= `gs.BurnFactor` < `_PrecisionFactor`, `_Precision`).
// Bump `_Precision` if more decimals are needed, along with the exponent.
// e.g. `DefaultBurnFactor = 0.510000000000000000` --> `_Precision = 2` and `_PrecisionFactor = 1e2`, etc.

// DefaultGenesis returns the default `cheqd` genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidList:      []*StateValue{},
		DidNamespace: DefaultDidNamespace,
		FeeParams: DefaultFeeParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	didIdMap := make(map[string]bool)

	for _, elem := range gs.DidList {
		did, err := elem.UnpackDataAsDid()
		if err != nil {
			return err
		}

		if _, ok := didIdMap[did.Id]; ok {
			return fmt.Errorf("duplicate id for did")
		}

		didIdMap[did.Id] = true
	}

	if err := gs.FeeParams.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
