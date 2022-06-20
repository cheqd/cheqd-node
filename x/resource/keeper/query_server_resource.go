package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cheqd/cheqd-node/x/resource/types"
)

func (q queryServer) Resource(c context.Context, req *types.QueryGetResourceRequest) (*types.QueryGetResourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	resource, err := q.GetResource(&ctx, req.CollectionId, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetResourceResponse{
		Resource: &resource,
	}, nil
}
