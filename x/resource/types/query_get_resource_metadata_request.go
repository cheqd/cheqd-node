package types

import "github.com/cheqd/cheqd-node/x/did/utils"

func (query *QueryResourceMetadataRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
	query.Id = utils.NormalizeUUID(query.Id)
}
