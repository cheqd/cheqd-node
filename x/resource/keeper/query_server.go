package keeper

import (
	"github.com/cheqd/cheqd-node/x/resource/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}
