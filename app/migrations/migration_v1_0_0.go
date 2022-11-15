package migrations

import (
	"crypto/sha256"
	. "github.com/cheqd/cheqd-node/x/did/utils"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesV1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesV1 "github.com/cheqd/cheqd-node/x/resource/types/v1"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateCheqdV1(sctx sdk.Context, mctx MigrationContext) error {
	// Resource Checksum migration
	err := MigrateDidV1(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateResourceV1(sctx, mctx)
	if err != nil {
		return err
	}
	// TODO: Add more migrations for resource module
	return nil
}

// Migration for the whole did module
func MigrateDidV1(sctx sdk.Context, mctx MigrationContext) error {
	// 
	err := MigrateDidProtobufV1(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateDidUUIDV1(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateDidIndyStyleIdsV1(sctx, mctx)
	if err != nil {
		return err
	}

	return nil
}

// Resoure migrations
func MigrateResourceV1(sctx sdk.Context, mctx MigrationContext) error {
	// Resource Checksum migration
	err := MigrateResourceChecksumV1(sctx, mctx)
	if err != nil {
		return err
	}

	// Resource Version Links migration
	err = MigrateResourceVersionLinksV1(sctx, mctx)
	if err != nil {
		return err
	}
	// TODO: Add more migrations for resource module
	return nil
}

// Migration because of protobuf changes 
func MigrateDidProtobufV1(sctx sdk.Context, mctx MigrationContext) error {
	err := MigrateDidProtobufDIDocV1(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateDidProtobufResourceV1(sctx, mctx)
	if err != nil {
		return err
	}
	return nil
}

func MigrateDidProtobufDIDocV1(sctx sdk.Context, mctx MigrationContext) error {
	var iterator sdk.Iterator
	var stateValue didtypesV1.StateValue
	// var err error

	store := prefix.NewStore(
		sctx.KVStore(sdk.NewKVStoreKey(didtypesV1.StoreKey)), 
		StrBytes(didtypesV1.DidKey))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})

	closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		mctx.codec.MustUnmarshal(iterator.Value(), &stateValue)

		newDidDocWithMetadata, err := StateValueToDIDDocWithMetadata(stateValue)
		if err != nil {
			return err
		}
		
		// Remove old DID Doc
		store.Delete(iterator.Key())

		// Set new DID Doc
		mctx.didKeeper.AddNewDidDocVersion(&sctx, &newDidDocWithMetadata)
	}

	return nil
}

func MigrateDidProtobufResourceV1(sctx sdk.Context, mctx MigrationContext) error {
	// Reset counter
	countStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))
	countKey := didutils.StrBytes(resourcetypes.ResourceCountKey)
	countStore.Delete(countKey)

	// Storages for old headers and data
	headerStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))
	dataStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))

	// Iterators for old headers and data
	headerIterator := sdk.KVStorePrefixIterator(
		headerStore, 
		didutils.StrBytes(resourcetypesV1.ResourceHeaderKey))
	dataIterator := sdk.KVStorePrefixIterator(
		dataStore, 
		didutils.StrBytes(resourcetypesV1.ResourceDataKey))
	
	closeIteratorOrPanic(headerIterator)
	closeIteratorOrPanic(dataIterator)

	for headerIterator.Valid() {
		if !dataIterator.Valid() {
			panic("number of headers and data don't match")
		}

		var headerV1 resourcetypesV1.ResourceHeader
		var dataV1 []byte
		
		mctx.codec.MustUnmarshal(headerIterator.Value(), &headerV1)
		dataV1 = dataIterator.Value()

		newMetadata := resourcetypes.Metadata{

			CollectionId: 		headerV1.CollectionId,
			Id: 				headerV1.Id,
			Name: 				headerV1.Name,
			Version: 			"",
			ResourceType: 		headerV1.ResourceType,
			AlsoKnownAs:		[]*resourcetypes.AlternativeUri{},
			MediaType: 			headerV1.MediaType,
			Created: 			headerV1.Created,
			Checksum: 			headerV1.Checksum,
			PreviousVersionId: 	headerV1.PreviousVersionId,
			NextVersionId: 		headerV1.NextVersionId,

		}

		resourceWithMetadata := resourcetypes.ResourceWithMetadata{
			Metadata: &newMetadata,
			Resource: &resourcetypes.Resource{
				Data: dataV1,
			},
		}

		// Remove old resource data and header
		headerStore.Delete(headerIterator.Key())
		dataStore.Delete(dataIterator.Key())

		// Write new resource
		mctx.resourceKeeper.SetResource(&sctx, &resourceWithMetadata)

		// Iterate next
		headerIterator.Next()
		dataIterator.Next()
	}
	return nil
}

// Migration because we need to fix the algo for checksum calculation
func MigrateResourceChecksumV1(sctx sdk.Context, mctx MigrationContext) error {
	metadataStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))
	dataStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))
	metadataIterator := sdk.KVStorePrefixIterator(
		metadataStore, 
		didutils.StrBytes(resourcetypesV1.ResourceHeaderKey))
	dataIterator := sdk.KVStorePrefixIterator(
		dataStore, 
		didutils.StrBytes(resourcetypesV1.ResourceDataKey))
	

	closeIteratorOrPanic(metadataIterator)
	closeIteratorOrPanic(dataIterator)

	for metadataIterator.Valid() {
		if !dataIterator.Valid() {
			panic("number of headers and data don't match")
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
		mctx.resourceKeeper.UpdateResourceMetadata(&sctx, &metadata)

		// Iterate next
		metadataIterator.Next()
		dataIterator.Next()
	}
	return nil
}

// Recalculate resource links
func MigrateResourceVersionLinksV1(sctx sdk.Context, mctx MigrationContext) error {
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
func MigrateDidUUIDV1(sctx sdk.Context, mctx MigrationContext) error {
	
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
	var iterator sdk.Iterator
	// var err error

	store := prefix.NewStore(
		sctx.KVStore(sdk.NewKVStoreKey(didtypesV1.StoreKey)), 
		StrBytes(didtypesV1.DidKey))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})

	closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		didDocWithMetadata = didtypes.DidDocWithMetadata{}

		mctx.codec.MustUnmarshal(iterator.Value(), &didDocWithMetadata)

		// Make all dids indy style
		MoveToIndyStyleIds(&didDocWithMetadata)

		// Remove old DID Doc
		store.Delete(iterator.Key())

		// Set new DID Doc
		mctx.didKeeper.AddNewDidDocVersion(&sctx, &didDocWithMetadata)
	}

	return nil
}

func MigrateDidIndyStyleIdsV1ResourceModule(sctx sdk.Context, mctx MigrationContext) error {
	metadataStore := sctx.KVStore(sdk.NewKVStoreKey(resourcetypesV1.StoreKey))
	metadataIterator := sdk.KVStorePrefixIterator(
		metadataStore, 
		didutils.StrBytes(resourcetypesV1.ResourceHeaderKey))
	

	closeIteratorOrPanic(metadataIterator)

	for metadataIterator.Valid() {

		var metadata resourcetypes.Metadata

		// Get metadata and data from storage
		mctx.codec.MustUnmarshal(metadataIterator.Value(), &metadata)

		// Get corresponding DidDoc

		metadata.Id = IndyStyleId(metadata.Id)
		metadata.CollectionId = IndyStyleId(metadata.CollectionId)

		// Update HeaderInfo
		mctx.resourceKeeper.UpdateResourceMetadata(&sctx, &metadata)

		// Iterate next
		metadataIterator.Next()
	}
	return nil
}
