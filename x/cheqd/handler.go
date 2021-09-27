package cheqd

import (
	"fmt"
	"reflect"

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
		parsed_msg, isMsgIdentity := msg.(*types.MsgWriteRequest)
		if !isMsgIdentity {
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
		switch parsed_msg.Data.TypeUrl {
		// this line is used by starport scaffolding # 1
		case reflect.TypeOf(types.MsgCreateCredDef{}).Name():
			res, err := msgServer.CreateCredDef(sdk.WrapSDKContext(ctx), parsed_msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case reflect.TypeOf(types.MsgCreateSchema{}).Name():
			res, err := msgServer.CreateSchema(sdk.WrapSDKContext(ctx), parsed_msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case reflect.TypeOf(types.MsgCreateDid{}).Name():
			res, err := msgServer.CreateDid(sdk.WrapSDKContext(ctx), parsed_msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case reflect.TypeOf(types.MsgUpdateDid{}).Name():
			res, err := msgServer.UpdateDid(sdk.WrapSDKContext(ctx), parsed_msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
