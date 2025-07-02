package v5

import (
	corestoretypes "cosmossdk.io/core/store"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService,
	cdc codec.BinaryCodec,
) error {
	store := storeService.OpenKVStore(ctx)
	return migrateParams(store, cdc)
}

func migrateParams(store corestoretypes.KVStore, cdc codec.BinaryCodec) error {
	var legacyParams types.LegacyFeeParams

	bz, err := store.Get(types.ParamStoreKeyFeeParams)
	if err != nil {
		return err
	}
	cdc.MustUnmarshal(bz, &legacyParams)

	// Now convert legacy to new format
	usdParams := types.DefaultUSDParams()

	newParams := types.FeeParams{
		BurnFactor: legacyParams.BurnFactor,
		Default:    usdParams.Default,
		Json:       usdParams.Json,
		Image:      usdParams.Image,
	}

	// Marshal and write
	bz, err = cdc.Marshal(&newParams)
	if err != nil {
		return err
	}

	return store.Set(types.ParamStoreKeyFeeParams, bz)
}
