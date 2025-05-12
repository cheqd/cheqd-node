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

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, legacySubspace exported.Subspace,
	cdc codec.BinaryCodec, countCollection collections.Item[uint64],
) error {
	store := storeService.OpenKVStore(ctx)
	if err := migrateParams(ctx, store, legacySubspace, cdc); err != nil {
		return err
	}

	return migrateResourceCount(ctx, runtime.KVStoreAdapter(store), countCollection)
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
