package v5

import (
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/util"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
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
	newParams := types.FeeParams{
		Image: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Image.Denom,
				MinAmount: &legacyParams.Image.Amount,
				MaxAmount: util.PtrInt(legacyParams.Image.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(280_000_000_000_000_000),
				MaxAmount: util.PtrInt(280_000_000_000_000_000),
			},
		},
		Json: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Json.Denom,
				MinAmount: &legacyParams.Json.Amount,
				MaxAmount: util.PtrInt(legacyParams.Json.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(49_000_000_000_000_000),
				MaxAmount: util.PtrInt(49_000_000_000_000_000),
			},
		},
		Default: []didtypes.FeeRange{
			{
				Denom:     legacyParams.Default.Denom,
				MinAmount: &legacyParams.Default.Amount,
				MaxAmount: util.PtrInt(legacyParams.Default.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(84_000_000_000_000_000),
				MaxAmount: util.PtrInt(84_000_000_000_000_000),
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
