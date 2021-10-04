package cheqd

import (
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgSercheqdpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		parsedMsg, isMsgIdentity := msg.(*types.MsgWriteRequest)
		if !isMsgIdentity {
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		if !k.Verify(ctx, parsedMsg) {
			errMsg := "Invalid signature"
			return nil, sdkerrors.Wrap(types.ErrInvalidSignature, errMsg)
		}

		switch parsedMsg.Data.TypeUrl {
		// this line is used by starport scaffolding # 1
		case "/cheqdid.cheqdnode.cheqd.MsgCreateCredDef":
			res, err := msgServer.CreateCredDef(sdk.WrapSDKContext(ctx), parsedMsg)
			return sdk.WrapServiceResult(ctx, res, err)

		case "/cheqdid.cheqdnode.cheqd.MsgCreateSchema":
			res, err := msgServer.CreateSchema(sdk.WrapSDKContext(ctx), parsedMsg)
			return sdk.WrapServiceResult(ctx, res, err)

		case "/cheqdid.cheqdnode.cheqd.MsgCreateDid":
			res, err := msgServer.CreateDid(sdk.WrapSDKContext(ctx), parsedMsg)
			return sdk.WrapServiceResult(ctx, res, err)

		case "/cheqdid.cheqdnode.cheqd.MsgUpdateDid":
			res, err := msgServer.UpdateDid(sdk.WrapSDKContext(ctx), parsedMsg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
