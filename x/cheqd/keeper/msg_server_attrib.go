package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateAttrib(goCtx context.Context, msg *types.MsgCreateAttrib) (*types.MsgCreateAttribResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendAttrib(
		ctx,
		msg.Creator,
		msg.Did,
		msg.Raw,
	)

	return &types.MsgCreateAttribResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateAttrib(goCtx context.Context, msg *types.MsgUpdateAttrib) (*types.MsgUpdateAttribResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var attrib = types.Attrib{
		Creator: msg.Creator,
		Id:      msg.Id,
		Did:     msg.Did,
		Raw:     msg.Raw,
	}

	// Checks that the element exists
	if !k.HasAttrib(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetAttribOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetAttrib(ctx, attrib)

	return &types.MsgUpdateAttribResponse{}, nil
}

func (k msgServer) DeleteAttrib(goCtx context.Context, msg *types.MsgDeleteAttrib) (*types.MsgDeleteAttribResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasAttrib(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetAttribOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAttrib(ctx, msg.Id)

	return &types.MsgDeleteAttribResponse{}, nil
}
