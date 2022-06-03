package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Did(c context.Context, req *types.QueryGetDidRequest) (*types.QueryGetDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	stateValue, err := k.GetDid(&ctx, req.Id)
	if err != nil {
		return nil, err
	}

	did, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return nil, err
	}

	return &types.QueryGetDidResponse{Did: did, Metadata: stateValue.Metadata}, nil
}
