package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetCollectionResourcesRequest) Normalize() *QueryGetCollectionResourcesRequest {
	return &QueryGetCollectionResourcesRequest{
		CollectionId: utils.NormalizeIdentifier(query.CollectionId),
	}
}
