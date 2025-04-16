package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidDoc(ctx context.Context, req *types.QueryDidDocRequest) (*types.QueryDidDocResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	didDoc, err := k.GetLatestDidDoc(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryDidDocResponse{Value: &didDoc}, nil
}
