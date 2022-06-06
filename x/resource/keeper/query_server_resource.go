package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Resource(c context.Context, req *types.QueryGetResourceRequest) (*types.QueryGetResourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	resource, err := k.GetResource(&ctx, req.CollectionId, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetResourceResponse{Resource: &resource}, nil
}

func (k Keeper) CollectionResources(c context.Context, req *types.QueryGetCollectionResourcesRequest) (*types.QueryGetCollectionResourcesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	// stateResource, err := k.GetAllResource(&ctx, req.CollectionId)
	// if err != nil {
	// 	return nil, err
	// }

	// resource, err := stateResource.UnpackDataAsResource()
	// if err != nil {
	// 	return nil, err
	// }

	return &types.QueryGetCollectionResourcesResponse{Resources: []*types.Resource{}}, nil
}

func (k Keeper) AllResourceVersions(c context.Context, req *types.QueryGetAllResourceVersionsRequest) (*types.QueryGetAllResourceVersionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	// stateResource, err := k.GetResource(&ctx, req.CollectionId, req.Id)
	// if err != nil {
	// 	return nil, err
	// }

	// resource, err := stateResource.UnpackDataAsResource()
	// if err != nil {
	// 	return nil, err
	// }

	return &types.QueryGetAllResourceVersionsResponse{Resources: []*types.Resource{}}, nil
}
