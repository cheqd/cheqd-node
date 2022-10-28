package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetResourceRequest) Normalize() {
	query.CollectionId = utils.NormalizeId(query.CollectionId)
	query.Id = utils.NormalizeId(query.Id)
}
