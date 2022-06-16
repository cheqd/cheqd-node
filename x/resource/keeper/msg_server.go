package keeper

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type msgServer struct {
	Keeper
	cheqdKeeper cheqdkeeper.Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper, cheqdKeeper cheqdkeeper.Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
		cheqdKeeper: cheqdKeeper,
	}
}

var _ types.MsgServer = msgServer{}
