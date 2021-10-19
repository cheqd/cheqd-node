package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateSchema(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	schemaMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgCreateSchema)

	if !isMsgIdentity {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
	}

	if err := schemaMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, msg, schemaMsg.GetSigners()); err != nil {
		return nil, err
	}

	// Checks that the element exists
	if err := k.HasDidDoc(ctx, schemaMsg.GetDid()); err != nil {
		return nil, err
	}

	k.AppendSchema(
		ctx,
		schemaMsg.Id,
		schemaMsg.Type,
		schemaMsg.Name,
		schemaMsg.Version,
		schemaMsg.AttrNames,
		schemaMsg.Controller,
	)

	return &types.MsgCreateSchemaResponse{
		Id: schemaMsg.Id,
	}, nil
}
