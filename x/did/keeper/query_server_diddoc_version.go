package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidDocVersion(c context.Context, req *types.QueryGetDidDocVersionRequest) (*types.QueryGetDidDocVersionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	didDoc, err := k.GetDidDocVersion(&ctx, req.Id, req.Version)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetDidDocVersionResponse{Value: &didDoc}, nil
}
