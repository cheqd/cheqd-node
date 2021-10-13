package keeper

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateSchema(goCtx context.Context, msg *types.MsgWriteRequest) (*types.MsgCreateSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	schemaMsg, isMsgIdentity := msg.Data.GetCachedValue().(*types.MsgCreateSchema)

	if !isMsgIdentity {
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	// Checks that the element exists
	if err := k.HasDidDoc(ctx, schemaMsg.Id); err != nil {
		return nil, err
	}

	k.AppendSchema(
		ctx,
		schemaMsg.Id,
		schemaMsg.Name,
		schemaMsg.Version,
		schemaMsg.AttrNames,
	)

	return &types.MsgCreateSchemaResponse{
		Id: schemaMsg.Id,
	}, nil
}
