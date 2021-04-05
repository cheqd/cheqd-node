package keeper

import (
	"github.com/verim-id/verim-cosmos/x/verimcosmos/types"
)

var _ types.QueryServer = Keeper{}
