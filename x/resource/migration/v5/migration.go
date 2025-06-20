package v5

import (
	corestoretypes "cosmossdk.io/core/store"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService,
	cdc codec.BinaryCodec) error {
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
	newParams := types.FeeParams{
		Image: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Image.Denom,
				MinAmount: &legacyParams.Image.Amount,
				MaxAmount: &legacyParams.Image.Amount,
			},
		},
		Json: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Json.Denom,
				MinAmount: &legacyParams.Json.Amount,
				MaxAmount: &legacyParams.Json.Amount,
			},
		},
		Default: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Default.Denom,
				MinAmount: &legacyParams.Default.Amount,
				MaxAmount: &legacyParams.Default.Amount,
			},
		},
		BurnFactor: legacyParams.BurnFactor,
	}

	// Marshal and write
	bz, err = cdc.Marshal(&newParams)
	if err != nil {
		return err
	}

	return store.Set(types.ParamStoreKeyFeeParams, bz)
}
