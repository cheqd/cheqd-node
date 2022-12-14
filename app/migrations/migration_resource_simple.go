package migrations

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceSimple(sctx sdk.Context, mctx MigrationContext, apply func(resourceWithMetadata *resourcetypes.ResourceWithMetadata)) error {
	sctx.Logger().Debug("Simple migration for resources. Start")
	store := sctx.KVStore(mctx.resourceStoreKey)

	sctx.Logger().Debug("Simple migration for resources. Remove old counters")
	// Reset counter
	mctx.didKeeperNew.SetDidDocCount(&sctx, 0)

	// Cache resources
	var metadatas []resourcetypes.Metadata

	sctx.Logger().Debug("Simple migration for resources. Iterate over all resource metadatas")
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

		sctx.Logger().Debug("old resource with metadata: %s", string(mctx.codec.MustMarshalJSON(&resourceWithMetadata)))

		// Remove old values
		metadataKey := resourcetypes.GetResourceMetadataKey(metadata.CollectionId, metadata.Id)
		dataKey := resourcetypes.GetResourceDataKey(metadata.CollectionId, metadata.Id)

		store.Delete(metadataKey)
		store.Delete(dataKey)

		// Migrate
		apply(&resourceWithMetadata)

		sctx.Logger().Debug("new resource with metadata: %s", string(mctx.codec.MustMarshalJSON(&resourceWithMetadata)))

		// Write new value
		err = mctx.resourceKeeperNew.SetResource(&sctx, &resourceWithMetadata)
		if err != nil {
			return err
		}
	}

	sctx.Logger().Debug("Simple migration for resources. End")

	return nil
}
