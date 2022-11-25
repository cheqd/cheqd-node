package migrations

import (
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceIndyStyle(sctx sdk.Context, mctx MigrationContext) error {
	var metadataKeys []IteratorKey

	store := sctx.KVStore(mctx.resourceStoreKey)
	metadataKeys = CollectAllKeys(
		sctx,
		mctx.resourceStoreKey,
		didutils.StrBytes(resourcetypes.ResourceMetadataKey))

	for _, metadataKey := range metadataKeys {

		var metadata resourcetypes.Metadata
		var data []byte

		dataKey := ResourceV2MetadataKeyToDataKey(metadataKey)

		// Get metadata and data from storage
		mctx.codec.MustUnmarshal(store.Get(metadataKey), &metadata)
		data = store.Get(dataKey)

		// Get corresponding DidDoc

		metadata.Id = IndyStyleId(metadata.Id)
		metadata.CollectionId = IndyStyleId(metadata.CollectionId)

		newResourceWithMetadata := resourcetypes.ResourceWithMetadata{
			Metadata: &metadata,
			Resource: &resourcetypes.Resource{
				Data: data,
			},
		}

		// Remove old values
		store.Delete(metadataKey)
		store.Delete(dataKey)

		// Update HeaderInfo
		err := mctx.resourceKeeper.SetResource(&sctx, &newResourceWithMetadata)
		if err != nil {
			return err
		}
	}
	return nil
}
