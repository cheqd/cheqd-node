package v4

import (
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

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
	latestResourceVersionCollection collections.Map[collections.Triple[string, string, string], string],
) error {
	store := storeService.OpenKVStore(ctx)
	if err := migrateParams(ctx, store, legacySubspace, cdc); err != nil {
		return err
	}

	kvStore := runtime.KVStoreAdapter(store)

	if err := migrateResourceCount(ctx, kvStore, countCollection); err != nil {
		return err
	}

	return migrateResources(ctx, kvStore, cdc, metadataCollection, dataCollection, latestResourceVersionCollection)
}

func migrateParams(ctx sdk.Context, store corestoretypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var currParams types.FeeParams
	legacySubspace.Get(ctx, types.ParamStoreKeyFeeParams, &currParams)

	if err := currParams.ValidateBasic(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&currParams)
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
	latestResourceVersionCollection collections.Map[collections.Triple[string, string, string], string],
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

		if metadata.NextVersionId == "" {
			if err := latestResourceVersionCollection.Set(ctx, collections.Join3(metadata.CollectionId, metadata.Name, metadata.ResourceType), metadata.Id); err != nil {
				return err
			}
		}

		// delete old record
		store.Delete(iterator.Key())
	}

	return nil
}
