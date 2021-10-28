package keeper

import (
	"context"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

func (k msgServer) CreateCredDef(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateCredDefResponse, error) {
	return nil, types.ErrNotImplemented.Wrap("CreateCredDef")

	/*
		ctx := sdk.UnwrapSDKContext(goCtx)
		prefix := types.DidPrefix + ":" + types.DidMethod + ":" + ctx.ChainID() + ":"

		var credDefMsg types.MsgCreateCredDef
		err := k.cdc.Unmarshal(msg.Data.Value, &credDefMsg)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}

		if err := credDefMsg.ValidateBasic(prefix); err != nil {
			return nil, err
		}

		if err := k.VerifySignature(&ctx, msg, credDefMsg.GetSigners()); err != nil {
			return nil, err
		}

		// Checks that the did doesn't exist
		if err := k.EnsureDidIsNotUsed(ctx, credDefMsg.GetDid()); err != nil {
			return nil, err
		}

		// TODO: implement cred def
	*/
}
