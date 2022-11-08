package migrations

import (
	"crypto/sha256"

	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateCheqdV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// TODO: implement for cheqd module
	return nil
}

func MigrateResourceV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// Resource Checksum migration
	err := MigrateResourceChecksumV1(ctx, cheqdKeeper, resourceKeeper)
	if err != nil {
		return err
	}

	err = MigrateResourceVersionLinksV1(ctx, cheqdKeeper, resourceKeeper)
	if err != nil {
		return err
	}
	// TODO: Add more migrations for resource module
	return nil
}

func MigrateResourceChecksumV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// TODO: Loading everything into memory is not the best approach.
	// Resources can be large. I would suggest to use iterator instead and load resources one by one.

	headerIterator := resourceKeeper.GetHeaderIterator(&ctx)
	store := ctx.KVStore(resourceKeeper.StoreKey)

	defer resourcekeeper.CloseIteratorOrPanic(headerIterator)

	for headerIterator.Valid() {
		// Vars
		var data_val []byte
		var header_val resourcetypes.ResourceHeader

		// Get the header
		resourceKeeper.Cdc.MustUnmarshal(headerIterator.Value(), &header_val)

		data_val = store.Get(resourcekeeper.GetResourceDataKeyBytes(header_val.CollectionId, header_val.Id))
		checksum := sha256.Sum256(data_val)
		header_val.Checksum = checksum[:]

		// Update header
		err := resourceKeeper.UpdateResourceHeader(&ctx, &header_val)
		if err != nil {
			return err
		}

		headerIterator.Next()
	}
	return nil
}

func MigrateResourceVersionLinksV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	// TODO: We have to reset links first. Then we can use GetLastResourceVersionHeader
	// but only because resources in state are corted by creation time.
	// Also, we need to avoid loading all resources in memory.

	headerIterator := resourceKeeper.GetHeaderIterator(&ctx)

	defer resourcekeeper.CloseIteratorOrPanic(headerIterator)

	for headerIterator.Valid() {
		// Vars
		var current_header resourcetypes.ResourceHeader

		// Get the header
		resourceKeeper.Cdc.MustUnmarshal(headerIterator.Value(), &current_header)

		previousResourceVersionHeader, found := resourceKeeper.GetLastResourceVersionHeader(
			&ctx, 
			current_header.CollectionId, 
			current_header.Name, 
			current_header.ResourceType)
		if found {
			// Set links
			previousResourceVersionHeader.NextVersionId = current_header.Id
			current_header.PreviousVersionId = previousResourceVersionHeader.Id

			// Update previous version
			err := resourceKeeper.UpdateResourceHeader(&ctx, &current_header)
			if err != nil {
				return err
			}

			// Update previous version
			err = resourceKeeper.UpdateResourceHeader(&ctx, &previousResourceVersionHeader)
			if err != nil {
				return err
			}
		}

		headerIterator.Next()
	}
	return nil
}

func MigrateCheqdUUIDV1(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	var iterator sdk.Iterator
	var stateValue cheqdtypes.StateValue
	var payload cheqdtypes.MsgCreateDidPayload
	var didDoc *cheqdtypes.Did
	var did cheqdtypes.Did
	var err error


	iterator = cheqdKeeper.GetStoreIterator(ctx)

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err.Error())
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		cheqdKeeper.MustUnmarshal(iterator.Value(), &stateValue)
		didDoc, err = stateValue.UnpackDataAsDid()
		if err != nil {
			return err
		}
		payload = cheqdtypes.MsgCreateDidPayloadFromDid(didDoc)
		payload.Normalize()
		did = payload.ToDid()
		stateValue, err = cheqdtypes.NewStateValue(&did, stateValue.GetMetadata())
		if err != nil {
			return err
		}
		cheqdKeeper.SetDid(&ctx, &stateValue)
	}
	return nil
}

