package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
)

type msgServer struct {
	Keeper
}

// NewMsgSercheqdpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgSercheqdpl(keeper Keeper) v1.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ v1.MsgServer = msgServer{}
