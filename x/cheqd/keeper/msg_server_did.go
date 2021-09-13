package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendDid(
		ctx,
		msg.Creator,
		msg.Verkey,
		msg.Alias,
	)

	return &types.MsgCreateDidResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateDid(goCtx context.Context, msg *types.MsgUpdateDid) (*types.MsgUpdateDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var did = types.Did{
		Creator: msg.Creator,
		Id:      msg.Id,
		Verkey:  msg.Verkey,
		Alias:   msg.Alias,
	}

	// Checks that the element exists
	if !k.HasDid(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetDidOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetDid(ctx, did)

	return &types.MsgUpdateDidResponse{}, nil
}

func (k msgServer) DeleteDid(goCtx context.Context, msg *types.MsgDeleteDid) (*types.MsgDeleteDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasDid(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetDidOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveDid(ctx, msg.Id)

	return &types.MsgDeleteDidResponse{}, nil
}
