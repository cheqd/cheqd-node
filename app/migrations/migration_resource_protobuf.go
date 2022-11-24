package migrations

import (
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceProtobuf(sctx sdk.Context, mctx MigrationContext) error {
	var headerKeys []IteratorKey
	// Reset counter
	countStore := sctx.KVStore(mctx.resourceStoreKey)
	countKey := didutils.StrBytes(resourcetypes.ResourceCountKey)
	countStore.Delete(countKey)

	// Storages for old headers and data
	headerStore := sctx.KVStore(mctx.resourceStoreKey)
	dataStore := sctx.KVStore(mctx.resourceStoreKey)

	headerKeys = CollectAllKeys(sctx, mctx.resourceStoreKey, didutils.StrBytes(resourcetypesv1.ResourceHeaderKey))

	for _, headerKey := range headerKeys {
		dataKey := ResourceV1HeaderkeyToDataKey(headerKey)

		var headerV1 resourcetypesv1.ResourceHeader
		var dataV1 []byte

		mctx.codec.MustUnmarshal(headerStore.Get(headerKey), &headerV1)
		dataV1 = dataStore.Get(dataKey)

		newMetadata := resourcetypes.Metadata{
			CollectionId:      headerV1.CollectionId,
			Id:                headerV1.Id,
			Name:              headerV1.Name,
			Version:           "",
			ResourceType:      headerV1.ResourceType,
			AlsoKnownAs:       []*resourcetypes.AlternativeUri{},
			MediaType:         headerV1.MediaType,
			Created:           headerV1.Created,
			Checksum:          headerV1.Checksum,
			PreviousVersionId: headerV1.PreviousVersionId,
			NextVersionId:     headerV1.NextVersionId,
		}

		resourceWithMetadata := resourcetypes.ResourceWithMetadata{
			Metadata: &newMetadata,
			Resource: &resourcetypes.Resource{
				Data: dataV1,
			},
		}

		// Remove old resource data and header
		headerStore.Delete(headerKey)
		dataStore.Delete(dataKey)

		// Write new resource
		err := mctx.resourceKeeper.SetResource(&sctx, &resourceWithMetadata)
		if err != nil {
			return err
		}
	}
	return nil
}
