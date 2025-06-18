package types

import (
	"github.com/cheqd/cheqd-node/x/did/utils"
)

func (query *QueryResourceRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
	query.Id = utils.NormalizeUUID(query.Id)
}

func (query *QueryLatestResourceVersionRequest) Normalize() {
	query.CollectionId = utils.NormalizeID(query.CollectionId)
	query.Name = utils.NormalizeID(query.Name)
	query.ResourceType = utils.NormalizeID(query.ResourceType)
}
