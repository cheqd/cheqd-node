package resource

import (
	"fmt"

	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"

	errorsmod "cosmossdk.io/errors"
	oracleKeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper, cheqdKeeper didkeeper.Keeper, oracleKeeper oracleKeeper.Keeper) baseapp.MsgServiceHandler {
	msgServer := keeper.NewMsgServer(k, cheqdKeeper, oracleKeeper)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateResource:
			res, err := msgServer.CreateResource(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
