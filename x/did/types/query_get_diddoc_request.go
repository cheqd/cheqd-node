package types

import (
	"github.com/cheqd/cheqd-node/x/did/utils"
)

func (query *QueryGetDidDocRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
}
