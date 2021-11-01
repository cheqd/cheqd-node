package types

import "fmt"

const DidNamespace = "testnet"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidList:      []*StateValue{},
		DidNamespace: DidNamespace,
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

	if len(gs.GetDidNamespace()) == 0 {
		return fmt.Errorf("DidNamespace is not set up")
	}

	return nil
}
