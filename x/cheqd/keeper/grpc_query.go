package keeper

import (
	"github.com/cheqd-id/cheqd-node/x/cheqd/types"
)

var _ types.QueryServer = Keeper{}
