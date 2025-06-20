package v6

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"github.com/cheqd/cheqd-node/x/did/exported"
	"github.com/cheqd/cheqd-node/x/did/types"
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
				MaxAmount: &legacyParams.CreateDid.Amount,
			},
		},
		UpdateDid: []types.FeeRange{
			{
				Denom:     legacyParams.CreateDid.Denom,
				MinAmount: &legacyParams.UpdateDid.Amount,
				MaxAmount: &legacyParams.UpdateDid.Amount,
			},
		},
		DeactivateDid: []types.FeeRange{
			{
				Denom:     legacyParams.DeactivateDid.Denom,
				MinAmount: &legacyParams.DeactivateDid.Amount,
				MaxAmount: &legacyParams.DeactivateDid.Amount,
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
