package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetCollectionResourcesRequest) Normalize() {
	query.CollectionId = utils.NormalizeId(query.CollectionId)
}
