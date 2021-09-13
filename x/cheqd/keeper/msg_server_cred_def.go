package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCred_def(goCtx context.Context, msg *types.MsgCreateCred_def) (*types.MsgCreateCred_defResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendCred_def(
		ctx,
		msg.Creator,
		msg.Schema_id,
		msg.Tag,
		msg.Signature_type,
		msg.Value,
	)

	return &types.MsgCreateCred_defResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateCred_def(goCtx context.Context, msg *types.MsgUpdateCred_def) (*types.MsgUpdateCred_defResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var cred_def = types.Cred_def{
		Creator:        msg.Creator,
		Id:             msg.Id,
		Schema_id:      msg.Schema_id,
		Tag:            msg.Tag,
		Signature_type: msg.Signature_type,
		Value:          msg.Value,
	}

	// Checks that the element exists
	if !k.HasCred_def(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetCred_defOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetCred_def(ctx, cred_def)

	return &types.MsgUpdateCred_defResponse{}, nil
}

func (k msgServer) DeleteCred_def(goCtx context.Context, msg *types.MsgDeleteCred_def) (*types.MsgDeleteCred_defResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasCred_def(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetCred_defOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveCred_def(ctx, msg.Id)

	return &types.MsgDeleteCred_defResponse{}, nil
}
