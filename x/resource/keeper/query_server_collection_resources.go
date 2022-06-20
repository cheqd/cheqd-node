package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m queryServer) CollectionResources(c context.Context, request *types.QueryGetCollectionResourcesRequest) (*types.QueryGetCollectionResourcesResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	resources := m.GetResourceCollection(&ctx, request.CollectionId)

	return &types.QueryGetCollectionResourcesResponse{
		Resources: resources,
	}, nil
}
