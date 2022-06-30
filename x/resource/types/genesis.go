package types

import (
	"fmt"
)

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
		collectionResourceId := resource.Header.CollectionId + ":" + resource.Header.Id

		if _, ok := resourceIdMap[collectionResourceId]; ok {
			return fmt.Errorf("duplicated id for resource within the same collection: %s", collectionResourceId)
		}

		resourceIdMap[collectionResourceId] = true
	}

	return nil
}
