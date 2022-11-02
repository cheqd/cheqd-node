package types

import (
	"fmt"
)

const DefaultDidNamespace = "testnet"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidDocs:      []*DidDocWithMetadata{},
		DidNamespace: DefaultDidNamespace,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateNoDuplicates()
	if err != nil {
		return err
	}

	return gs.ValidateBasic()
}

func (gs GenesisState) ValidateNoDuplicates() error {
	cache := make(map[string]bool)

	for _, didDoc := range gs.DidDocs {

		if _, ok := cache[didDoc.DidDoc.Id]; ok {
			return fmt.Errorf("duplicated didDoc found with id %s", didDoc.DidDoc.Id)
		}

		cache[didDoc.DidDoc.Id] = true
	}

	return nil
}

func (gs GenesisState) ValidateBasic() error {
	for _, didDoc := range gs.DidDocs {
		err := didDoc.DidDoc.Validate(nil)
		if err != nil {
			return err
		}
	}

	return nil
}
