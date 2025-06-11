package v4

import (
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/exported"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LegacyResourceDataKey returns the byte representation of legacy resource data key
func LegacyResourceDataKey(collectionID string, id string) []byte {
	return []byte(types.ResourceDataKey + collectionID + ":" + id)
}

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, legacySubspace exported.Subspace,
	cdc codec.BinaryCodec, countCollection collections.Item[uint64],
	metadataCollection collections.Map[collections.Pair[string, string], types.Metadata],
	dataCollection collections.Map[collections.Pair[string, string], []byte],
) error {
	store := storeService.OpenKVStore(ctx)
	if err := migrateParams(ctx, store, legacySubspace, cdc); err != nil {
		return err
	}

	kvStore := runtime.KVStoreAdapter(store)

	if err := migrateResourceCount(ctx, kvStore, countCollection); err != nil {
		return err
	}

	return migrateResources(ctx, kvStore, cdc, metadataCollection, dataCollection)
}

func migrateParams(ctx sdk.Context, store corestoretypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var legacyParams types.LegacyFeeParams
	legacySubspace.Get(ctx, types.ParamStoreKeyFeeParams, &legacyParams)

	// Now convert legacy to new format
	newParams := types.FeeParams{
		Image: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Image.Denom,
				MinAmount: legacyParams.Image.Amount,
				MaxAmount: &legacyParams.Image.Amount,
			},
		},
		Json: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Json.Denom,
				MinAmount: legacyParams.Json.Amount,
				MaxAmount: &legacyParams.Json.Amount,
			},
		},
		Default: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Default.Denom,
				MinAmount: legacyParams.Default.Amount,
				MaxAmount: &legacyParams.Default.Amount,
			},
		},
		BurnFactor: legacyParams.BurnFactor,
	}

	// Marshal and write
	bz, err := cdc.Marshal(&newParams)
	if err != nil {
		return err
	}

	return store.Set(types.ParamStoreKeyFeeParams, bz)
}

func migrateResourceCount(ctx sdk.Context, store storetypes.KVStore, countCollection collections.Item[uint64]) error {
	bz := store.Get([]byte(types.ResourceCountKey))
	if bz == nil {
		return nil
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot decode resource count")
	}

	return countCollection.Set(ctx, count)
}

func migrateResources(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec,
	metadataCollection collections.Map[collections.Pair[string, string], types.Metadata],
	dataCollection collections.Map[collections.Pair[string, string], []byte],
) error {
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.ResourceMetadataKey))

	for ; iterator.Valid(); iterator.Next() {
		var metadata types.Metadata
		cdc.MustUnmarshal(iterator.Value(), &metadata)

		// set resource metadata in metadata collection
		if err := metadataCollection.Set(ctx, collections.Join(metadata.CollectionId, metadata.Id), metadata); err != nil {
			return err
		}

		dataKey := LegacyResourceDataKey(metadata.CollectionId, metadata.Id)
		data := store.Get(dataKey)
		if data != nil {
			// set resource data in data collection
			if err := dataCollection.Set(ctx, collections.Join(metadata.CollectionId, metadata.Id), data); err != nil {
				return err
			}
			// delete old record
			store.Delete(dataKey)
		}

		// delete old record
		store.Delete(iterator.Key())
	}

	return nil
}
