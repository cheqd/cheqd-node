package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetResourceRequest) Normalize() *QueryGetResourceRequest {
	return &QueryGetResourceRequest{
		CollectionId: utils.NormalizeIdentifier(query.CollectionId),
		Id:           utils.NormalizeIdentifier(query.Id),
	}
}
