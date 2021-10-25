package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CredDef(c context.Context, req *types.QueryGetCredDefRequest) (*types.QueryGetCredDefResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	state, err := k.GetCredDef(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetCredDefResponse{CredDef: state.GetCredDef()}, nil
}
