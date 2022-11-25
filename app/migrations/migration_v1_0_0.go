package migrations

import (
	"crypto/sha256"
	"errors"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksumV2(sctx sdk.Context, mctx MigrationContext) error {
	metadataStore := sctx.KVStore(mctx.resourceStoreKey)
	dataStore := sctx.KVStore(mctx.resourceStoreKey)
	metadataIterator := sdk.KVStorePrefixIterator(
		metadataStore,
		didutils.StrBytes(resourcetypes.ResourceMetadataKey))
	dataIterator := sdk.KVStorePrefixIterator(
		dataStore,
		didutils.StrBytes(resourcetypes.ResourceDataKey))

	closeIteratorOrPanic(metadataIterator)
	closeIteratorOrPanic(dataIterator)

	for metadataIterator.Valid() {
		if !dataIterator.Valid() {
			return errors.New("number of headers and data don't match")
		}

		var metadata resourcetypes.Metadata
		var data []byte

		// Get metadata and data from storage
		mctx.codec.MustUnmarshal(metadataIterator.Value(), &metadata)
		data = dataIterator.Value()

		// Fix checksum
		checksum := sha256.Sum256(data)
		metadata.Checksum = checksum[:]

		// Update HeaderInfo
		err := mctx.resourceKeeper.UpdateResourceMetadata(&sctx, &metadata)
		if err != nil {
			return err
		}

		// Iterate next
		metadataIterator.Next()
		dataIterator.Next()
	}
	return nil
}

// Recalculate resource links
func MigrateResourceVersionLinksV2(sctx sdk.Context, mctx MigrationContext) error {
	// 	// TODO: We have to reset links first. Then we can use GetLastResourceVersionHeader
	// 	// but only because resources in state are corted by creation time.
	// 	// Also, we need to avoid loading all resources in memory.

	// 	headerIterator := resourceKeeper.GetHeaderIterator(&ctx)

	// 	defer resourcekeeper.CloseIteratorOrPanic(headerIterator)

	// 	for headerIterator.Valid() {
	// 		// Vars
	// 		var current_header resourcetypes.ResourceHeader

	// 		// Get the header
	// 		resourceKeeper.Cdc.MustUnmarshal(headerIterator.Value(), &current_header)

	// 		previousResourceVersionHeader, found := resourceKeeper.GetLastResourceVersionHeader(
	// 			&ctx,
	// 			current_header.CollectionId,
	// 			current_header.Name,
	// 			current_header.ResourceType)
	// 		if found {
	// 			// Set links
	// 			previousResourceVersionHeader.NextVersionId = current_header.Id
	// 			current_header.PreviousVersionId = previousResourceVersionHeader.Id

	// 			// Update previous version
	// 			err := resourceKeeper.UpdateResourceHeader(&ctx, &current_header)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			// Update previous version
	// 			err = resourceKeeper.UpdateResourceHeader(&ctx, &previousResourceVersionHeader)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}

	// 		headerIterator.Next()
	// 	}
	return nil
}

// Migration for making UUID case-insensitive
func MigrateDidUUIDV2(sctx sdk.Context, mctx MigrationContext) error {
	// 	var iterator sdk.Iterator
	// 	var stateValue didtypes.StateValue
	// 	var payload didtypes.MsgCreateDidPayload
	// 	var didDoc *didtypes.Did
	// 	var did didtypes.Did
	// 	var err error

	// 	iterator = didkeeper.GetStoreIterator(ctx)

	// 	defer func(iterator sdk.Iterator) {
	// 		err := iterator.Close()
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 	}(iterator)

	// 	for ; iterator.Valid(); iterator.Next() {
	// 		didkeeper.MustUnmarshal(iterator.Value(), &stateValue)
	// 		didDoc, err = stateValue.UnpackDataAsDid()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		payload = didtypes.MsgCreateDidPayloadFromDid(didDoc)
	// 		payload.Normalize()
	// 		did = payload.ToDid()
	// 		stateValue, err = didtypes.NewStateValue(&did, stateValue.GetMetadata())
	// 		if err != nil {
	// 			return err
	// 		}
	// 		didkeeper.SetDid(&ctx, &stateValue)
	// 	}
	return nil
}

// Migration for making ids in Indy-style
func MigrateDidIndyStyleIdsV1(sctx sdk.Context, mctx MigrationContext) error {
	err := MigrateDidIndyStyleIdsV1DidModule(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateDidIndyStyleIdsV1ResourceModule(sctx, mctx)
	if err != nil {
		return err
	}
	return nil
}

func MigrateDidIndyStyleIdsV1DidModule(sctx sdk.Context, mctx MigrationContext) error {
	// This migration should be run after protobuf that's why we use new DidDocWithMetadata
	var didDocWithMetadata didtypes.DidDocWithMetadata
	var didKeys []IteratorKey

	didKeys = CollectAllKeys(
		sctx,
		mctx.didStoreKey,
		didtypes.GetLatestDidDocVersionPrefix())

	store := sctx.KVStore(mctx.didStoreKey)

	for _, didKey := range didKeys {
		didDocWithMetadata = didtypes.DidDocWithMetadata{}

		mctx.codec.MustUnmarshal(store.Get(didKey), &didDocWithMetadata)

		// Make all dids indy style
		MoveToIndyStyleIds(&didDocWithMetadata)

		// Remove old DID Doc
		store.Delete(didKey)

		// Set new DID Doc
		err := mctx.didKeeper.AddNewDidDocVersion(&sctx, &didDocWithMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func MigrateDidIndyStyleIdsV1ResourceModule(sctx sdk.Context, mctx MigrationContext) error {

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
