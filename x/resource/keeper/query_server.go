package keeper

import (
	"context"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type queryServer struct {
	Keeper
	didKeeper didkeeper.Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper, cheqdKeeper didkeeper.Keeper) types.QueryServer {
	return &queryServer{
		Keeper:    keeper,
		didKeeper: cheqdKeeper,
	}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.ParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
