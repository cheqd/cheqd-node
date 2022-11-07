package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllDidDocVersions(c context.Context, req *types.QueryGetAllDidDocVersionsRequest) (*types.QueryGetAllDidDocVersionsResponse, error) {
	// TODO: Implement
	return nil, status.Error(codes.Unimplemented, "method AllDidDocVersions not implemented")
}
