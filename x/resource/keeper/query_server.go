package keeper

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type queryServer struct {
	Keeper
	cheqdKeeper cheqdkeeper.Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper, cheqdKeeper cheqdkeeper.Keeper) types.QueryServer {
	return &queryServer{
		Keeper:      keeper,
		cheqdKeeper: cheqdKeeper,
	}
}

var _ types.QueryServer = queryServer{}
