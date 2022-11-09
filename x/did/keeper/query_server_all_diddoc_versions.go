package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllDidDocVersions(c context.Context, req *types.QueryGetAllDidDocVersionsRequest) (*types.QueryGetAllDidDocVersionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	versions, err := k.GetAllDidDocVersions(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllDidDocVersionsResponse{Versions: versions}, nil
}
