package migrations

import (
	// "crypto/sha256"

	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	. "github.com/cheqd/cheqd-node/x/did/utils"

	// didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	// resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
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
	var stateValue didtypesv1.StateValue
	// var err error

	store := prefix.NewStore(
		sctx.KVStore(sdk.NewKVStoreKey(didtypes.StoreKey)), 
		StrBytes(didtypes.DidKey))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})

	closeIteratorOrPanic(iterator)

	for ; iterator.Valid(); iterator.Next() {
		mctx.codec.MustUnmarshal(iterator.Value(), &stateValue)

		newDidDocWithMetadata, err := StateValueToDIDDocWithMetadata(stateValue)
		if err != nil {
			return err
		}
		// Set new DID Doc
		mctx.didKeeper.SetDidDoc(&sctx, &newDidDocWithMetadata)
		
		// Remove old DID Doc
		store.Delete(iterator.Key())
	}

	return nil
}

func MigrateDidProtobufResourceV1(sctx sdk.Context, mctx MigrationContext) error {
	metadataIterator := sdk.KVStorePrefixIterator(
		sctx.KVStore(sdk.NewKVStoreKey(resourcetypes.StoreKey)), 
		resourcetypes.KeyPrefix(resourcetypes.ResourceMetadataKey))

	closeIteratorOrPanic(metadataIterator)

	for ; metadataIterator.Valid(); metadataIterator.Next() {

		var metadata resourcetypes.Metadata
		mctx.codec.MustUnmarshal(metadataIterator.Value(), &metadata)

		new_metadata := resourcetypes.Metadata{
			CollectionId: 		metadata.CollectionId,
			Id: 				metadata.Id,
			Name: 				metadata.Name,
			Version: 			metadata.Version,
			ResourceType: 		metadata.ResourceType,
			AlsoKnownAs:		metadata.AlsoKnownAs,
			MediaType: 			metadata.MediaType,
			Created: 			metadata.Created,
			Checksum: 			metadata.Checksum,
			PreviousVersionId: 	metadata.PreviousVersionId,
			NextVersionId: 		metadata.NextVersionId,

		}

		mctx.resourceKeeper.UpdateResourceMetadata(&sctx, &new_metadata)
	}

	return nil
	
}

func MigrateResourceChecksumV1(sctx sdk.Context, mctx MigrationContext) error {
// 	// TODO: Loading everything into memory is not the best approach.
// 	// Resources can be large. I would suggest to use iterator instead and load resources one by one.

// 	headerIterator := resourceKeeper.GetHeaderIterator(&ctx)
// 	store := ctx.KVStore(resourceKeeper.StoreKey)

// 	defer resourcekeeper.CloseIteratorOrPanic(headerIterator)

// 	for headerIterator.Valid() {
// 		// Vars
// 		var data_val []byte
// 		var header_val resourcetypes.ResourceHeader

// 		// Get the header
// 		resourceKeeper.Cdc.MustUnmarshal(headerIterator.Value(), &header_val)

// 		data_val = store.Get(resourcekeeper.GetResourceDataKeyBytes(header_val.CollectionId, header_val.Id))
// 		checksum := sha256.Sum256(data_val)
// 		header_val.Checksum = checksum[:]

// 		// Update header
// 		err := resourceKeeper.UpdateResourceHeader(&ctx, &header_val)
// 		if err != nil {
// 			return err
// 		}

// 		headerIterator.Next()
// 	}
	return nil
}

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

func MigrateDidIndyStyleIdsV1(sctx sdk.Context, mctx MigrationContext) error {
	return nil
}

