package types

import "fmt"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidList: []*StateValue{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	didIdMap := make(map[string]bool)

	for _, elem := range gs.DidList {
		did, err := elem.GetDid()
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
