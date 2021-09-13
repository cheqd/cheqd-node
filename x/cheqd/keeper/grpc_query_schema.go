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

func (k Keeper) SchemaAll(c context.Context, req *types.QueryAllSchemaRequest) (*types.QueryAllSchemaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var schemas []*types.Schema
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	schemaStore := prefix.NewStore(store, types.KeyPrefix(types.SchemaKey))

	pageRes, err := query.Paginate(schemaStore, req.Pagination, func(key []byte, value []byte) error {
		var schema types.Schema
		if err := k.cdc.UnmarshalBinaryBare(value, &schema); err != nil {
			return err
		}

		schemas = append(schemas, &schema)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSchemaResponse{Schema: schemas, Pagination: pageRes}, nil
}

func (k Keeper) Schema(c context.Context, req *types.QueryGetSchemaRequest) (*types.QueryGetSchemaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var schema types.Schema
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasSchema(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSchemaIDBytes(req.Id)), &schema)

	return &types.QueryGetSchemaResponse{Schema: &schema}, nil
}
