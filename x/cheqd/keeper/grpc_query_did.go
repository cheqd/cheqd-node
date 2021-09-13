package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidAll(c context.Context, req *types.QueryAllDidRequest) (*types.QueryAllDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var dids []*types.Did
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	didStore := prefix.NewStore(store, types.KeyPrefix(types.DidKey))

	pageRes, err := query.Paginate(didStore, req.Pagination, func(key []byte, value []byte) error {
		var did types.Did
		if err := k.cdc.UnmarshalBinaryBare(value, &did); err != nil {
			return err
		}

		dids = append(dids, &did)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDidResponse{Did: dids, Pagination: pageRes}, nil
}

func (k Keeper) Did(c context.Context, req *types.QueryGetDidRequest) (*types.QueryGetDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var did types.Did
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasDid(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetDidIDBytes(req.Id)), &did)

	return &types.QueryGetDidResponse{Did: &did}, nil
}
