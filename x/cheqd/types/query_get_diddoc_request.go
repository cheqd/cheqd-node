package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetDidDocRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
}
