package keeper

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type msgServer struct {
	Keeper
	didKeeper didkeeper.Keeper
}

// NewMsgServerImpl returns an implementation of the x/auth MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper, cheqdKeeper didkeeper.Keeper) types.MsgServer {
	return &msgServer{
		Keeper:    keeper,
		didKeeper: cheqdKeeper,
	}
}

var _ types.MsgServer = msgServer{}
