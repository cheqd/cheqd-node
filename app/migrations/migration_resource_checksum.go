package migrations

import (
	"crypto/sha256"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksum(sctx sdk.Context, mctx MigrationContext) error {
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
		metadataKey := resourcetypes.GetResourceMetadataKey(metadata.CollectionId, metadata.Id)
		dataKey := resourcetypes.GetResourceDataKey(metadata.CollectionId, metadata.Id)

		// Read data
		data := store.Get(dataKey)

		// Remove old values
		store.Delete(metadataKey)
		store.Delete(dataKey)

		// Migrate, fix checksum
		checksum := sha256.Sum256(data)
		metadata.Checksum = checksum[:]

		// Write new value
		newResourceWithMetadata := resourcetypes.ResourceWithMetadata{
			Metadata: &metadata,
			Resource: &resourcetypes.Resource{
				Data: data,
			},
		}

		err := mctx.resourceKeeperNew.SetResource(&sctx, &newResourceWithMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}
