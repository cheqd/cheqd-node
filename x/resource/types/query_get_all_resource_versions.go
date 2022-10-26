package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetAllResourceVersionsRequest) Normalize() *QueryGetAllResourceVersionsRequest {
	return &QueryGetAllResourceVersionsRequest{
		CollectionId: utils.NormalizeIdentifier(query.CollectionId),
		Name:         query.Name,
		ResourceType: query.ResourceType,
	}
}
