package v6

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/did/exported"
	"github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, legacySubspace exported.Subspace,
	cdc codec.BinaryCodec, countCollection collections.Item[uint64],
	docCollection collections.Map[collections.Pair[string, string], types.DidDocWithMetadata],
) error {
	store := storeService.OpenKVStore(ctx)
	return migrateParams(store, cdc)
}

func migrateParams(store corestoretypes.KVStore, cdc codec.BinaryCodec) error {
	var legacyParams types.LegacyFeeParams
	bz, err := store.Get(types.ParamStoreKey)
	// Marshal and write to the new store
	if err != nil {
		return err
	}
	cdc.MustUnmarshal(bz, &legacyParams)

	// Now convert legacy to new format
	newParams := types.FeeParams{
		CreateDid: []types.FeeRange{
			{
				Denom:     legacyParams.CreateDid.Denom,
				MinAmount: &legacyParams.CreateDid.Amount,
				MaxAmount: util.PtrInt(legacyParams.CreateDid.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(700_000_000_000_000_000),
				MaxAmount: util.PtrInt(700_000_000_000_000_000),
			},
		},
		UpdateDid: []types.FeeRange{
			{
				Denom:     legacyParams.CreateDid.Denom,
				MinAmount: &legacyParams.UpdateDid.Amount,
				MaxAmount: util.PtrInt(legacyParams.UpdateDid.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(350_000_000_000_000_000),
				MaxAmount: util.PtrInt(350_000_000_000_000_000),
			},
		},
		DeactivateDid: []types.FeeRange{
			{
				Denom:     legacyParams.DeactivateDid.Denom,
				MinAmount: &legacyParams.DeactivateDid.Amount,
				MaxAmount: util.PtrInt(legacyParams.DeactivateDid.Amount.Mul(math.NewInt(2)).Int64()),
			},
			{
				Denom:     oracletypes.UsdDenom,
				MinAmount: util.PtrInt(140_000_000_000_000_000),
				MaxAmount: util.PtrInt(140_000_000_000_000_000),
			},
		},
		BurnFactor: legacyParams.BurnFactor,
	}

	// Marshal and write to the new store
	bz, err = cdc.Marshal(&newParams)
	if err != nil {
		return err
	}

	err = store.Set(types.ParamStoreKey, bz)
	if err != nil {
		return err
	}
	return nil
}
