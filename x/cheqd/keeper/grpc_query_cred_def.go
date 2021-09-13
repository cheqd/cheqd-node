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

func (k Keeper) Cred_defAll(c context.Context, req *types.QueryAllCred_defRequest) (*types.QueryAllCred_defResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var cred_defs []*types.Cred_def
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	cred_defStore := prefix.NewStore(store, types.KeyPrefix(types.Cred_defKey))

	pageRes, err := query.Paginate(cred_defStore, req.Pagination, func(key []byte, value []byte) error {
		var cred_def types.Cred_def
		if err := k.cdc.UnmarshalBinaryBare(value, &cred_def); err != nil {
			return err
		}

		cred_defs = append(cred_defs, &cred_def)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCred_defResponse{Cred_def: cred_defs, Pagination: pageRes}, nil
}

func (k Keeper) Cred_def(c context.Context, req *types.QueryGetCred_defRequest) (*types.QueryGetCred_defResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var cred_def types.Cred_def
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasCred_def(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cred_defKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCred_defIDBytes(req.Id)), &cred_def)

	return &types.QueryGetCred_defResponse{Cred_def: &cred_def}, nil
}
