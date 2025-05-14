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
) error {
	store := storeService.OpenKVStore(ctx)
	if err := migrateParams(ctx, store, legacySubspace, cdc); err != nil {
		return err
	}

	return migrateDidCount(ctx, runtime.KVStoreAdapter(store), countCollection)
}

func migrateParams(ctx sdk.Context, store corestoretypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var currParams types.FeeParams
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
