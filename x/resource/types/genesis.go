package types

import (
	"fmt"
)

const DefaultResourceNamespace = "testnet"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ResourceList: []*Resource{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	resourceIdMap := make(map[string]bool)

	for _, resource := range gs.ResourceList {

		if _, ok := resourceIdMap[resource.Id]; ok {
			return fmt.Errorf("duplicated id for resource")
		}

		resourceIdMap[resource.Id] = true
	}

	return nil
}
