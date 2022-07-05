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

func (m queryServer) AllResourceVersions(c context.Context, req *types.QueryGetAllResourceVersionsRequest) (*types.QueryGetAllResourceVersionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Validate corresponding DIDDoc exists
	namespace := m.cheqdKeeper.GetDidNamespace(&ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, req.CollectionId)
	if !m.cheqdKeeper.HasDid(&ctx, did) {
		return nil, cheqdtypes.ErrDidDocNotFound.Wrap(did)
	}

	// Get all versions
	versions := m.GetAllResourceVersions(&ctx, req.CollectionId, req.Name)

	return &types.QueryGetAllResourceVersionsResponse{
		Resources: versions,
	}, nil
}
