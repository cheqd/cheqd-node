package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
