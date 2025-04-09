package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
)

func (k Keeper) SetParams(ctx context.Context, params types.FeeParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyFeeParams, &params)
}

func (k Keeper) GetParams(ctx context.Context) (params types.FeeParams) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyFeeParams, &params)
	return params
}
