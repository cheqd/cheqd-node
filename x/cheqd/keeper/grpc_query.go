package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
)

var _ v1.QueryServer = Keeper{}
