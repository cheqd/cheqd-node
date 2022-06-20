package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m queryServer) AllResourceVersions(c context.Context, req *types.QueryGetAllResourceVersionsRequest) (*types.QueryGetAllResourceVersionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	versions := m.GetAllResourceVersions(&ctx, req.CollectionId, req.Name, req.ResourceType, req.MimeType)

	return &types.QueryGetAllResourceVersionsResponse{
		Resources: versions,
	}, nil
}
