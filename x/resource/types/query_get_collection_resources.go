package types

import (
	"github.com/cheqd/cheqd-node/x/did/utils"
)

func (query *QueryGetCollectionResourcesRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
}
