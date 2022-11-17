package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidDoc(c context.Context, req *types.QueryGetDidDocRequest) (*types.QueryGetDidDocResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	didDoc, err := k.GetLatestDidDoc(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetDidDocResponse{Value: &didDoc}, nil
}
