package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
