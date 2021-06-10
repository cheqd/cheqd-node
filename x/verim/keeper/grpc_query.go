package keeper

import (
	"github.com/verim-id/verim-node/x/verim/types"
)

var _ types.QueryServer = Keeper{}
