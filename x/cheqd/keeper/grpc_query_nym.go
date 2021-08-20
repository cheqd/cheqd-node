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

func (k Keeper) NymAll(c context.Context, req *types.QueryAllNymRequest) (*types.QueryAllNymResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nyms []*types.Nym
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	nymStore := prefix.NewStore(store, types.KeyPrefix(types.NymKey))

	pageRes, err := query.Paginate(nymStore, req.Pagination, func(key []byte, value []byte) error {
		var nym types.Nym
		if err := k.cdc.Unmarshal(value, &nym); err != nil {
			return err
		}

		nyms = append(nyms, &nym)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNymResponse{Nym: nyms, Pagination: pageRes}, nil
}

func (k Keeper) Nym(c context.Context, req *types.QueryGetNymRequest) (*types.QueryGetNymResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nym types.Nym
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasNym(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NymKey))
	k.cdc.MustUnmarshal(store.Get(GetNymIDBytes(req.Id)), &nym)

	return &types.QueryGetNymResponse{Nym: &nym}, nil
}
