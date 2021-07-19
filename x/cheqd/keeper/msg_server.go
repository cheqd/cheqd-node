package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

type msgServer struct {
	Keeper
}

// NewMsgSercheqdpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgSercheqdpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
