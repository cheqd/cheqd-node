package keeper

import (
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.FeeParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyFeeParams, &params)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.FeeParams) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyFeeParams, &params)
	return params
}
