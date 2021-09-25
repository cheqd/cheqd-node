package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateSchema(goCtx context.Context, msg *types.MsgCreateSchema) (*types.MsgCreateSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendSchema(
		ctx,
		msg.Creator,
		msg.Name,
		msg.Version,
		msg.AttrNames,
	)

	return &types.MsgCreateSchemaResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateSchema(goCtx context.Context, msg *types.MsgUpdateSchema) (*types.MsgUpdateSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var schema = types.Schema{
		Creator:   msg.Creator,
		Id:        msg.Id,
		Name:      msg.Name,
		Version:   msg.Version,
		AttrNames: msg.AttrNames,
	}

	// Checks that the element exists
	if !k.HasSchema(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetSchemaOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetSchema(ctx, schema)

	return &types.MsgUpdateSchemaResponse{}, nil
}

func (k msgServer) DeleteSchema(goCtx context.Context, msg *types.MsgDeleteSchema) (*types.MsgDeleteSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasSchema(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetSchemaOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveSchema(ctx, msg.Id)

	return &types.MsgDeleteSchemaResponse{}, nil
}
