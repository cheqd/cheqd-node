package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ types.QueryServer = Keeper{}
