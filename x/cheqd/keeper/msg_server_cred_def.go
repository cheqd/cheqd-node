package keeper

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCredDef(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateCredDefResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefix := types.DidPrefix + ctx.ChainID() + ":"

	credDefMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgCreateCredDef)
	if !isMsgIdentity {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
	}

	if err := credDefMsg.ValidateBasic(prefix); err != nil {
		return nil, err
	}

	if err := k.VerifySignature(&ctx, msg, credDefMsg.GetSigners()); err != nil {
		return nil, err
	}

	// Checks that the element exists
	if err := k.HasDidDoc(ctx, credDefMsg.GetDid()); err != nil {
		return nil, err
	}

	switch value := credDefMsg.Value.(type) {
	case *types.MsgCreateCredDef_ClType:
		k.AppendCredDef(
			ctx,
			credDefMsg.Id,
			credDefMsg.SchemaId,
			credDefMsg.Tag,
			credDefMsg.Type,
			(*types.CredDef_ClType)(value),
			credDefMsg.Controller,
		)

		return &types.MsgCreateCredDefResponse{
			Id: credDefMsg.Id,
		}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidCredDefValue, "unsupported cred def value")
	}
}
