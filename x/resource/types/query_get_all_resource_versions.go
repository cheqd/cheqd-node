package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetAllResourceVersionsRequest) Normalize() {
	query.CollectionId = utils.NormalizeId(query.CollectionId)
}
