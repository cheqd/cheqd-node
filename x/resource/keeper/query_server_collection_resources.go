package keeper

import (
	"context"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) CollectionResources(ctx context.Context, req *types.QueryCollectionResourcesRequest) (*types.QueryCollectionResourcesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	// Validate corresponding DIDDoc exists
	namespace, err := q.didKeeper.GetDidNamespace(ctx)
	if err != nil {
		return nil, err
	}
	did := didutils.JoinDID(didtypes.DidMethod, namespace, req.CollectionId)
	hasDidDoc, err := q.didKeeper.HasDidDoc(ctx, did)
	if err != nil {
		return nil, err
	}
	if !hasDidDoc {
		return nil, didtypes.ErrDidDocNotFound.Wrap(did)
	}
	// Get all resources
	resources, err := q.GetResourceCollection(ctx, req.CollectionId)
	if err != nil {
		return nil, types.ErrResourceNotAvail.Wrap(err.Error())
	}

	return &types.QueryCollectionResourcesResponse{
		Resources: resources,
	}, nil
}
