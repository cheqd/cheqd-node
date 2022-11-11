package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllDidDocVersionsMetadata(c context.Context, req *types.QueryGetAllDidDocVersionsMetadataRequest) (*types.QueryGetAllDidDocVersionsMetadataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	versions, err := k.GetAllDidDocVersions(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllDidDocVersionsMetadataResponse{Versions: versions}, nil
}
