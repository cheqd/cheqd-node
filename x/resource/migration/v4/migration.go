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
	var legacyParams types.LegacyFeeParams
	legacySubspace.Get(ctx, types.ParamStoreKeyFeeParams, &legacyParams)

	// Now convert legacy to new format
	newParams := types.FeeParams{
		Image: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Image.Denom,
				MinAmount: legacyParams.Image.Amount,
				MaxAmount: legacyParams.Image.Amount,
			},
		},
		Json: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Json.Denom,
				MinAmount: legacyParams.Json.Amount,
				MaxAmount: legacyParams.Json.Amount,
			},
		},
		Default: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Default.Denom,
				MinAmount: legacyParams.Default.Amount,
				MaxAmount: legacyParams.Default.Amount,
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
