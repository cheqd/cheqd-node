package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.FeeParams) error {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.FeeParamsKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(byteKey, bz)

	return nil
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.FeeParams) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.FeeParamsKey)
	bz := store.Get(byteKey)
	if bz == nil {
		return *types.DefaultFeeParams()
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetAuthority returns the x/cheqd module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
