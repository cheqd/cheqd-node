package types

import "github.com/cheqd/cheqd-node/x/cheqd/utils"

func (query *QueryGetResourceMetadataRequest) Normalize() {
	query.CollectionId = utils.NormalizeId(query.CollectionId)
	query.Id = utils.NormalizeUUID(query.Id)
}
