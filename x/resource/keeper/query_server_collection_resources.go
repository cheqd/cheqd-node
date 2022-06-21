package keeper

import (
	"context"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
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

	// Validate corresponding DIDDoc exists
	namespace := m.cheqdKeeper.GetDidNamespace(&ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, request.CollectionId)
	if !m.cheqdKeeper.HasDid(&ctx, did) {
		return nil, cheqdtypes.ErrDidDocNotFound.Wrap(did)
	}

	// Get all resources
	resources := m.GetResourceCollection(&ctx, request.CollectionId)

	return &types.QueryGetCollectionResourcesResponse{
		Resources: resources,
	}, nil
}
