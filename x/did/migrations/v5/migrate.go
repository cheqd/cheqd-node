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
	var currParams types.LegacyFeeParams
	legacySubspace.Get(ctx, types.ParamStoreKey, &currParams)

	if err := currParams.ValidateBasic(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&currParams)
	if err != nil {
		return err
	}

	return store.Set(types.ParamStoreKey, bz)
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
