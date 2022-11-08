package keeper

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type queryServer struct {
	Keeper
	didKeeper didkeeper.Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper, cheqdKeeper didkeeper.Keeper) types.QueryServer {
	return &queryServer{
		Keeper:    keeper,
		didKeeper: cheqdKeeper,
	}
}

var _ types.QueryServer = queryServer{}
