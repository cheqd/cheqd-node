package keeper

import (
	"context"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) CollectionResources(c context.Context, req *types.QueryCollectionResourcesRequest) (*types.QueryCollectionResourcesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	req.Normalize()

	// Validate corresponding DIDDoc exists
	namespace := q.didKeeper.GetDidNamespace(&ctx)
	did := didutils.JoinDID(didtypes.DidMethod, namespace, req.CollectionId)
	if !q.didKeeper.HasDidDoc(&ctx, did) {
		return nil, didtypes.ErrDidDocNotFound.Wrap(did)
	}

	// Get all resources
	resources, err := q.GetResourceCollection(&ctx, req.CollectionId)
	if err != nil {
		return nil, types.ErrResourceNotAvail.Wrap(err.Error())
	}

	return &types.QueryCollectionResourcesResponse{
		Resources: resources,
	}, nil
}
