package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetResourceCount get the total number of resource
func (k Keeper) GetResourceCount(ctx context.Context) (uint64, error) {
	count, err := k.ResourceCount.Get(ctx)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// SetResourceCount set the total number of resource
func (k Keeper) SetResourceCount(ctx context.Context, count uint64) error {
	return k.ResourceCount.Set(ctx, count)
}

func (k Keeper) AddNewResourceVersion(ctx context.Context, resource *types.ResourceWithMetadata, previousVersionResourceId string) error {
	// Find previous version and upgrade backward and forward version links
	previousResourceVersionHeader, found, err := k.GetLastResourceVersionMetadata(ctx, resource.Metadata.CollectionId, resource.Metadata.Name, resource.Metadata.ResourceType, previousVersionResourceId)
	if err != nil {
		return err
	}
	if found {
		// Set links
		previousResourceVersionHeader.NextVersionId = resource.Metadata.Id
		resource.Metadata.PreviousVersionId = previousResourceVersionHeader.Id

		// Update previous version
		err := k.UpdateResourceMetadata(ctx, &previousResourceVersionHeader)
		if err != nil {
			return err
		}
	}

	err = k.SetResource(ctx, resource)
	return err
}

// SetResource create or update a specific resource in the store
func (k Keeper) SetResource(ctx context.Context, resource *types.ResourceWithMetadata) error {
	hasResource, err := k.ResourceMetadata.Has(ctx, collections.Join(resource.Metadata.CollectionId, resource.Metadata.Id))
	if err != nil {
		return err
	}

	if !hasResource {
		count, err := k.GetResourceCount(ctx)
		if err != nil {
			return err
		}
		if err := k.SetResourceCount(ctx, count+1); err != nil {
			return err
		}
	}

	// Set metadata
	if err := k.ResourceMetadata.Set(ctx, collections.Join(resource.Metadata.CollectionId, resource.Metadata.Id), *resource.Metadata); err != nil {
		return err
	}

	// Set data
	return k.ResourceData.Set(ctx, collections.Join(resource.Metadata.CollectionId, resource.Metadata.Id), resource.Resource.Data)
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx context.Context, collectionID string, id string) (types.ResourceWithMetadata, error) {
	hasResource, err := k.ResourceMetadata.Has(ctx, collections.Join(collectionID, id))
	if err != nil {
		return types.ResourceWithMetadata{}, err
	}

	if !hasResource {
		return types.ResourceWithMetadata{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionID + ":" + id)
	}
	// Get metadata
	metadata, err := k.ResourceMetadata.Get(ctx, collections.Join(collectionID, id))
	if err != nil {
		return types.ResourceWithMetadata{}, err
	}

	// Get data
	data, err := k.ResourceData.Get(ctx, collections.Join(collectionID, id))
	if err != nil {
		return types.ResourceWithMetadata{}, err
	}

	return types.ResourceWithMetadata{
		Metadata: &metadata,
		Resource: &types.Resource{Data: data},
	}, nil
}

func (k Keeper) GetResourceMetadata(ctx context.Context, collectionID string, id string) (types.Metadata, error) {
	hasResource, err := k.ResourceMetadata.Has(ctx, collections.Join(collectionID, id))
	if err != nil {
		return types.Metadata{}, err
	}
	if !hasResource {
		return types.Metadata{}, sdkerrors.ErrNotFound.Wrap("resource " + collectionID + ":" + id)
	}

	return k.ResourceMetadata.Get(ctx, collections.Join(collectionID, id))
}

// HasResource checks if the resource exists in the store
func (k Keeper) HasResource(ctx context.Context, collectionID string, id string) bool {
	has, err := k.ResourceMetadata.Has(ctx, collections.Join(collectionID, id))
	if err != nil {
		return false
	}
	return has
}

func (k Keeper) GetResourceCollection(ctx context.Context, collectionID string) ([]*types.Metadata, error) {
	var resources []*types.Metadata

	rng := collections.NewPrefixedPairRange[string, string](collectionID)
	err := k.ResourceMetadata.Walk(ctx, rng, func(_ collections.Pair[string, string], metadata types.Metadata) (bool, error) {
		metadataCopy := metadata // Create a copy to avoid reference issues
		resources = append(resources, &metadataCopy)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (k Keeper) GetLastResourceVersionMetadata(ctx context.Context, collectionID, name, resourceType string, previousVersionResourceId string) (types.Metadata, bool, error) {
	// Case 1: Use provided previousVersionResourceId
	if previousVersionResourceId != "" {
		lastVersion, err := k.ResourceMetadata.Get(ctx, collections.Join(collectionID, previousVersionResourceId))
		if err != nil {
			return types.Metadata{}, false, err
		}
		if lastVersion.NextVersionId == "" &&
			lastVersion.Name == name &&
			lastVersion.ResourceType == resourceType {
			return lastVersion, true, nil
		}
		return types.Metadata{}, false, nil
	}

	// Case 2: Resolve last version from chain (Expensive)
	var lastVersion types.Metadata
	found := false
	rng := collections.NewPrefixedPairRange[string, string](collectionID)
	err := k.ResourceMetadata.Walk(ctx, rng, func(_ collections.Pair[string, string], metadata types.Metadata) (bool, error) {
		if metadata.Name == name &&
			metadata.ResourceType == resourceType &&
			metadata.NextVersionId == "" {
			lastVersion = metadata
			found = true
			return true, nil // stop iteration
		}
		return false, nil // continue
	})
	if err != nil {
		return types.Metadata{}, false, err
	}
	return lastVersion, found, nil
}

// UpdateResourceMetadata update the metadata of a resource. Returns an error if the resource doesn't exist
func (k Keeper) UpdateResourceMetadata(ctx context.Context, metadata *types.Metadata) error {
	hasResource, err := k.ResourceMetadata.Has(ctx, collections.Join(metadata.CollectionId, metadata.Id))
	if err != nil {
		return err
	}
	if !hasResource {
		return sdkerrors.ErrNotFound.Wrap("resource " + metadata.CollectionId + ":" + metadata.Id)
	}

	return k.ResourceMetadata.Set(ctx, collections.Join(metadata.CollectionId, metadata.Id), *metadata)
}

func (k Keeper) IterateAllResourceMetadatas(ctx context.Context, callback func(metadata types.Metadata) (continue_ bool)) error {
	err := k.ResourceMetadata.Walk(
		ctx,
		nil, // nil range means full range in x/collections
		func(_ collections.Pair[string, string], metadata types.Metadata) (bool, error) {
			if !callback(metadata) {
				return true, nil
			}
			return false, nil
		},
	)

	return err
}

// GetAllResources returns all resources as a list
// Loads everything in memory. Use only for genesis export!
func (k Keeper) GetAllResources(ctx context.Context) (list []*types.ResourceWithMetadata, iterErr error) {
	var resources []*types.ResourceWithMetadata

	err := k.IterateAllResourceMetadatas(ctx, func(metadata types.Metadata) bool {
		resource, err := k.GetResource(ctx, metadata.CollectionId, metadata.Id)
		if err != nil {
			iterErr = err
			return false
		}

		resources = append(resources, &resource)
		return true
	})
	if err != nil {
		return nil, err
	}

	if iterErr != nil {
		return nil, iterErr
	}

	return resources, nil
}
