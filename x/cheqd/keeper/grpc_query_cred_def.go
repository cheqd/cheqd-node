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

func (k Keeper) CredDefAll(c context.Context, req *types.QueryAllCredDefRequest) (*types.QueryAllCredDefResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var credDefs []*types.CredDef
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	credDefStore := prefix.NewStore(store, types.KeyPrefix(types.CredDefKey))

	pageRes, err := query.Paginate(credDefStore, req.Pagination, func(key []byte, value []byte) error {
		var credDef types.CredDef
		if err := k.cdc.UnmarshalBinaryBare(value, &credDef); err != nil {
			return err
		}

		credDefs = append(credDefs, &credDef)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCredDefResponse{CredDef: credDefs, Pagination: pageRes}, nil
}

func (k Keeper) CredDef(c context.Context, req *types.QueryGetCredDefRequest) (*types.QueryGetCredDefResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var credDef types.CredDef
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasCredDef(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCredDefIDBytes(req.Id)), &credDef)

	return &types.QueryGetCredDefResponse{CredDef: &credDef}, nil
}
