package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) baseapp.MsgServiceHandler {
	msgServer := NewMsgServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateDidDoc:
			res, err := msgServer.CreateDidDoc(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUpdateDidDoc:
			res, err := msgServer.UpdateDidDoc(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDeactivateDidDoc:
			res, err := msgServer.DeactivateDidDoc(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
