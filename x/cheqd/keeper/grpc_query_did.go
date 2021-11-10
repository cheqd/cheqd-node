package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Did(c context.Context, req *v1.QueryGetDidRequest) (*v1.QueryGetDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	state, err := k.GetDid(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	did, err := state.GetDid()
	if err != nil {
		return nil, err
	}

	return &v1.QueryGetDidResponse{Did: did, Metadata: state.Metadata}, nil
}
