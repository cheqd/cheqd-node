package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

const DefaultDidNamespace = "testnet"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidList:      []*StateValue{},
		DidNamespace: DefaultDidNamespace,
	}
}

func (gs *GenesisState) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, elem := range gs.DidList {
		err := elem.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
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
			return fmt.Errorf("duplicated id for did")
		}

		didIdMap[did.Id] = true
	}

	return nil
}
