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

func (k Keeper) AttribAll(c context.Context, req *types.QueryAllAttribRequest) (*types.QueryAllAttribResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var attribs []*types.Attrib
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	attribStore := prefix.NewStore(store, types.KeyPrefix(types.AttribKey))

	pageRes, err := query.Paginate(attribStore, req.Pagination, func(key []byte, value []byte) error {
		var attrib types.Attrib
		if err := k.cdc.UnmarshalBinaryBare(value, &attrib); err != nil {
			return err
		}

		attribs = append(attribs, &attrib)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAttribResponse{Attrib: attribs, Pagination: pageRes}, nil
}

func (k Keeper) Attrib(c context.Context, req *types.QueryGetAttribRequest) (*types.QueryGetAttribResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var attrib types.Attrib
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasAttrib(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttribKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetAttribIDBytes(req.Id)), &attrib)

	return &types.QueryGetAttribResponse{Attrib: &attrib}, nil
}
