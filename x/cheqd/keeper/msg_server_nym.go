package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd-id/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateNym(goCtx context.Context, msg *types.MsgCreateNym) (*types.MsgCreateNymResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendNym(
		ctx,
		msg.Creator,
		msg.Alias,
		msg.Verkey,
		msg.Did,
		msg.Role,
	)

	return &types.MsgCreateNymResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateNym(goCtx context.Context, msg *types.MsgUpdateNym) (*types.MsgUpdateNymResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var nym = types.Nym{
		Creator: msg.Creator,
		Id:      msg.Id,
		Alias:   msg.Alias,
		Verkey:  msg.Verkey,
		Did:     msg.Did,
		Role:    msg.Role,
	}

	// Checks that the element exists
	if !k.HasNym(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetNymOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetNym(ctx, nym)

	return &types.MsgUpdateNymResponse{}, nil
}

func (k msgServer) DeleteNym(goCtx context.Context, msg *types.MsgDeleteNym) (*types.MsgDeleteNymResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasNym(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetNymOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveNym(ctx, msg.Id)

	return &types.MsgDeleteNymResponse{}, nil
}
