package migrations

import (
	"crypto/sha256"
	"errors"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration for the whole did module
func MigrateDidV1(sctx sdk.Context, mctx MigrationContext) error {
	//
	err := MigrateDidProtobufV1(sctx, mctx)
	if err != nil {
		return err
	}

	err = MigrateDidUUIDV2(sctx, mctx)
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
	err := MigrateResourceChecksumV2(sctx, mctx)
	if err != nil {
		return err
	}

	// Resource Version Links migration
	err = MigrateResourceVersionLinksV2(sctx, mctx)
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
	var didKeys []IteratorKey

	ir := codectypes.NewInterfaceRegistry()

	ir.RegisterInterface("StateValueData", (*didtypesV1.StateValueData)(nil))
	ir.RegisterImplementations((*didtypesV1.StateValueData)(nil), &didtypesV1.Did{})

	CdcV1 := codec.NewProtoCodec(ir)

	didKeys = CollectAllKeys(sctx, mctx.didStoreKey, StrBytes(didtypesV1.DidKey))

	store := prefix.NewStore(
		sctx.KVStore(mctx.didStoreKey),
		StrBytes(didtypesV1.DidKey))

	for _, didKey := range didKeys {
		var stateValue didtypesV1.StateValue
		var newDidDocWithMetadata didtypes.DidDocWithMetadata
		CdcV1.MustUnmarshal(store.Get(didKey), &stateValue)

		newDidDocWithMetadata, err := StateValueToDIDDocWithMetadata(&stateValue)

		if err != nil {
			return err
		}

		// Remove old DID Doc
		store.Delete(didKey)

		// Set new DID Doc
		err = mctx.didKeeper.AddNewDidDocVersion(&sctx, &newDidDocWithMetadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func MigrateDidProtobufResourceV1(sctx sdk.Context, mctx MigrationContext) error {

	var headerKeys []IteratorKey
	// Reset counter
	countStore := sctx.KVStore(mctx.resourceStoreKey)
	countKey := didutils.StrBytes(resourcetypes.ResourceCountKey)
	countStore.Delete(countKey)

	// Storages for old headers and data
	headerStore := sctx.KVStore(mctx.resourceStoreKey)
	dataStore := sctx.KVStore(mctx.resourceStoreKey)

	headerKeys = CollectAllKeys(sctx, mctx.resourceStoreKey, StrBytes(resourcetypesv1.ResourceHeaderKey))

	for _, headerKey := range headerKeys {
		// ToDo: Make it more readable and understandable.
		// For now it's because Dids were set using just id as a key, but resources used 2 storages with prefixes for keys
		headerKey := []byte(resourcetypesv1.ResourceHeaderKey + string(headerKey))
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

	store := prefix.NewStore(
		sctx.KVStore(mctx.didStoreKey),
		StrBytes(didtypesV1.DidKey))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})

	closeIteratorOrPanic(iterator)

	// for ; iterator.Valid(); iterator.Next() {
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
	

	for metadataIterator.Valid() {

		var metadata resourcetypes.Metadata

		// Get metadata and data from storage
		mctx.codec.MustUnmarshal(metadataIterator.Value(), &metadata)

		// Get corresponding DidDoc

		metadata.Id = IndyStyleId(metadata.Id)
		metadata.CollectionId = IndyStyleId(metadata.CollectionId)

		// Update HeaderInfo
		err := mctx.resourceKeeper.UpdateResourceMetadata(&sctx, &metadata)
		if err != nil {
			return err
		}

		// Iterate next
		metadataIterator.Next()
	}
	return nil
}
