package keeper

import (
	"github.com/cheqd/cheqd-node/x/did/types"
)

type QueryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

var _ types.QueryServer = QueryServer{}
