package v5

import (
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"

	storetypes "cosmossdk.io/store/types"
	"github.com/cheqd/cheqd-node/x/did/exported"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, legacySubspace exported.Subspace,
	cdc codec.BinaryCodec, countCollection collections.Item[uint64],
	docCollection collections.Map[collections.Pair[string, string], types.DidDocWithMetadata],
) error {
	store := storeService.OpenKVStore(ctx)
	if err := migrateParams(ctx, store, legacySubspace, cdc); err != nil {
		return err
	}

	kvStore := runtime.KVStoreAdapter(store)

	if err := migrateDidCount(ctx, kvStore, countCollection); err != nil {
		return err
	}

	return migrateDidDocuments(ctx, kvStore, cdc, docCollection)
}

func migrateParams(ctx sdk.Context, store corestoretypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var legacyParams types.LegacyFeeParams
	// Protect against missing param key (which causes panic in Get)
	legacySubspace.Get(ctx, types.ParamStoreKey, &legacyParams)
	// Now convert legacy to new format
	newParams := types.FeeParams{
		CreateDid: []types.FeeRange{
			{
				Denom:     legacyParams.CreateDid.Denom,
				MinAmount: &legacyParams.CreateDid.Amount,
				MaxAmount: &legacyParams.CreateDid.Amount,
			},
		},
		UpdateDid: []types.FeeRange{
			{
				Denom:     legacyParams.CreateDid.Denom,
				MinAmount: &legacyParams.UpdateDid.Amount,
				MaxAmount: &legacyParams.UpdateDid.Amount,
			},
		},
		DeactivateDid: []types.FeeRange{
			{
				Denom:     legacyParams.DeactivateDid.Denom,
				MinAmount: &legacyParams.DeactivateDid.Amount,
				MaxAmount: &legacyParams.DeactivateDid.Amount,
			},
		},
		BurnFactor: legacyParams.BurnFactor,
	}

	// Marshal and write to the new store
	bz, err := cdc.Marshal(&newParams)
	if err != nil {
		return err
	}

	err = store.Set(types.ParamStoreKey, bz)
	if err != nil {
		return err
	}
	return nil
}

func migrateDidCount(ctx sdk.Context, store storetypes.KVStore, countCollection collections.Item[uint64]) error {
	bz := store.Get([]byte(types.DidDocCountKey))
	if bz == nil {
		return nil
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot decode did doc count")
	}

	return countCollection.Set(ctx, count)
}

func migrateDidDocuments(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec,
	docCollection collections.Map[collections.Pair[string, string], types.DidDocWithMetadata],
) error {
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.DidDocVersionKey))

	for ; iterator.Valid(); iterator.Next() {
		var didDoc types.DidDocWithMetadata
		cdc.MustUnmarshal(iterator.Value(), &didDoc)

		// set document in collection
		if err := docCollection.Set(ctx, collections.Join(didDoc.DidDoc.Id, didDoc.Metadata.VersionId), didDoc); err != nil {
			return err
		}

		// delete old record
		store.Delete(iterator.Key())
	}

	return nil
}
