package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"github.com/golang/protobuf/jsonpb" //nolint

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

func (k Keeper) DidJson(c context.Context, req *types.QueryGetDidJsonRequest) (*types.QueryGetDidJsonResponse, error) {
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

	var m jsonpb.Marshaler
	didJson, err := m.MarshalToString(did)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetDidJsonResponse{Value: didJson}, nil
}
