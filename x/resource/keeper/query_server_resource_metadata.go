package keeper

import (
	"context"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cheqd/cheqd-node/x/resource/types"
)

func (q queryServer) ResourceMetadata(c context.Context, req *types.QueryGetResourceMetadataRequest) (*types.QueryGetResourceMetadataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	// Validate corresponding DIDDoc exists
	namespace := q.didKeeper.GetDidNamespace(&ctx)
	did := didutils.JoinDID(didtypes.DidMethod, namespace, req.CollectionId)
	if !q.didKeeper.HasDidDoc(&ctx, did) {
		return nil, didtypes.ErrDidDocNotFound.Wrap(did)
	}

	metadata, err := q.GetResourceMetadata(&ctx, req.CollectionId, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetResourceMetadataResponse{
		Resource: &metadata,
	}, nil
}
