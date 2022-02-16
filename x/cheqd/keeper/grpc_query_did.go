package keeper

import (
	"context"
	"encoding/json"
	"strings"

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

	didJson, err := json.Marshal(did)
	if err != nil {
		return nil, err
	}

	formatedDidJson := strings.ReplaceAll(string(didJson), "key_agreement", "")
	formatedDidJson = strings.ReplaceAll(string(formatedDidJson), "context", "@context")

	//result := strings.ReplaceAll(string(bz), "key_agreement", "")
	//println(result)

	//return []byte("result"), nil

	return &types.QueryGetDidJsonResponse{Value: formatedDidJson}, nil
}
