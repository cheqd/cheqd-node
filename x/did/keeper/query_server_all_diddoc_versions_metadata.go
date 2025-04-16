package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllDidDocVersionsMetadata(ctx context.Context, req *types.QueryAllDidDocVersionsMetadataRequest) (*types.QueryAllDidDocVersionsMetadataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	versions, err := k.GetAllDidDocVersions(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryAllDidDocVersionsMetadataResponse{Versions: versions}, nil
}
