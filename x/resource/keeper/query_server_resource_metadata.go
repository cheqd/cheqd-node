package keeper

import (
	"context"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cheqd/cheqd-node/x/resource/types"
)

func (q queryServer) ResourceMetadata(c context.Context, req *types.QueryResourceMetadataRequest) (*types.QueryResourceMetadataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	// Validate corresponding DIDDoc exists
	namespace, err := q.didKeeper.GetDidNamespace(c)
	if err != nil {
		return nil, err
	}
	did := didutils.JoinDID(didtypes.DidMethod, namespace, req.CollectionId)
	hasDidDoc, err := q.didKeeper.HasDidDoc(c, did)
	if err != nil {
		return nil, err
	}
	if !hasDidDoc {
		return nil, didtypes.ErrDidDocNotFound.Wrap(did)
	}

	metadata, err := q.GetResourceMetadata(c, req.CollectionId, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryResourceMetadataResponse{
		Resource: &metadata,
	}, nil
}
