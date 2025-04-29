package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidDocVersion(ctx context.Context, req *types.QueryDidDocVersionRequest) (*types.QueryDidDocVersionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	didDoc, err := k.GetDidDocVersion(ctx, req.Id, req.Version)
	if err != nil {
		return nil, err
	}

	return &types.QueryDidDocVersionResponse{Value: &didDoc}, nil
}
