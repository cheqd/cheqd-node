package types

import (
	"github.com/cheqd/cheqd-node/x/did/utils"
)

func (query *QueryDidDocVersionRequest) Normalize() {
	query.Id = utils.NormalizeDID(query.Id)
	query.Version = utils.NormalizeUUID(query.Version)
}
