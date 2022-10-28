package keeper

import (
	"context"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cheqd/cheqd-node/x/resource/types"
)

func (q queryServer) Resource(c context.Context, req *types.QueryGetResourceRequest) (*types.QueryGetResourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	req.Normalize()

	ctx := sdk.UnwrapSDKContext(c)

	// Validate corresponding DIDDoc exists
	namespace := q.cheqdKeeper.GetDidNamespace(&ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, req.CollectionId)
	if !q.cheqdKeeper.HasDid(&ctx, did) {
		return nil, cheqdtypes.ErrDidDocNotFound.Wrap(did)
	}

	resource, err := q.GetResource(&ctx, req.CollectionId, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetResourceResponse{
		Resource: &resource,
	}, nil
}
