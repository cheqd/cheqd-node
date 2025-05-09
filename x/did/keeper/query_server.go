package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := k.Paramstore.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
