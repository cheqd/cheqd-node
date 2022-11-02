package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Resources: []*ResourceWithMetadata{},
		FeeParams: &FeeParams{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.ValidateNoDuplicates()
}

func (gs GenesisState) ValidateNoDuplicates() error {
	// Group resources by collection
	resourcesByCollection := make(map[string][]*ResourceWithMetadata)

	for _, resource := range gs.Resources {
		existing := resourcesByCollection[resource.Metadata.CollectionId]
		resourcesByCollection[resource.Metadata.CollectionId] = append(existing, resource)
	}

	// Check that there are no collisions within each collection
	for _, resources := range resourcesByCollection {
		resourceIdMap := make(map[string]bool)

		for _, resource := range resources {
			if _, ok := resourceIdMap[resource.Metadata.Id]; ok {
				return fmt.Errorf("duplicated id for resource within the same collection. collection: %s, id: %s", resource.Metadata.CollectionId, resource.Metadata.Id)
			}

			resourceIdMap[resource.Metadata.Id] = true
		}
	}

	return nil
}
