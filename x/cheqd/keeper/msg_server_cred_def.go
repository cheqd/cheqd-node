package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCredDef(goCtx context.Context, msg *types.MsgCreateCredDef) (*types.MsgCreateCredDefResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	id := k.AppendCredDef(
		ctx,
		msg.Id,
		msg.SchemaId,
		msg.Tag,
		msg.SignatureType,
		msg.Value,
	)

	return &types.MsgCreateCredDefResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateCredDef(goCtx context.Context, msg *types.MsgUpdateCredDef) (*types.MsgUpdateCredDefResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var credDef = types.CredDef{
		Creator:       msg.Creator,
		Id:            msg.Id,
		Schema_id:     msg.Schema_id,
		Tag:           msg.Tag,
		SignatureType: msg.SignatureType,
		Value:         msg.Value,
	}

	// Checks that the element exists
	if !k.HasCredDef(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetCredDefOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetCredDef(ctx, credDef)

	return &types.MsgUpdateCredDefResponse{}, nil
}

func (k msgServer) DeleteCredDef(goCtx context.Context, msg *types.MsgDeleteCredDef) (*types.MsgDeleteCredDefResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasCredDef(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetCredDefOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveCredDef(ctx, msg.Id)

	return &types.MsgDeleteCredDefResponse{}, nil
}
