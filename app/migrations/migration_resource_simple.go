package migrations

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceSimple(sctx sdk.Context, mctx MigrationContext, apply func(resourceWithMetadata *resourcetypes.ResourceWithMetadata)) error {
	sctx.Logger().Debug("MigrateResourceSimple: Starting migration")
	store := sctx.KVStore(mctx.resourceStoreKey)

	sctx.Logger().Debug("MigrateResourceSimple: Removing old counters")
	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Cache resources
	var metadatas []resourcetypes.Metadata

	sctx.Logger().Debug("MigrateResourceSimple: Iterating over all resource metadatas")
	mctx.resourceKeeperNew.IterateAllResourceMetadatas(&sctx, func(metadata resourcetypes.Metadata) bool {
		metadatas = append(metadatas, metadata)
		return true
	})

	// Iterate and migrate resources
	for _, metadata := range metadatas {

		sctx.Logger().Debug("MigrateResourceSimple: Starting migration for resource: " + metadata.Id)
		// Read value
		resourceWithMetadata, err := mctx.resourceKeeperNew.GetResource(&sctx, metadata.CollectionId, metadata.Id)
		if err != nil {
			return err
		}

		// Remove old values
		metadataKey := resourcetypes.GetResourceMetadataKey(metadata.CollectionId, metadata.Id)
		dataKey := resourcetypes.GetResourceDataKey(metadata.CollectionId, metadata.Id)

		store.Delete(metadataKey)
		store.Delete(dataKey)

		// Migrate
		apply(&resourceWithMetadata)

		// Write new value
		err = mctx.resourceKeeperNew.SetResource(&sctx, &resourceWithMetadata)
		if err != nil {
			return err
		}

		sctx.Logger().Debug("MigrateResourceSimple: Migration finished for resource: " + resourceWithMetadata.Metadata.Id)
	}

	sctx.Logger().Debug("MigrateResourceSimple: Migration finished")

	return nil
}
