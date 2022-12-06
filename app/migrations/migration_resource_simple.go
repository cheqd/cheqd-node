package migrations

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceSimple(sctx sdk.Context, mctx MigrationContext, apply func(resourceWithMetadata *resourcetypes.ResourceWithMetadata)) error {
	store := sctx.KVStore(mctx.resourceStoreKey)

	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Cache resources
	var metadatas []resourcetypes.Metadata

	mctx.resourceKeeperNew.IterateAllResourceMetadatas(&sctx, func(metadata resourcetypes.Metadata) bool {
		metadatas = append(metadatas, metadata)
		return true
	})

	// Iterate and migrate resources
	for _, metadata := range metadatas {
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
	}

	return nil
}
