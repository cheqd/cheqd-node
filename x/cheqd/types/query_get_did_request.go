package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetDidRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
}
